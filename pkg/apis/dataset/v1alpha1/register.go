// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the dataset v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=dataset.njuics.cn
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

var (
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	AddToScheme = SchemeBuilder.AddToScheme
)

const (
	// GroupName is the group name use in this package.
	GroupName = "spark.k8s.io"
	// Kind is the kind name.
	Kind = "Dataset"
	// GroupVersion is the version.
	GroupVersion = "v1alpha1"
	// Plural is the Plural for Workspace.
	Plural = "datasets"
	// Singular is the singular for Workspace.
	Singular = "dataset"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

	// SchemeGroupVersionKind is the GroupVersionKind of the resource.
	SchemeGroupVersionKind = SchemeGroupVersion.WithKind(Kind)
)

// Resource takes an unqualified resource and returns a Group-qualified GroupResource.
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Dataset{},
		&DatasetList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
