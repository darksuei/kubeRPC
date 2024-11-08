package serviceDiscovery

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/darksuei/kubeRPC/config"
	corev1 "k8s.io/api/core/v1"
)

func RegisterService(service *corev1.Service) error {
	fields, err := config.Rdb.HGetAll("service:" + service.Name).Result()

	if err == nil && len(fields) > 0 {
		return nil
	}

	fmt.Printf("Registering service: %s...\n", service.Name)

	redisKey := "service:" + service.Name
	host := service.Name + "." + os.Getenv("NAMESPACE") + ".svc.cluster.local"

	// Set the service
	err = config.Rdb.WithContext(context.Background()).HSet(redisKey, "serviceName", service.Name).Err()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Set the host
	err = config.Rdb.WithContext(context.Background()).HSet(redisKey, "host", host).Err()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return nil
}
