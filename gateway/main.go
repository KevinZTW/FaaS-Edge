package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("init\n")
	InitFunctionManager()

	router := httprouter.New()
	router.POST("/filter", K8SSchedulerFilter)

	// TODO: Add the ingress to the OpenFaaS gateway to handle large number of requests
	router.GET("/test-auto-scaling", TestAutoScaling)

	//TODO: request through this gateway

	log.Fatal(http.ListenAndServe(":8888", router))
}
