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

	log.Fatal(http.ListenAndServe(":8888", router))
}
