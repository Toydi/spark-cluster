package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DatasetSpec defines the desired state of Dataset
type DatasetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

	// Size is the size of the dataset.
	Size string `json:"size"`

	// Description is the description to this dataset
	Description string `json:"description"`

	// Shared or not 
	Shared string
	
	// ReadOnly shows if this data set is editable
	Readonly string `json:"readonly"`

	// ThumbnailURL is the link of dataset thumbnail
	ThumbnailURL string `json:"thumbnailURL"`
}

// DatasetPhase defines all phase of dataset lifecycle.
type DatasetPhase string

const (
	// DatasetPhasePending means one or some of the containers, storages,
	// or services are creating.
	DatasetPhasePending = "Pending"

	// DatasetPhaseRunning means dataset have been successfully scheduled and launched
	// and it is running without error.
	DatasetPhaseRunning = "Running"

	// DatasetPhaseFailed means some pods of dataset have failed.
	DatasetPhaseFailed = "Failed"
)

// DatasetStatus defines the observed state of Dataset
type DatasetStatus struct {
	// Pod status is the status of dataset pod
	PodStatus v1.PodStatus `json:"podStatus"`

	// Phase show the running phase of dataset.
	Phase DatasetPhase `json:"phase"`

	// CreateTime represents time when the dataset was created.
	CreateTime *metav1.Time `json:"createTime,omitempty"`

	// StartTime represents time when the dataset was started.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Endpoint is the endpoint to access dataset pod.
	Endpoint string `json:"endpoint"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Dataset is the Schema for the datasets API
// +k8s:openapi-gen=true
type Dataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasetSpec   `json:"spec,omitempty"`
	Status DatasetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DatasetList contains a list of Dataset
type DatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dataset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dataset{}, &DatasetList{})
}
