package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	applisters "k8s.io/client-go/listers/apps/v1"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset             kubernetes.Interface
	deploymentLister      applisters.DeploymentLister
	deploymentCacheSynced cache.InformerSynced
	queue                 workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface,
	depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:             clientset,
		deploymentLister:      depInformer.Lister(),
		deploymentCacheSynced: depInformer.Informer().HasSynced,
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "expose"),
	}

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDelete,
		},
	)

	return c
}

func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println("[INFO]: Starting controller...")
	if !cache.WaitForCacheSync(ch, c.deploymentCacheSynced) {
		fmt.Println("[ERROR]: Waiting for chace to be synced ...")
	}

	go wait.Until(c.Worker, 1*time.Second, ch)

	<-ch
}

func (c *controller) Worker() {
	for c.processItem() {
	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()

	if shutdown {
		return false
	}

	defer c.queue.Forget(item)

	key, err := cache.MetaNamespaceKeyFunc(item)

	if err != nil {
		fmt.Printf("[ERROR]: unable to get key from cache.\n%s", err.Error())
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("[ERROR]; splitting key into namespace and name.\n%s", err.Error())
		return false
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry
		fmt.Println("[ERROR]: syncing deployments\n", err.Error())
		return false
	}

	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()

	dep, err := c.deploymentLister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Println("[ERROR]: Error getting deployment from lister.")
		return err
	}

	svc, err := c.CreateService(dep, ctx, ns, name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = c.CreateIngress(ctx, svc)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
