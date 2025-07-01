package utils

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FilterAnnotations(filterAnnotations map[string]string, targetAnnotation string) (time.Duration, error) {
	outputDuration := time.Duration(0)
	for annotation, value := range filterAnnotations {
		if annotation == targetAnnotation {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return outputDuration, err
			}
			outputDuration = time.Duration(intValue)
			break
		}
	}
	errMsg := fmt.Sprintf("No %s annotation was found", targetAnnotation)
	return outputDuration, errors.New(errMsg)
}

func GetPodReadyTime(pod corev1.Pod) (int64, error) {
	for _, podCondition := range pod.Status.Conditions {
		if podCondition.Type == corev1.PodReady {
			return podCondition.LastTransitionTime.Unix(), nil
		}
	}
	return 0, errors.New("Pod is not ready")
}

func DeletePod(clientset kubernetes.Clientset, namespaceName string, podContext context.Context, containerName string) {
	err := clientset.CoreV1().Pods(namespaceName).Delete(podContext, containerName, *metav1.NewDeleteOptions(0))
					if err != nil {
									panic(err)
					}
}
