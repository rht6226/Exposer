package main

import (
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	applisters "k8s.io/client-go/listers/apps/v1"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset            kubernetes.Interface
	deploymentLister     applisters.DeploymentLister
	deplymentCacheSynced cache.InformerSynced
	queue                workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface,
	depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:            clientset,
		deploymentLister:     depInformer.Lister(),
		deplymentCacheSynced: depInformer.Informer().HasSynced,
		queue:                workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "expose"),
	}

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    HandleAdd,
			DeleteFunc: HandleDelete,
		},
	)

	return c
}
