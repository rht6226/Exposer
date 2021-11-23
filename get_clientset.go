package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(kubeconfig *string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Println("[ERROR]: No kubeconfig file provided.\n", err.Error())
		fmt.Println("[INFO]: Now trying to load config from the service account")

		config, err = rest.InClusterConfig()

		if err != nil {
			fmt.Println("[ERROR]: Cannot load config from the cluster.")
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Println("[ERROR]: Cannot create clientset from config.")
		return nil, err
	}

	return clientset, nil
}