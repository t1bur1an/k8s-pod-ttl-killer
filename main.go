package main

import (
	"context"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	"github.com/t1bur1an/k8s-pod-ttl-killer/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	envConfig := config.ReadConfig()
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
	for {
		slog.Info("Check cluster pods: started")
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			if utils.DeletePodCheck(pod) {
				slog.Info("Pod to delete", "pod", pod.GetName(), "namespace", pod.Namespace)
				podContext := context.Background()
				utils.DeletePod(clientset, pod, podContext)
			}
		}
		slog.Info("Check cluster pods: done", "clusterPodsCount", len(pods.Items))
		time.Sleep(time.Duration(envConfig.CheckIntervalSeconds)*time.Second)
	}
}
