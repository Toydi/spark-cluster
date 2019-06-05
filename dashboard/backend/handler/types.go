package handler

import (
	"github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
)

const (
	Namespace = "bdkit"
)

const (
	oauthScopeProfile = "profile"
	oauthScopeEmail   = "email"
)

type SparkClusterList struct {
	SparkClusters []v1alpha1.SparkCluster `json:"sparkclusters"`
}

type User struct {
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	UserID    string `json:"userid"`
	Email     string `json:"email"`
	Signature string `json:"signature"`
}
