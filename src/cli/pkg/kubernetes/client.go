package kubernetes

import (
	"fmt"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetConfig returns the kubernetes config.
func GetConfig() (*rest.Config, error) {
	// get config
	// 1. in-cluster config
	// 2. `KUBECONFIG` env
	// 3. `~/.kube/config`
	// 4. /etc/kubernetes/admin.conf
	// TODO: add support for flags —-kubeconfig
	config, err := rest.InClusterConfig()
	if err == nil && config != nil {
		return config, nil
	}
	configLocal := clientcmd.NewDefaultClientConfigLoadingRules()
	startConfig, err := configLocal.GetStartingConfig()
	if err != nil {
		return nil, err
	}
	return clientcmd.NewDefaultClientConfig(*startConfig, nil).ClientConfig()
}

// GetClientForConfig returns the kubernetes client for the given config.
func GetClientForConfig(config *rest.Config) (*k8s.Clientset, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	return k8s.NewForConfig(config)
}

// NewClient returns the kubernetes client.
func NewClient() (*k8s.Clientset, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return GetClientForConfig(config)
}
