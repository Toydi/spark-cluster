package apis

import (
	sparkcluster "github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
	dataset "github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, sparkcluster.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes,dataset.SchemeBuilder.AddToScheme)
}
