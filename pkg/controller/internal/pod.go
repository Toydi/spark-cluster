package internal

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodsForLabels(namespace string, labels labels.Set, c client.Client) ([]v1.Pod, error) {
	pods := &v1.PodList{}
	err := c.List(context.TODO(),
		client.InNamespace(namespace).
			MatchingLabels(labels), pods)

	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func PodsForMember(pods []v1.Pod, mtype string) ([]v1.Pod, error) {
	var result []v1.Pod

	memberSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			LabelMemberKey: strings.ToLower(mtype),
		},
	}

	for _, pod := range pods {
		selector, err := metav1.LabelSelectorAsSelector(memberSelector)
		if err != nil {
			return nil, err
		}
		if !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}
		result = append(result, pod)
	}

	return result, nil
}

func MappingPodsByPhase(pods []v1.Pod) map[v1.PodPhase]int {
	result := make(map[v1.PodPhase]int)
	for _, pod := range pods {
		if len(pod.Status.Phase) == 0 {
			continue
		}
		if _, ok := result[pod.Status.Phase]; !ok {
			result[pod.Status.Phase] = 1
		} else {
			result[pod.Status.Phase]++
		}
	}
	return result
}

func ServicesForLabels(namespace string, labels labels.Set, c client.Client) ([]v1.Service, error) {
	services := &v1.ServiceList{}
	err := c.List(context.TODO(),
		client.InNamespace(namespace).
			MatchingLabels(labels), services)
	if err != nil {
		return nil, err
	}
	return services.Items, nil
}
