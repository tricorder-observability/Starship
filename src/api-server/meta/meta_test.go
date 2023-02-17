package meta

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/tricorder/src/testing/pg"
)

var testAgainstK8s = flag.Bool("test_against_k8s", false, "If true, test against Kubernetes pointed to "+
	"by the local Kube config")

func init() {
	// This is needed to define the default test flags.
	testing.Init()
	flag.Parse()
}

// Tests that json.Marshal can handle golang object.
// Tests that runtime.Encode() with jsonserializer produces the same results
// as json.Marshal().
func TestCompareRuntimeEncodeVsJsonMarshal(t *testing.T) {
	assert := assert.New(t)

	pod := &corev1.Pod{}
	pod.Name = "pod1"
	pod.UID = types.UID("aaa")
	pod.CreationTimestamp = metav1.Now()

	encoder := jsonserializer.NewSerializerWithOptions(nil, nil, nil, jsonserializer.SerializerOptions{})
	encoded, err := runtime.Encode(encoder, pod)
	assert.Nil(err)

	encoded2, err := json.Marshal(pod)
	assert.Nil(err)

	assert.Equal(string(encoded), string(encoded2)+"\n")
}

// Tests that resources(Namespace|Pod|Enpoints|Service|ReplicaSet|Deployment) are the same between K8s and DB.
// MUST set arg 'K8s', otherwise testing will skip.
// bazel test :meta_test --test_arg=--test_against_k8s=true
func TestCompareResourcesBetweenK8sAndDB(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var clientset kubernetes.Interface
	var err error

	if *testAgainstK8s {
		clientset, err = kubernetes.NewForConfig(ctrl.GetConfigOrDie())
		assert.Nil(err)
	} else {
		clientset = fake.NewSimpleClientset()
	}

	cleaner, pgClient, err := pg.LaunchContainer()
	require.Nil(err)
	defer func() {
		assert.Nil(cleaner())
	}()

	watcher := NewResourceWatcher(clientset, pgClient)
	go func() {
		err = watcher.StartWatching()
		assert.Nil(err)
	}()

	// Create a test namespace and clean it after testing
	testNS := &corev1.Namespace{}
	testNS.Name = "ns1"
	testNS.UID = types.UID("ns1")
	defer func() {
		assert.Nil(clientset.CoreV1().Namespaces().Delete(context.TODO(), testNS.Name, metav1.DeleteOptions{}))
	}()

	// Test namespace correctness between K8s and DB
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), testNS, metav1.CreateOptions{})
	assert.Nil(err)
	// Ensure to be created completely
	time.Sleep(2 * time.Second)
	// Query from K8s
	nsInK8s, err := clientset.CoreV1().Namespaces().Get(context.TODO(), testNS.Name, metav1.GetOptions{})
	assert.Nil(err)
	// Query from DB
	nsInDB := corev1.Namespace{}
	err = pgClient.JSON().Get("namespaces", &nsInDB, fmt.Sprintf("WHERE data#>>'{metadata,name}'='%s'", testNS.Name))
	assert.Nil(err)
	// Compare K8s and DB for correctness
	assert.Equal(nsInK8s.Name, nsInDB.Name)
	assert.Equal(nsInK8s.Status, nsInDB.Status)

	// Test pod correctness between K8s and DB
	pod1 := &corev1.Pod{}
	pod1.Name = "pod1"
	pod1.UID = types.UID("uid1")
	container1 := corev1.Container{Name: "test", Image: "xxx"}
	pod1.Spec.Containers = []corev1.Container{container1}
	_, err = clientset.CoreV1().Pods(testNS.Name).Create(context.TODO(), pod1, metav1.CreateOptions{})
	assert.Nil(err)
	// Ensure to be created completely
	time.Sleep(2 * time.Second)
	// Query from K8s
	podInK8s, err := clientset.CoreV1().Pods(testNS.Name).Get(context.TODO(), pod1.Name, metav1.GetOptions{})
	// Query from DB
	podInDB := corev1.Pod{}
	err = pgClient.JSON().Get("pods", &podInDB, fmt.Sprintf("WHERE data#>>'{metadata,name}'='%s'", pod1.Name))
	assert.Nil(err)
	// Compare K8s and DB for correctness
	assert.Equal(podInK8s.Name, podInDB.Name)
	assert.Equal(podInK8s.Status, podInDB.Status)

	// Test endpoints correctness between K8s and DB
	endpoints := &corev1.Endpoints{}
	endpoints.Name = "endpoints1"
	endpoints.UID = types.UID("endpoints_uid1")
	_, err = clientset.CoreV1().Endpoints(testNS.Name).Create(context.TODO(), endpoints, metav1.CreateOptions{})
	assert.Nil(err)
	// Ensure to be created completely
	time.Sleep(2 * time.Second)
	// Query from K8s
	epInK8s, err := clientset.CoreV1().Endpoints(testNS.Name).Get(context.TODO(), endpoints.Name, metav1.GetOptions{})
	assert.Nil(err)
	// Query from DB
	epInDB := corev1.Endpoints{}
	err = pgClient.JSON().Get("endpoints", &epInDB, fmt.Sprintf("WHERE data#>>'{metadata,name}'='%s'", endpoints.Name))
	assert.Nil(err)
	// Compare K8s and DB for correctness
	assert.Equal(epInK8s.Name, epInDB.Name)
	assert.Equal(epInK8s.String(), epInDB.String())

	// Test service watched and inserted into DB
	service := &corev1.Service{}
	service.Name = "service1"
	service.UID = types.UID("service_uid1")
	sp := corev1.ServicePort{Name: "sp1", Port: 80}
	service.Spec.Ports = []corev1.ServicePort{sp}
	_, err = clientset.CoreV1().Services(testNS.Name).Create(context.TODO(), service, metav1.CreateOptions{})
	assert.Nil(err)
	// Ensure to be created completely
	time.Sleep(2 * time.Second)
	// Query from K8s
	svcInK8s, err := clientset.CoreV1().Services(testNS.Name).Get(context.TODO(), service.Name, metav1.GetOptions{})
	assert.Nil(err)
	// Query from DB
	svcInDB := corev1.Service{}
	err = pgClient.JSON().Get("services", &svcInDB, fmt.Sprintf("WHERE data#>>'{metadata,name}'='%s'", service.Name))
	assert.Nil(err)
	// Compare K8s and DB for correctness
	assert.Equal(svcInK8s.Name, svcInDB.Name)
	assert.Equal(svcInK8s.Status, svcInDB.Status)

	// Test replicaSet watched and inserted into DB
	replicaSet := &appsv1.ReplicaSet{}
	replicaSet.Name = "replicaset1"
	replicaSet.UID = types.UID("replicaset_uid1")
	replicaSet.Spec.Template.Spec.Containers = []corev1.Container{container1}
	label := map[string]string{"a": "b"}
	replicaSet.Spec.Template.Labels = label
	replicaSet.Spec.Selector = &metav1.LabelSelector{MatchLabels: label}

	_, err = clientset.AppsV1().ReplicaSets(testNS.Name).Create(context.TODO(), replicaSet, metav1.CreateOptions{})
	assert.Nil(err)
	// Ensure to be created completely
	time.Sleep(2 * time.Second)
	// Query from K8s
	rsInK8s, err := clientset.AppsV1().ReplicaSets(testNS.Name).Get(context.TODO(), replicaSet.Name, metav1.GetOptions{})
	assert.Nil(err)
	// Query from DB
	rsInDB := appsv1.ReplicaSet{}
	err = pgClient.JSON().Get("replicasets", &rsInDB, fmt.Sprintf("WHERE data#>>'{metadata,name}'='%s'", replicaSet.Name))
	assert.Nil(err)
	// Compare K8s and DB for correctness
	assert.Equal(rsInK8s.Name, rsInDB.Name)
	assert.Equal(rsInK8s.Status, rsInDB.Status)

	// Test deployment watched and inserted into DB
	deployment := &appsv1.Deployment{}
	deployment.Name = "deployment1"
	deployment.UID = types.UID("deployment_uid1")
	deployment.Spec.Template.Spec.Containers = []corev1.Container{container1}
	deployment.Spec.Template.Labels = label
	deployment.Spec.Selector = &metav1.LabelSelector{MatchLabels: label}
	_, err = clientset.AppsV1().Deployments(testNS.Name).Create(context.TODO(), deployment, metav1.CreateOptions{})
	assert.Nil(err)
	// Ensure be created
	time.Sleep(2 * time.Second)
	// Query from K8s
	deploymentInK8s, err := clientset.AppsV1().Deployments(testNS.Name).Get(
		context.TODO(),
		deployment.Name,
		metav1.GetOptions{})
	assert.Nil(err)
	// Query from DB
	deploymentInDB := appsv1.Deployment{}
	err = pgClient.JSON().Get("deployments", &deploymentInDB,
		fmt.Sprintf("WHERE data#>>'{metadata,name}'='%s'", deployment.Name))
	assert.Nil(err)
	// Compare K8s and DB for correctness
	assert.Equal(deploymentInK8s.Name, deploymentInDB.Name)
	assert.Equal(deploymentInK8s.Status, deploymentInDB.Status)
}

// TestInitResourceTable test data still exists after initResourceTable again
func TestInitResourceTable(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	cleaner, pgClient, err := pg.LaunchContainer()
	require.Nil(err)
	defer func() {
		assert.Nil(cleaner())
	}()
	assert.Nil(initResourceTables(pgClient))

	pod1 := &corev1.Pod{}
	pod1.Name = "pod1"
	pod1.UID = types.UID("uid1")
	value, _ := json.Marshal(pod1)
	require.Nil(pgClient.JSON().Upsert(PodTable, string(pod1.UID), value))

	assert.Nil(initResourceTables(pgClient))
	pod := &corev1.Pod{}
	require.Nil(pgClient.JSON().Get(PodTable, pod))
	assert.Equal(pod1.Name, pod.Name)
}
