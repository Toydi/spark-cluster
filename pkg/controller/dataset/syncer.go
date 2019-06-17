package dataset

import (
	"reflect"

	"github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
	"github.com/spark-cluster/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewPodSyncer(ds *v1alpha1.Dataset, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	pod := newDataServerPod(ds)
	return syncer.NewObjectSyncer("pod", ds, pod, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1.Pod)
		if !reflect.DeepEqual(out.Spec, pod.Spec) {
			out.Spec = pod.Spec
		}
		return nil
	})
}

func NewServiceSyncer(ds *v1alpha1.Dataset, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := newDataServerService(ds)
	return syncer.NewObjectSyncer("service", ds, svc, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1.Service)
		if !reflect.DeepEqual(out.Spec, svc.Spec) {
			out.Spec = svc.Spec
		}
		return nil
	})
}

func NewPersistentVolumeClaimSyncer(ds *v1alpha1.Dataset, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	pvc := newPersistentVolumeClaim(ds)
	return syncer.NewObjectSyncer("pvc", ds, pvc, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1.PersistentVolumeClaim)
		if !reflect.DeepEqual(out.Spec, pvc.Spec) {
			out.Spec = pvc.Spec
		}
		return nil
	})
}
