package main

import (
	"func/workload"
	"log"
)

func main() {

	log.Println("hello world")

	ch := make(chan bool)
	go func() {
		log.Println("New workload")
		workload := workload.NewBasicWorkload()
		workload.Run()
		workload.ReportOutCome()
	}()
	<-ch
}
