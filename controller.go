package main

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	applisters "k8s.io/client-go/listers/apps/v1"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

// type controller
type controller struct {
	clientset             kubernetes.Interface            // Kubernetes Clinetset to interact with different API versions
	deploymentLister      applisters.DeploymentLister     // DeploymentLister to get all the deployments
	deploymentCacheSynced cache.InformerSynced            // If the chache is synced or not
	queue                 workqueue.RateLimitingInterface // Queue to store the jobs
}

// Create a new controller by providing the clientset
// and an informer which looks for changes in Deployment resource.
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

// Run runs the controller in a go routine till the stop channel is closed.
// It also ensures that the cache is synced before the worker is called.
func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println("[INFO]: Starting controller...")

	if !cache.WaitForCacheSync(ch, c.deploymentCacheSynced) {
		fmt.Println("[ERROR]: Waiting for chace to be synced ...")
	}

	go wait.Until(c.Worker, 1*time.Second, ch)

	<-ch
}

// Worker Infinitely calls the process item function so that the
// controller can keep processing items added to the queue.
func (c *controller) Worker() {
	for c.processItem() {
	}
}

// processItem takes items one ata a time from the queue.
// For each item it calls syncDeploymnt with proper argumets
// to ensure that proper services, ingress are created/deleted.
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

	// check if the object has been deleted from the k8s cluster
	ctx := context.Background()
	_, err = c.clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("[INFO]: Deployment %s was deleted\n", name)

		// delete service
		err = c.clientset.CoreV1().Services(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("[ERROR] deleting the service\n", err.Error())
			return false
		}

		// delete ingress
		err = c.clientset.NetworkingV1().Ingresses(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Println("[ERROR] deleting the ingress\n", err.Error())
			return false
		}

		return true
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry
		fmt.Println("[ERROR]: syncing deployments\n", err.Error())
		return false
	}

	return true
}

// Creates the service and ingress for the givn deployment.
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
