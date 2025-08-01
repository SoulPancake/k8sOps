package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"time"

	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	flag.Parse()
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kube.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		pvcs, _ := clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), meta.ListOptions{
			LabelSelector: "auto-resize=true",
		})

		for _, pvc := range pvcs.Items {
			// Simulate checking usage (in real case, you'd get metrics from cAdvisor/Prometheus)
			usagePercent := simulateStorageUsage(pvc.Name)

			if usagePercent > 80 {
				fmt.Printf("PVC %s/%s is %d%% full. Resizing...\n", pvc.Namespace, pvc.Name, usagePercent)
				newSize := "2Gi" // hardcoded for simplicity

				pvc.Spec.Resources.Requests[v1.ResourceStorage] = resourceMustParse(newSize)
				_, err := clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Update(context.TODO(), &pvc, meta.UpdateOptions{})
				if err != nil {
					fmt.Println("Update failed:", err)
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func simulateStorageUsage(name string) int {
	// Simulate increasing usage per call
	return time.Now().Second() % 100
}

func resourceMustParse(str string) resource.Quantity {
	q, err := resource.ParseQuantity(str)
	if err != nil {
		panic(err)
	}
	return q
}
