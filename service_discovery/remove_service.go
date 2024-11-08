package serviceDiscovery

import (
	"context"
	"fmt"
	"log"

	"github.com/darksuei/kubeRPC/config"
	corev1 "k8s.io/api/core/v1"
)

func RemoveService(service *corev1.Service) error {
	fields, err := config.Rdb.HGetAll("service:" + service.Name).Result()

	if err != nil || len(fields) == 0 {
		return nil
	}

	fmt.Printf("Removing service: %s...\n", service.Name)

	redisKey := "service:" + service.Name

	err = config.Rdb.WithContext(context.Background()).Del(redisKey).Err()

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return nil
}
