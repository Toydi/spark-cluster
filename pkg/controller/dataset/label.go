package dataset

import (
	"github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
	"github.com/spark-cluster/pkg/controller/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	labelTypeValue = "dataset"
)

func DefaultLabels(ds *v1alpha1.Dataset) labels.Set {
	l := labels.Set{}
	l[internal.LabelNameKey] = ds.Name
	l[internal.LabelTypeKey] = labelTypeValue

	return l
}

func AddUserLabel(ds *v1alpha1.Dataset, user string) {
	if ds.ObjectMeta.Labels == nil {
		ds.ObjectMeta.Labels = make(map[string]string)
	}
	ds.ObjectMeta.Labels[internal.LabelUserKey] = user
}

func SelectorForUser(user string) labels.Selector {
	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			internal.LabelUserKey: user,
		},
	}

	labelSelector, _ := metav1.LabelSelectorAsSelector(selector)

	return labelSelector
}
