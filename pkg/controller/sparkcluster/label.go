package sparkcluster

import (
	"context"
	"fmt"

	sparkv1alpha1 "github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
	"github.com/spark-cluster/pkg/controller/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func masterName(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + Master
}

func MasterName(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + Master
}

func uiServiceName(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + "ui-service"
}

func masterLabel(instance *sparkv1alpha1.SparkCluster) map[string]string {
	return map[string]string{"app": masterName(instance)}
}

func masterPvc(instance *sparkv1alpha1.SparkCluster) string {
	return instance.Spec.ClusterPrefix + "-" + MasterPvc
}

func slaveName(instance *sparkv1alpha1.SparkCluster, index int) string {
	return instance.Spec.ClusterPrefix + "-" + Slave + "-" + fmt.Sprintf("%d", index)
}

func slaveLabel(instance *sparkv1alpha1.SparkCluster, index int) map[string]string {
	return map[string]string{"app": slaveName(instance, index), "type": instance.Spec.ClusterPrefix + "-slave"}
}

func slavePvc(instance *sparkv1alpha1.SparkCluster, index int) string {
	return instance.Spec.ClusterPrefix + "-" + SlavePvc + "-" + fmt.Sprintf("%d", index)
}
func (r *ReconcileSparkCluster) updateLabels(instance *sparkv1alpha1.SparkCluster) error {
	flag := false
	log.Info("enter the updatelabels function")
	if instance.Labels == nil {
		flag = true
		instance.Labels = map[string]string{
			"app":  instance.Spec.ClusterPrefix + "-" + "cluster",
			"type": "spark-cluster",
		}
	} else if _, ok := instance.Labels["type"]; !ok {
		flag = true
		instance.Labels["type"] = "spark-cluster"
		log.Info("Instance labels  updating", " label: ", instance.Labels)
	}
	if flag {
		if err := r.client.Update(context.TODO(), instance); err != nil {
			log.Error(err, "update labels error", "instance name:", instance.Name)
			return err
		}
		log.Info("Instance labels  updated", " label: ", instance.Labels)
	}
	return nil
}
func AddUserLabel(sc *sparkv1alpha1.SparkCluster, user string) {
	if sc.ObjectMeta.Labels == nil {
		sc.ObjectMeta.Labels = make(map[string]string)
	}
	sc.ObjectMeta.Labels[internal.LabelUserKey] = user
}

func AddSharedLabel(sc *sparkv1alpha1.SparkCluster, shared string) {
	if sc.ObjectMeta.Labels == nil {
		sc.ObjectMeta.Labels = make(map[string]string)
	}
	sc.ObjectMeta.Labels[internal.LabelShareKey] = shared
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

func SelectorForShare(flag string) labels.Selector{
	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			internal.LabelShareKey: flag,
		},
	}
	labelSelector, _ := metav1.LabelSelectorAsSelector(selector)

	return labelSelector	
}
