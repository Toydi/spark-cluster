package controller

import (
	"github.com/spark-cluster/pkg/controller/dataset"
	"github.com/spark-cluster/pkg/controller/sparkcluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sparkcluster.Add)
	AddToManagerFuncs = append(AddToManagerFuncs, dataset.Add)
}
