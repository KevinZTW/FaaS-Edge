package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
	"log"
)

func filter(args extender.ExtenderArgs) *extender.ExtenderFilterResult {
	fmt.Print("filter call\n")
	var filteredNodes []v1.Node
	failedNodes := make(extender.FailedNodesMap)
	pod := args.Pod

	//TODO: hard code the label as name
	log.Println("register function: " + args.Pod.ObjectMeta.Labels[LabelName] + " pod as: " + pod.Name)

	functionManager.RegisterFunction(args.Pod.ObjectMeta.Labels[LabelName], pod.Name)

	result := extender.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}
