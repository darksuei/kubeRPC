package serviceDiscovery

import (
	"context"
	"log"
	"os"

	"github.com/darksuei/kubeRPC/helpers"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func WatchKubeServices(clientset *kubernetes.Clientset) {
	namespace := os.Getenv("NAMESPACE")

	watch, err := clientset.CoreV1().Services(namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error setting up watch: %v", err)
		return
	}

	for event := range watch.ResultChan() {
		svc := event.Object.(*v1.Service)

		excludedServices, err := helpers.ParseJSONArrayFromEnv("EXCLUDE_SERVICES")

		if err != nil {
			log.Fatalf("Error parsing excluded services: %v", err)
		}

		if helpers.StringInSlice(excludedServices, svc.Name) {
			continue
		}

		switch event.Type {
		case "ADDED":
			RegisterService(svc)

		case "DELETED":
			RemoveService(svc)
		}
	}
}
