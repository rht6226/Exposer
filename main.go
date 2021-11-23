package main

import (
	"flag"

	"k8s.io/client-go/kubernetes"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/ranand/.kube/config",
		"Location of the kubeconfig file")

	clientSet, err := GetClientSet(kubeconfig)
	if err != nil {
		panic(err)
	}

	informer := GetInformer(clientSet.(*kubernetes.Clientset))

	ch := make(chan struct{})
	c := NewController(clientSet, informer.Apps().V1().Deployments())

	informer.Start(ch)
	c.Run(ch)
}
