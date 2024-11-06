package service_discovery

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetKubeServices(clientset *kubernetes.Clientset) error {
	namespace := os.Getenv("NAMESPACE")

	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("error listing services: %v", err)
	}

	for _, service := range services.Items {
		fmt.Printf("Found service: %s\n", service.Name)
		// All services that have been found should be registered in Redis
		// We need the service name, hostname and port
	}
	return nil
}
