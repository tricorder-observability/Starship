package kubernetes

import (
	"context"
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

func GetStarshipServiceURL(client *k8s.Clientset) (string, error) {
	service, err := GetStarshipService(client)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}
