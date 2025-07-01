package main

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	"github.com/t1bur1an/k8s-pod-ttl-killer/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config.ReadConfig()
	var kubeconfig string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
		if utils.DeletePodCheck(pod) {
			slog.Info("Found pod to delete %s namespace %s", pod.GetName(), pod.Namespace)
			podContext := context.Background()
			utils.DeletePod(clientset, pod, podContext)
		}
	}
}
