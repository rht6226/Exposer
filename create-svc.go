package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *controller) CreateService(dep *appsv1.Deployment, ctx context.Context, ns, name string) (*corev1.Service, error) {
	labels := dep.Spec.Template.Labels

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.Name,
			Namespace: dep.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}

	createdSvc, err := c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})

	if err != nil {
		fmt.Println("[ERROR]: Error creating service.")
		return nil, err
	}

	return createdSvc, nil
}
