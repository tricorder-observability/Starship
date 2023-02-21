package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	matav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

var (
	DefaultNamespace   = "tricorder"
	DefaultServiceName = "api-server"
)

func GetService(client *k8s.Clientset, namespace, serviceName string) (*v1.Service, error) {
	return client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, matav1.GetOptions{})
}

func GetStarshipService(client *k8s.Clientset) (*v1.Service, error) {
	return GetService(client, DefaultNamespace, DefaultServiceName)
}

func getStarshipServiceURL(client *k8s.Clientset) (string, error) {
	service, err := GetStarshipService(client)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}

func getStarshipServicePort(client *k8s.Clientset) (int32, error) {
	service, err := GetStarshipService(client)
	if err != nil {
		return 0, err
	}
	for _, port := range service.Spec.Ports {
		if port.Name == "serverhttp" {
			return port.Port, nil
		}
	}
	return service.Spec.Ports[0].Port, nil
}

func GetAPIAddress() (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", err
	}
	ip, err := getStarshipServiceURL(client)
	if err != nil {
		return "", err
	}
	port, err := getStarshipServicePort(client)
	if err != nil {
		return "", err
	}
	apiAddress := fmt.Sprintf("%s:%d", ip, port)
	return apiAddress, nil
}
