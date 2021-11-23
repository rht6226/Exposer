package main

import (
	"flag"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/ranand/.kube/config",
		"Location of the kubeconfig file")

	clientSet, err := GetClientSet(kubeconfig)
	if err != nil {
		panic(err)
	}

	informer := GetInformer(clientSet)

	print(informer)
}
