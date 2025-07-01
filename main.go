package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/t1bur1an/k8s-pod-ttl-killer/utils"
	corev1 "k8s.io/api/core/v1"
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
	for i, pod := range pods.Items {
		podReadyTimestamp, err := utils.GetPodReadyTime(pod)
		if err != nil {
			fmt.Printf("Pod %s not in ready state", pod.Spec.Containers[0].Name)
		}
		for _, podCondition := range pod.Status.Conditions {
			if podCondition.Type == corev1.PodReady {
				ttl, err := utils.FilterAnnotations(pod.Annotations, "pod-killer-ttl")
				if err != nil {
					fmt.Printf("Got an error with annotations: %s", err.Error())
				}
				fmt.Printf("Pod %d pod time info, timestamp: %d, ttl: %d\n", i, podReadyTimestamp, ttl)
			}
		}
	}
}
