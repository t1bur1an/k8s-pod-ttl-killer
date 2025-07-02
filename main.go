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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	envConfig := config.ReadConfig()
	slog.Info("Try to read InClusterConfig")
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		slog.Error("InClusterConfig", "error", err.Error())
	}
	if kubeconfig == nil {
		slog.Info("Try to read kube config")
		home := homedir.HomeDir()
		if home != "" {
			kubeconfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
			if err != nil {
				panic(err.Error())
			}
		}
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
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
