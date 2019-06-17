package handler

import (
	sparkv1alha1"github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
	datasetv1alha1"github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
)

const (
	Namespace = "bdkit"
	ResourceNamespace = "bdkit-resources"
)

const (
	oauthScopeProfile = "profile"
	oauthScopeEmail   = "email"
)

type SparkClusterList struct {
	SparkClusters []sparkv1alha1.SparkCluster `json:"sparkclusters"`
	GitRepos []string `json:"gitrepos"`
} 

type DatasetList struct{
	Datasets []datasetv1alha1.Dataset `json:"datasets"`
}

type User struct {
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	UserID    string `json:"userid"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}
