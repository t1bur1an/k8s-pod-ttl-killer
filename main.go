package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	"github.com/t1bur1an/k8s-pod-ttl-killer/k8s"
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

	go k8s.CheckClusterPodsPoll(clientset)

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	httpListenString := fmt.Sprintf(
		"%s:%s",
		envConfig.HTTPListenAddress,
		envConfig.HTTPListenPort)
	slog.Info(
		"Starting http server",
		"address", envConfig.HTTPListenAddress,
		"port", envConfig.HTTPListenPort)
	err = http.ListenAndServe(httpListenString, nil)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
}
