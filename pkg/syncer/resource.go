package syncer

import (
	"fmt"

	"reflect"

	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewDeploySyncer(name string, owner, obj runtime.Object, c client.Client, scheme *runtime.Scheme) Interface {
	template := obj.(*appsv1.Deployment)
	metaobj := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      template.ObjectMeta.Name,
			Namespace: template.ObjectMeta.Namespace,
		},
	}
	return NewObjectSyncer(name, owner, metaobj, c, scheme, func(existing runtime.Object) error {
		out := existing.(*appsv1.Deployment)

		out.Spec.Template.ObjectMeta = template.Spec.Template.ObjectMeta
		selector := metav1.SetAsLabelSelector(template.Spec.Template.ObjectMeta.Labels)
		if !reflect.DeepEqual(selector, out.Spec.Selector) {
			if out.ObjectMeta.CreationTimestamp.IsZero() {
				out.Spec.Selector = selector
			} else {
				return fmt.Errorf("deployment selector is immutable")
			}
		}

		err := mergo.Merge(&out.Spec.Template.Spec, template.Spec.Template.Spec, mergo.WithTransformers(PodSpecTransformer))
		if err != nil {
			return err
		}
		return nil
	})
}

func NewServiceSyncer(name string, owner, obj runtime.Object, c client.Client, scheme *runtime.Scheme) Interface {
	template := obj.(*v1.Service)
	metaobj := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      template.ObjectMeta.Name,
			Namespace: template.ObjectMeta.Namespace,
		},
	}
	return NewObjectSyncer(name, owner, metaobj, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1.Service)

		out.Labels = template.Labels

		if len(out.Spec.Ports) != 1 {
			out.Spec.Ports = make([]v1.ServicePort, 1)
		}
		out.Spec.Ports[0].Name = template.Spec.Ports[0].Name
		out.Spec.Ports[0].Port = template.Spec.Ports[0].Port
		out.Spec.Ports[0].TargetPort = template.Spec.Ports[0].TargetPort
		out.Spec.Type = template.Spec.Type
		out.Spec.Selector = template.Spec.Selector

		return nil
	})
}
