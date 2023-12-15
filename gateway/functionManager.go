package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"sync"
	"time"
)

const LabelName = "faas_function"

var functionManager *FunctionManager

type Function struct {
	Name       string
	DeployedAt Node
	PodName    string
	IP         string
	UpdatedAt  *time.Time
}

// new funciton
func NewFunction(name string, podName string) *Function {
	return &Function{
		Name:      name,
		PodName:   podName,
		UpdatedAt: nil,
	}
}

// k8s node
type Node struct {
	Name string
	UID  string
}

type FunctionManager struct {
	clientset *kubernetes.Clientset
	mu        sync.Mutex

	FunctionMap map[string]map[string]*Function // map[functionName]map[podName]*Function
}

func InitFunctionManager() {
	functionManager = NewFunctionManager()
	fetchFunctionInfo()
}

func fetchFunctionInfo() {
	fmt.Println("start to get pod details")

	//allUpdated := true
	//now := time.Now()
	//for name, functions := range functionManager.FunctionMap {
	//for _, function := range functions {
	//threshold := now.Add(-1 * time.Minute)

	//fmt.Println(function.UpdatedAt.Before(threshold))
	//if function.UpdatedAt.Before(threshold) {
	//	allUpdated = false
	//}
	//	}
	//}

	//if allUpdated {
	//	fmt.Println("All updated")
	//	continue
	//}

	pods, err := functionManager.clientset.CoreV1().Pods("openfaas-fn").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		functionName := pod.Labels[LabelName]

		if _, ok := functionManager.FunctionMap[functionName]; !ok {
			functionManager.FunctionMap[functionName] = make(map[string]*Function)
		}

		if _, ok := functionManager.FunctionMap[functionName][pod.Name]; !ok {
			functionManager.FunctionMap[functionName][pod.Name] = NewFunction(functionName, pod.Name)
		}
		f := functionManager.FunctionMap[functionName][pod.Name]
		f.IP = pod.Status.PodIP
		now := time.Now()
		f.UpdatedAt = &now
		f.PodName = pod.Name

		f.DeployedAt = Node{
			Name: pod.Spec.NodeName,
			UID:  string(pod.UID),
		}
		f.UpdatedAt = &now

		log.Printf("Pod name: %s\n", pod.Name)
		log.Printf("Pod IP: %s\n", pod.Status.PodIP)

	}
}

func NewFunctionManager() *FunctionManager {
	f := &FunctionManager{
		FunctionMap: make(map[string]map[string]*Function),
	}

	var err error
	config := getConfig(false)
	f.clientset, err = kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (f *FunctionManager) GetRandomPodFunction(name string) *Function {
	if podMap, ok := functionManager.FunctionMap[name]; !ok {
		return nil
	} else if len(podMap) == 0 {
		return nil
	} else {
		for _, function := range podMap {
			return function
		}
		return nil
	}
}

func (f *FunctionManager) RegisterFunction(name string, podName string) {
	if _, ok := functionManager.FunctionMap[name]; !ok {
		functionManager.FunctionMap[name] = make(map[string]*Function)
	}

	if _, ok := functionManager.FunctionMap[name][podName]; !ok {
		functionManager.FunctionMap[name][podName] = NewFunction(name, podName)
	}
}

var config *rest.Config

func getConfig(inCluster bool) *rest.Config {

	var err error
	if config != nil {
		return config
	}

	if inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
		return config
	}

	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}
