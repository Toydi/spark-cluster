package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"io/ioutil"

	"github.com/gorilla/mux"
	datasetv1alpha1 "github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
	"github.com/spark-cluster/pkg/controller/dataset"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListDatasets(w http.ResponseWriter, r *http.Request) {
	dsl := &datasetv1alpha1.DatasetList{}
	// opts := &client.ListOptions{}
	// // opts.SetLabelSelector(fmt.Sprintf("type=%s", "spark-cluster"))
	// // opts.InNamespace(Namespace)
	// opts.LabelSelector = sparkcluster.SelectorForUser(user)
	// err := handler.client.List(context.TODO(), opts, sc)


	err := handler.client.List(context.TODO(), &client.ListOptions{}, dsl)

	if err != nil {
		log.Warningf("failed to list datasets: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(DatasetList{Datasets: dsl.Items}, w, http.StatusOK)
	}
}

func (handler *APIHandler) CreateDataset(w http.ResponseWriter, r *http.Request) {
	ds := new(datasetv1alpha1.Dataset)
	user := r.Header.Get("User")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &ds); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}

	if len(ds.Namespace) == 0 {
		ds.Namespace = ResourceNamespace
	}
	dataset.AddUserLabel(ds, user)

	err = handler.client.Create(context.TODO(), ds)
	if err != nil {
		log.Warningf("Failed to create dataset %v: %v", ds.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(ds, w, http.StatusCreated)
	}
}

func (handler *APIHandler) DeleteDataset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["dataset"]

	ds := new(datasetv1alpha1.Dataset)
	ds.Name = name
	ds.Namespace = ResourceNamespace

	err := handler.client.Delete(context.TODO(), ds)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("", w, http.StatusOK)
	}
}
