package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(kubeconfig *string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Println("[ERROR]: Building config from file\n", err.Error())

		// For this to work we are using default service account.
		// This account does not have lots of previledges so we will have to create a role and bind it.
		config, err = rest.InClusterConfig()

		if err != nil {
			fmt.Println("[ERROR]: Building reading InClusterConfig")
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("[ERROR]: Creating new clientset")
		return nil, err
	}

	return clientset, nil
}
