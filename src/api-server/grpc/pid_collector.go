package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/tricorder/src/utils/log"

	pb "github.com/tricorder/src/api-server/pb"
	"github.com/tricorder/src/utils/pg"
	"github.com/tricorder/src/utils/retry"
)

const (
	procInfoTableName = "process_info"
)

// idPath is UUID location in processInfo table based on Postgres json path expression
var idPath = []string{"container", "id"}

// PIDCollector implements the ProcessCollector gRPC service.
type PIDCollector struct {
	clientset kubernetes.Interface
	pgClient  *pg.Client
	// pod informer stops when this channel is closed.
	// https://pkg.go.dev/k8s.io/client-go@v0.26.1/tools/cache#SharedInformer
	podInformerQuitChan chan struct{}
}

func NewPIDCollector(clientset kubernetes.Interface, pgClient *pg.Client) *PIDCollector {
	c := new(PIDCollector)
	c.clientset = clientset
	c.pgClient = pgClient
	c.podInformerQuitChan = make(chan struct{})

	err := retry.ExpBackOffWithLimit(func() error {
		return pgClient.CreateTable(pg.GetJSONBTableSchema(procInfoTableName))
	})
	if err != nil {
		log.Fatalf("While preparing to start process info collector , failed to create table, error: %v", err)
	}

	return c
}

// ReportProcess implements the gRPC method ReportProcess of the ProcessCollector service.
func (s *PIDCollector) ReportProcess(stream pb.ProcessCollector_ReportProcessServer) error {
	for {
		pw, err := stream.Recv()
		if status.Code(err) == codes.Unavailable {
			log.Warnf("Agent disconnected, error: %v", err)
			// The informer will be stopped when stopCh is closed.
			close(s.podInformerQuitChan)
			return fmt.Errorf("streaming connection with agent is broken (agent died?), error: %v", err)
		}

		switch msg := pw.Msg.(type) {
		case *pb.ProcessWrapper_NodeName:
			go func() {
				err := s.podWatch(msg.NodeName, stream)
				if err != nil {
					log.Errorf("While starting watching pods on node '%s', failed to start, error: %v", msg.NodeName, err)
				}
			}()
		case *pb.ProcessWrapper_Process:
			// Receive process info update from agents, and write the info as JSON blob into data table process_info
			value, _ := protojson.Marshal(msg.Process)
			if err = s.pgClient.JSON().Upsert(procInfoTableName, msg.Process.Container.Id, value, idPath...); err != nil {
				log.Errorf("While reporting process info, failed to Upsert, error: %v", err)
			}
		default:
			log.Errorln("pw.Msg must be one of NodeName, ProcessInfo")
		}
	}
}

// podWatch include watch pods with node name, and handle the event(onAdd/onUpdate/onDelete)
func (s *PIDCollector) podWatch(nodeName string, stream pb.ProcessCollector_ReportProcessServer) error {
	// List pods with node name
	list, err := s.clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		log.Errorf("col.clientset.CoreV1().Pods().List error %s", err)
		return err
	}
	for _, pod := range list.Items {
		pushContainerInfoToAgent(&pod, stream)
	}

	// Watch pods with node name
	podInformer := informers.NewSharedInformerFactoryWithOptions(s.clientset, 12*time.Hour,
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fmt.Sprintf("spec.nodeName=%s", nodeName)
		})).Core().V1().Pods().Informer()

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if ok {
				pushContainerInfoToAgent(pod, stream)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod, ok := newObj.(*corev1.Pod)
			if ok {
				pushContainerInfoToAgent(pod, stream)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if ok {
				for _, c := range pod.Status.ContainerStatuses {
					if err := s.pgClient.JSON().Delete(procInfoTableName, c.ContainerID, idPath...); err != nil {
						log.Errorf("While watching Pod update, failed to delete container, error: %v", err)
					}
				}
			}
		},
	})

	podInformer.Run(s.podInformerQuitChan)
	return nil
}

func pushContainerInfoToAgent(pod *corev1.Pod, stream pb.ProcessCollector_ReportProcessServer) {
	if pod.Status.Phase != corev1.PodRunning {
		log.Warnf("pod %s must be running: %s", pod.Name, pod.Status.Phase)
		return
	}
	for _, container := range pod.Status.ContainerStatuses {
		ci := &pb.ContainerInfo{
			Id:       container.ContainerID,
			Name:     container.Name,
			PodUid:   string(pod.UID),
			PodName:  pod.Name,
			QosClass: string(pod.Status.QOSClass),
		}
		if err := stream.Send(ci); err != nil {
			log.Errorf("pushContainerInfoToAgent error %v", err)
		}
	}
}
