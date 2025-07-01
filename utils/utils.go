package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FilterAnnotations(filterAnnotations map[string]string, targetAnnotation string) (int64, error) {
	outputDuration := int64(0)
	for annotation, value := range filterAnnotations {
		if annotation == targetAnnotation {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return 0, err
			}
			outputDuration = int64(intValue)
			break
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
	slog.Info(
		"Pod time info",
		"pod", pod.GetName(),
		"podReadyTimestamp",
		podReadyTimestamp, "ttl", ttl)

	utcTimeNow := time.Now().UTC().Unix()
	return (podReadyTimestamp + ttl) >= utcTimeNow
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
