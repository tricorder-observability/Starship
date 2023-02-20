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

	"github.com/tricorder/src/utils/log"

	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"
)

var quitCh = make(chan struct{})

type ResourceWatcher struct {
	eg        errgroup.Group
	clientset kubernetes.Interface
	pgClient  *pg.Client
}

func NewResourceWatcher(clientset kubernetes.Interface, pgClient *pg.Client) *ResourceWatcher {
	watcher := new(ResourceWatcher)
	watcher.clientset = clientset
	watcher.pgClient = pgClient
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
		upsert(w.pgClient, "nodes", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Nodes().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if ok {
				upsert(w.pgClient, "nodes", node, node.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			node, ok := newObj.(*corev1.Node)
			if ok {
				upsert(w.pgClient, "nodes", node, node.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			node, ok := obj.(*corev1.Node)
			if ok {
				deleteByID(w.pgClient, "nodes", node.UID)
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
		upsert(w.pgClient, "namespaces", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Namespaces().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if ok {
				upsert(w.pgClient, "namespaces", ns, ns.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ns, ok := newObj.(*corev1.Namespace)
			if ok {
				upsert(w.pgClient, "namespaces", ns, ns.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if ok {
				deleteByID(w.pgClient, "namespaces", ns.UID)
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
		upsert(w.pgClient, "pods", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if ok {
				upsert(w.pgClient, "pods", pod, pod.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod, ok := newObj.(*corev1.Pod)
			if ok {
				upsert(w.pgClient, "pods", pod, pod.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if ok {
				deleteByID(w.pgClient, "pods", pod.UID)
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
		upsert(w.pgClient, "endpoints", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Endpoints().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ep, ok := obj.(*corev1.Endpoints)
			if ok {
				upsert(w.pgClient, "endpoints", ep, ep.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ep, ok := newObj.(*corev1.Endpoints)
			if ok {
				upsert(w.pgClient, "endpoints", ep, ep.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ep, ok := obj.(*corev1.Endpoints)
			if ok {
				deleteByID(w.pgClient, "endpoints", ep.UID)
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
		upsert(w.pgClient, "services", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Core().V1().Services().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc, ok := obj.(*corev1.Service)
			if ok {
				upsert(w.pgClient, "services", svc, svc.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			svc, ok := newObj.(*corev1.Service)
			if ok {
				upsert(w.pgClient, "services", svc, svc.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			svc, ok := obj.(*corev1.Service)
			if ok {
				deleteByID(w.pgClient, "services", svc.UID)
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
		upsert(w.pgClient, "replicasets", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Apps().V1().ReplicaSets().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			rs, ok := obj.(*appsv1.ReplicaSet)
			if ok {
				upsert(w.pgClient, "replicasets", rs, rs.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			rs, ok := newObj.(*appsv1.ReplicaSet)
			if ok {
				upsert(w.pgClient, "replicasets", rs, rs.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			rs, ok := obj.(*appsv1.ReplicaSet)
			if ok {
				deleteByID(w.pgClient, "replicasets", rs.UID)
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
		upsert(w.pgClient, "deployments", &obj, obj.UID)
	}

	factory := informers.NewSharedInformerFactory(w.clientset, 12*time.Hour)
	informer := factory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment, ok := obj.(*appsv1.Deployment)
			if ok {
				upsert(w.pgClient, "deployments", deployment, deployment.UID)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deployment, ok := newObj.(*appsv1.Deployment)
			if ok {
				upsert(w.pgClient, "deployments", deployment, deployment.UID)
			}
		},
		DeleteFunc: func(obj interface{}) {
			deployment, ok := obj.(*appsv1.Deployment)
			if ok {
				deleteByID(w.pgClient, "deployments", deployment.UID)
			}
		},
	})

	informer.Run(quitCh)
	return nil
}
