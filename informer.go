package main

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

// Create an informer for given clientset.
func GetInformer(clientset *kubernetes.Clientset) informers.SharedInformerFactory {
	return informers.NewSharedInformerFactory(clientset, 10*time.Minute)
}
