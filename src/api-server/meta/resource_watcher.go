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

package meta

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/tricorder/src/api-server/http/dao"
	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/utils/cond"
	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"
)

var quitCh = make(chan struct{})

type ResourceWatcher struct {
	eg        errgroup.Group
	clientset kubernetes.Interface
	pgClient  *pg.Client
	nodeAgent *dao.NodeAgentDao
	waitCond  *cond.Cond
}

func NewResourceWatcher(clientset kubernetes.Interface, pgClient *pg.Client,
	nodeAgent *dao.NodeAgentDao, waitCond *cond.Cond,
) *ResourceWatcher {
	watcher := new(ResourceWatcher)
	watcher.clientset = clientset
	watcher.pgClient = pgClient
	watcher.nodeAgent = nodeAgent
	watcher.waitCond = waitCond
	err := retry.ExpBackOffWithLimit(func() error {
		return initResourceTables(pgClient)
	})
	if err != nil {
		log.Fatalf("While preparing to start ResourceWatcher , failed to initializing resource tables, error: %v", err)
	}

	return watcher
}

// StartWatching launches all of the Go routine to watch for updates for each and every type of objects.
func (w *ResourceWatcher) StartWatching() error {
	w.eg.Go(w.node)
	w.eg.Go(w.namespace)
	w.eg.Go(w.pod)
	w.eg.Go(w.endpoints)
	w.eg.Go(w.service)
	w.eg.Go(w.replicaSet)
	w.eg.Go(w.deployment)

	return w.eg.Wait()
}

func (w *ResourceWatcher) node() error {
	list, err := w.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.CoreV1().Nodes().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, NodeTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Nodes().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if ok {
				upsert(w.pgClient, NodeTable, node, node.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			node, ok := newObj.(*corev1.Node)
			if ok {
				upsert(w.pgClient, NodeTable, node, node.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if ok {
				deleteByID(w.pgClient, NodeTable, node.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}

func (w *ResourceWatcher) namespace() error {
	list, err := w.clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.CoreV1().Namespaces().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, NameSpaceTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Namespaces().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if ok {
				upsert(w.pgClient, NameSpaceTable, ns, ns.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ns, ok := newObj.(*corev1.Namespace)
			if ok {
				upsert(w.pgClient, NameSpaceTable, ns, ns.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if ok {
				deleteByID(w.pgClient, NameSpaceTable, ns.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}

func (w *ResourceWatcher) pod() error {
	list, err := w.clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.CoreV1().Pods().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, PodTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if ok {
				upsert(w.pgClient, PodTable, pod, pod.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod, ok := newObj.(*corev1.Pod)
			if ok {
				upsert(w.pgClient, PodTable, pod, pod.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if ok {
				deleteByID(w.pgClient, PodTable, pod.UID)
				if na, err := w.nodeAgent.QueryByPodID(string(pod.UID)); err == nil {
					if err = w.nodeAgent.UpdateStateByID(na.AgentID, int(pb.AgentState_TERMINATED)); err != nil {
						log.Errorf("while deleting pod, failed to nodeAgent UpdateStateByID %s, error %s", na.AgentID, err)
					} else {
						// TODO(yzhao): We seem needing a more structured way of communicating needed actions between API Server's
						// HTTP and gRPC components. See https://github.com/tricorder-observability/starship/issues/150.
						w.waitCond.Broadcast()
					}
				}
			}
		},
	})

	informer.Run(quitCh)
	return nil
}

func (w *ResourceWatcher) endpoints() error {
	list, err := w.clientset.CoreV1().Endpoints(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.CoreV1().Endpoints().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, EndPointTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Endpoints().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ep, ok := obj.(*corev1.Endpoints)
			if ok {
				upsert(w.pgClient, EndPointTable, ep, ep.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ep, ok := newObj.(*corev1.Endpoints)
			if ok {
				upsert(w.pgClient, EndPointTable, ep, ep.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ep, ok := obj.(*corev1.Endpoints)
			if ok {
				deleteByID(w.pgClient, EndPointTable, ep.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}

func (w *ResourceWatcher) service() error {
	list, err := w.clientset.CoreV1().Services(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.CoreV1().Services().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, ServiceTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Services().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc, ok := obj.(*corev1.Service)
			if ok {
				upsert(w.pgClient, ServiceTable, svc, svc.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			svc, ok := newObj.(*corev1.Service)
			if ok {
				upsert(w.pgClient, ServiceTable, svc, svc.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			svc, ok := obj.(*corev1.Service)
			if ok {
				deleteByID(w.pgClient, ServiceTable, svc.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}

func (w *ResourceWatcher) replicaSet() error {
	list, err := w.clientset.AppsV1().ReplicaSets(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.AppsV1().ReplicaSets().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, ReplicSetTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Apps().V1().ReplicaSets().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			rs, ok := obj.(*appsv1.ReplicaSet)
			if ok {
				upsert(w.pgClient, ReplicSetTable, rs, rs.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			rs, ok := newObj.(*appsv1.ReplicaSet)
			if ok {
				upsert(w.pgClient, ReplicSetTable, rs, rs.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			rs, ok := obj.(*appsv1.ReplicaSet)
			if ok {
				deleteByID(w.pgClient, ReplicSetTable, rs.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}

func (w *ResourceWatcher) deployment() error {
	list, err := w.clientset.AppsV1().Deployments(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("clientset.AppsV1().Deployments().List error %s", err)
		return err
	}
	for _, obj := range list.Items {
		upsert(w.pgClient, DeploymentTable, &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment, ok := obj.(*appsv1.Deployment)
			if ok {
				upsert(w.pgClient, DeploymentTable, deployment, deployment.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deployment, ok := newObj.(*appsv1.Deployment)
			if ok {
				upsert(w.pgClient, DeploymentTable, deployment, deployment.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			deployment, ok := obj.(*appsv1.Deployment)
			if ok {
				deleteByID(w.pgClient, DeploymentTable, deployment.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}
