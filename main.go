package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/t1bur1an/k8s-pod-ttl-killer/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		podReadyTimestamp, err := utils.GetPodReadyTime(pod)
		if err != nil {
			fmt.Printf("Pod %s not in ready state\n", pod.GetName())
		} else {
			ttl, err := utils.FilterAnnotations(pod.Annotations, "pod-killer-ttl")
			if err != nil {
				fmt.Printf("Pod %s got an error with annotations: %s\n", pod.GetName(), err.Error())
			} else {
				fmt.Printf("Pod %s pod time info, timestamp: %d, ttl: %d\n", pod.GetName(), podReadyTimestamp, ttl)
			}
		}
	}
}
