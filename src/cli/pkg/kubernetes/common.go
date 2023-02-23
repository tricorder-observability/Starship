package kubernetes

import (
	"context"
	"fmt"

	"github.com/tricorder/src/utils/errors"

	v1 "k8s.io/api/core/v1"
	matav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

const (
	DefaultNamespace       = "tricorder"
	DefaultServiceName     = "api-server"
	DefaultServicePortName = "serverhttp"
)

// getService returns a Service by name and namespace
func getService(client *k8s.Clientset, namespace, serviceName string) (*v1.Service, error) {
	return client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, matav1.GetOptions{})
}

// GetStarshipService returns the starship Service by name and namespace
func getStarshipService(client *k8s.Clientset) (*v1.Service, error) {
	return getService(client, DefaultNamespace, DefaultServiceName)
}

// getStarshipServiceURL returns the starship Service ClusterIP
func getStarshipServiceURL(client *k8s.Clientset) (string, error) {
	service, err := getStarshipService(client)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}

// getStarshipServicePort returns the starship Service serverhttp port
func getStarshipServicePort(client *k8s.Clientset) (int32, error) {
	service, err := getStarshipService(client)
	if err != nil {
		return 0, err
	}
	for _, port := range service.Spec.Ports {
		if port.Name == DefaultServicePortName {
			return port.Port, nil
		}
	}
	return service.Spec.Ports[0].Port, nil
}

// GetStarshipAPIAddress returns the starship Service ClusterIP and serverhttp port
func GetStarshipAPIAddress() (string, error) {
	client, err := NewClient()
	if err != nil {
		return "", errors.Wrap("init kubernetes client for config", "new client", err)
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
