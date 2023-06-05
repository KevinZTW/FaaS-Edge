package function

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func GetPodDetails() (IP string, name string) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err.Error(), ""
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err.Error(), ""
	}

	IP = ""
	retry := 0
	for {
		pods, err := clientset.CoreV1().Pods("openfaas-fn").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			// log.Printf("Pod name: %s\n", pod.Name)
			// log.Printf("Pod IP: %s\n", pod.Status.PodIP)
			// log.Printf("Pod name: %v equals to hostname: %v => %v\n", pod.Name, os.Getenv("HOSTNAME"), pod.Name == os.Getenv("HOSTNAME"))
			if pod.Name == os.Getenv("HOSTNAME") {
				IP = pod.Status.PodIP
			}
		}
		if IP != "" || retry > 2 {
			break
		}
		retry++
		time.Sleep(1 * time.Second)
	}

	name = os.Getenv("HOSTNAME")
	return IP, name
}

// Handle a serverless request
func Handle(req []byte) string {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A new newwalrus appears")

	// connect to host cache-controller

	ip, name := GetPodDetails()
	log.Info("Sending request to cache-controller with IP: " + ip + " and name: " + name)

	resp, err := http.Get("http://10.101.201.15:3037/register?addr=" + ip + "&name=" + name)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)

	return fmt.Sprintf("Hello, this is kelly v5. You said: %s", string(req))
}
