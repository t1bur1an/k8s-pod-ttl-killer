package k8s

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/t1bur1an/k8s-pod-ttl-killer/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	PodsDeleted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pods_deleted",
		Help: "The total number of deleted pods",
	},
		[]string{"namespace", "annotation"})
	envConfig = config.ReadConfig()
)

func GetKubeConfig() (*rest.Config, error) {
	slog.Info("Trying to read InClusterConfig")
	kubeconfig, err := rest.InClusterConfig()
	if err == nil {
		return kubeconfig, nil
	}
	slog.Warn("InClusterConfig failed", "error", err.Error())

	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		return nil, err // InClusterConfig не работает и KUBECONFIG не задан
	}

	slog.Info("Trying to load kubeconfig from KUBECONFIG env", "path", kubeconfigPath)
	kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	return kubeconfig, nil
}

func CheckClusterPodsPoll(clientset *kubernetes.Clientset) {
	for {
		slog.Info("Check cluster pods: started")
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			slog.Error("Failed to list pods", "error", err)
			continue
		}

		for _, pod := range pods.Items {
			if DeletePodCheck(pod) {
				slog.Info("Pod to delete", "pod", pod.GetName(), "namespace", pod.Namespace)
				podContext := context.Background()
				if err := DeletePod(podContext, clientset, &pod, envConfig.TTLAnnotation); err != nil {
					slog.Error("Failed to delete pod", "pod", pod.GetName(), "namespace", pod.Namespace, "error", err)
				}
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

func DeletePod(ctx context.Context, clientset *kubernetes.Clientset, pod *corev1.Pod, ttlAnnotation string) error {
	err := clientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod %s/%s: %w", pod.Namespace, pod.Name, err)
	}

	slog.Info("Deleted pod",
		"namespace", pod.Namespace,
		"name", pod.Name,
		"reason", "TTL expired",
	)

	PodsDeleted.WithLabelValues(pod.Namespace, ttlAnnotation).Inc()
	return nil
}
