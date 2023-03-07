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
func getService(client k8s.Interface, namespace, serviceName string) (*v1.Service, error) {
	return client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, matav1.GetOptions{})
}

// GetStarshipService returns the starship Service by name and namespace
func getStarshipService(client k8s.Interface) (*v1.Service, error) {
	return getService(client, DefaultNamespace, DefaultServiceName)
}

// getStarshipServiceURL returns the starship Service ClusterIP
func getStarshipServiceURL(client k8s.Interface) (string, error) {
	service, err := getStarshipService(client)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}

// getStarshipServicePort returns the starship Service serverhttp port
func getStarshipServicePort(client k8s.Interface) (int32, error) {
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
	apiServerAddress := fmt.Sprintf("%s:%d", ip, port)
	return apiServerAddress, nil
}
