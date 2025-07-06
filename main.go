package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	"github.com/t1bur1an/k8s-pod-ttl-killer/k8s"
	"k8s.io/client-go/kubernetes"
)

func main() {
	envConfig := config.ReadConfig()
	kubeconfig, err := k8s.GetKubeConfig()
	if err != nil {
		slog.Error("Error to get kube config")
		os.Exit(2)
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

	// add prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

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
