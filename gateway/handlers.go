package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
)

func K8SSchedulerFilterHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)

	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Error: err.Error(),
		}
	} else {
		b, err := json.MarshalIndent(extenderArgs, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(string(b))
		extenderFilterResult = filter(extenderArgs)
	}

	if response, err := json.Marshal(extenderFilterResult); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func TestAutoScalingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	for i := 0; i < 10000; i++ {
		// send request to http://localhost:8080/function/nodeinfo
		go func() {
			res, err := http.Get("http://localhost:8080/function/nodeinfo")
			if err != nil {
				panic(err)
			}

			// read response body
			defer res.Body.Close()
			//body, err := ioutil.ReadAll(res.Body)
			//if err != nil {
			//	panic(err)
			//}
			//fmt.Println(string(body))
		}()

	}
}

func FunctionRequestHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	functionName := params.ByName("name")
	fmt.Println("function name: " + functionName)

	function := functionManager.GetRandomPodFunction(functionName)
	if function == nil {
		writer.Write([]byte("Not found function with name " + functionName + "\n"))
		return
	}

	res, err := http.Get("http://" + function.IP + ":8080")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	writer.Write(body)
}
