package k8s

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"

	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetKubeConfig() (*rest.Config, error) {
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
	return kubeconfig, nil
}

func CheckClusterPodsPoll(clientset *kubernetes.Clientset) {
	envConfig := config.ReadConfig()
	for {
		slog.Info("Check cluster pods: started")
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			if DeletePodCheck(pod) {
				slog.Info("Pod to delete", "pod", pod.GetName(), "namespace", pod.Namespace)
				podContext := context.Background()
				DeletePod(clientset, pod, podContext)
			}
		}
		slog.Info("Check cluster pods: done", "clusterPodsCount", len(pods.Items))
		time.Sleep(time.Duration(envConfig.CheckIntervalSeconds) * time.Second)
	}
}

func FilterAnnotations(filterAnnotations map[string]string, targetAnnotation string) (int64, error) {
	outputDuration := int64(0)
	for annotation, value := range filterAnnotations {
		if annotation == targetAnnotation {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return 0, err
			}
			outputDuration = int64(intValue)
			return outputDuration, nil
		}
	}
	errMsg := fmt.Sprintf("No %s annotation was found", targetAnnotation)
	return outputDuration, errors.New(errMsg)
}

func GetPodReadyTimestamp(pod corev1.Pod) (int64, error) {
	for _, podCondition := range pod.Status.Conditions {
		if podCondition.Type == corev1.PodReady {
			return podCondition.LastTransitionTime.Unix(), nil
		}
	}
	return 0, errors.New("Pod is not ready")
}

func DeletePodCheck(pod corev1.Pod) bool {
	envConfig := config.ReadConfig()
	ttl, err := FilterAnnotations(pod.Annotations, envConfig.TTLAnnotation)
	if err != nil {
		return false
	}
	podReadyTimestamp, err := GetPodReadyTimestamp(pod)
	if err != nil {
		slog.Info("Not in ready state", "pod", pod.GetName())
		return false
	}
	utcTimeNow := time.Now().UTC().Unix()
	slog.Info(
		"Pod:",
		"pod", pod.GetName(),
		"annotation", envConfig.TTLAnnotation,
		"value", ttl,
		"podReadyTimestamp",
		podReadyTimestamp,
		"ttl", ttl,
		"utcTimeNow", utcTimeNow)

	return (podReadyTimestamp + ttl) <= utcTimeNow
}

func DeletePod(clientset *kubernetes.Clientset, pod corev1.Pod, podContext context.Context) {
	err := clientset.CoreV1().Pods(pod.Namespace).Delete(podContext, pod.GetName(), *metav1.NewDeleteOptions(0))
	if err != nil {
		panic(err)
	}
	slog.Info("Deleted",
		"pod", pod.GetName(),
		"namespace", pod.Namespace)
}
