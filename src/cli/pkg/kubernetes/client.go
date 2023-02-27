// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package kubernetes

import (
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getConfig returns the kubernetes config.
func getConfig() (*rest.Config, error) {
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

// NewClient returns the kubernetes client.
func NewClient() (*k8s.Clientset, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}
	return k8s.NewForConfig(config)
}
