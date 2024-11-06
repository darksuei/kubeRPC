package service_discovery

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func WatchKubeServices(clientset *kubernetes.Clientset) {
	watch, err := clientset.CoreV1().Services("your-namespace").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Errorf("error setting up watch: %v", err)
		return
	}

	for event := range watch.ResultChan() {
		switch event.Type {
		case "ADDED":
			svc := event.Object.(*v1.Service)
			fmt.Printf("New service added: %s\n", svc.Name)
			// Register the new service

		case "DELETED":
			svc := event.Object.(*v1.Service)
			fmt.Printf("Service deleted: %s\n", svc.Name)
			// Remove the service from the registry

		case "MODIFIED":
			svc := event.Object.(*v1.Service)
			fmt.Printf("Service modified: %s\n", svc.Name)
			// Update the service in the registry if necessary
		}
	}
}
