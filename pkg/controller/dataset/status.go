package dataset

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
	"github.com/spark-cluster/pkg/controller/internal"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileDataset) updateStatus(dataset *v1alpha1.Dataset) error {
	if dataset.Status.CreateTime == nil {
		now := metav1.Now()
		dataset.Status.CreateTime = &now
	}

	labels := DefaultLabels(dataset)

	pods, err := internal.PodsForLabels(dataset.Namespace, labels, r.client)
	if err != nil {
		return err
	}

	services, err := internal.ServicesForLabels(dataset.Namespace, labels, r.client)
	if err != nil {
		return err
	}

	podStatuses := internal.MappingPodsByPhase(pods)
	if podStatuses[v1.PodRunning] == 1 {
		// All pods are running, set start time
		if dataset.Status.StartTime == nil {
			now := metav1.Now()
			dataset.Status.StartTime = &now
		}

		// if all pods are running, we set the workspace status to running and update endpoint field
		dataset.Status.Phase = v1alpha1.DatasetPhaseRunning
		if err := r.updateEndpoints(dataset, pods, services); err != nil {
			log.Errorf("failed to update workspace endpoint: %v", err)
		}
	} else if podStatuses[v1.PodFailed] > 0 {
		dataset.Status.Phase = v1alpha1.DatasetPhaseFailed
	} else {
		dataset.Status.Phase = v1alpha1.DatasetPhasePending
	}

	return r.syncStatus(dataset)
}

func (r *ReconcileDataset) syncStatus(dataset *v1alpha1.Dataset) error {
	oldDataset := &v1alpha1.Dataset{}
	r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dataset.Name,
		Namespace: dataset.Namespace,
	}, oldDataset)

	if !reflect.DeepEqual(oldDataset.Status, dataset.Status) {
		return r.client.Update(context.TODO(), dataset)
	}

	return nil
}

func (r *ReconcileDataset) updateEndpoints(dataset *v1alpha1.Dataset, pods []v1.Pod, services []v1.Service) error {
	if len(pods) == 0 || len(services) == 0 {
		return fmt.Errorf("Dataset is not running: %v pods, %v services", len(pods), len(services))
	}

	dataset.Status.Endpoint = getEndpoint(pods[0], services[0])
	return nil
}

func getEndpoint(pod v1.Pod, service v1.Service) string {
	if pod.Spec.NodeName == "" || len(service.Spec.Ports) == 0 {
		return ""
	}
	nodeip:="114.212.189."
	s:=strings.Split(pod.Spec.NodeName, "n")
	if len(s)>1{
		nodeip+=s[1]
	}
	return fmt.Sprintf("%s:%d", nodeip, service.Spec.Ports[0].NodePort)
}
