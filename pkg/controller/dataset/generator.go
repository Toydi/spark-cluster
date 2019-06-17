package dataset

import (
	"fmt"

	"github.com/spark-cluster/pkg/controller/internal"

	"github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	dataServerImage     = "registry.njuics.cn/qr/quickshare"
	dataServerPort      = 8888
	dataServerPortName  = "data-server"
	dataServerMountPath = "/data"
	dataServerCommand   = "./quickshare"

	defaultStorageSize = "1Gi"

	podNamePrefix = "bdkit-dataset"
	PodMountPath  = "/workspace"
)

func GeneralName(datasetName string) string {
	return fmt.Sprintf("%s-%s", podNamePrefix, datasetName)
}

func newDataServerPod(dataset *v1alpha1.Dataset) *v1.Pod {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GeneralName(dataset.Name),
			Namespace: dataset.Namespace,
			Labels:    DefaultLabels(dataset),
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    dataset.Name,
					Image:   dataServerImage,
					Command: []string{dataServerCommand, dataServerMountPath},
				},
			},
		},
	}

	internal.Volume{
		Name:      GeneralName(dataset.Name),
		MountPath: dataServerMountPath,
	}.AddToPod(pod)

	return pod
}

func newDataServerService(dataset *v1alpha1.Dataset) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GeneralName(dataset.Name),
			Namespace: dataset.Namespace,
			Labels:    DefaultLabels(dataset),
		},
		Spec: v1.ServiceSpec{
			Selector: DefaultLabels(dataset),
			Type:     v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
				{
					Name: dataServerPortName,
					Port: dataServerPort,
				},
			},
		},
	}
}

func newPersistentVolumeClaim(dataset *v1alpha1.Dataset) *v1.PersistentVolumeClaim {
	if len(dataset.Spec.Size) == 0 {
		dataset.Spec.Size = defaultStorageSize
	}
	storageClass := internal.GetStorageClassName()

	return &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GeneralName(dataset.Name),
			Namespace: dataset.Namespace,
			Labels:    DefaultLabels(dataset),
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			StorageClassName: &storageClass,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(dataset.Spec.Size),
				},
			},
		},
	}
}
