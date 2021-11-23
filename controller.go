package main

import (
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
			AddFunc:    HandleAdd,
			DeleteFunc: HandleDelete,
		},
	)

	return c
}

func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println("[INFO]: Starting controller...")
	if !cache.WaitForCacheSync(ch, c.deploymentCacheSynced) {
		fmt.Println("[ERROR]: Waiting for chace to be synced ...")
	}

	wait.Until(c.Worker, 1*time.Second, ch)
}

func (c *controller) Worker() {
	fmt.Println("[INFO]: Working ...")
}
