package main

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

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func main() {
	for {
		time.Sleep(1 * time.Second)
	}
}

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

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
	for {
		if IP != "" {
			break
		} else {
			log.Printf("No IP for now.\n")
		}

		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			log.Printf("Pod name: %s\n", pod.Name)
			log.Printf("Pod IP: %s\n", pod.Status.PodIP)

			pod, _ := clientset.CoreV1().Pods("").Get(context.TODO(), pod.Name, metav1.GetOptions{})
			if pod.Name == os.Getenv("HOSTNAME") {
				IP = pod.Status.PodIP
			}
		}

		log.Printf("Waits...\n")
		time.Sleep(1 * time.Second)
	}

	name = os.Getenv("HOSTNAME")
	log.Printf("Trying os.Getenv(\"HOSTNAME/IP\"): [%s][%s]\n", name, IP)

	return IP, name
}

// Handle a serverless request
func Handle(req []byte) string {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A new newwalrus appears")

	// connect to host cache-controller

	ip, name := GetPodDetails()
	fmt.Println(ip, name)
	// ip, name := "test-ip", "test-name"

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

	return fmt.Sprintf("Hello, this is kelly v4. You said: %s", string(req))
}
