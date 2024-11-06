package service_discovery

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Println("Error creating in-cluster config, trying local config..")
		kubeconfig := os.Getenv("KUBECONFIG_FILE")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

		if err != nil {
			return nil, fmt.Errorf("Error creating local config: %v", err)
		}
	}
	return kubernetes.NewForConfig(config)
}
