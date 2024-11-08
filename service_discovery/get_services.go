package serviceDiscovery

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/darksuei/kubeRPC/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetKubeServices(clientset *kubernetes.Clientset) error {
	namespace := os.Getenv("NAMESPACE")

	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf("error listing services: %v", err)
	}

	excludedServices, err := helpers.ParseJSONArrayFromEnv("EXCLUDE_SERVICES")

	if err != nil {
		log.Fatalf("Error parsing excluded services: %v", err)
	}

	for _, service := range services.Items {
		if helpers.StringInSlice(excludedServices, service.Name) {
			continue
		}

		RegisterService(&service)
	}
	return nil
}
