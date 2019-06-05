package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"io/ioutil"

	log "github.com/sirupsen/logrus"
	sparkclusterv1alpha1 "github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
	"github.com/spark-cluster/pkg/controller/sparkcluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListSparkCluster(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("User")

	sc := &sparkclusterv1alpha1.SparkClusterList{}
	//TODO() add label selector
	opts := &client.ListOptions{}
	// opts.SetLabelSelector(fmt.Sprintf("type=%s", "spark-cluster"))
	// opts.InNamespace(Namespace)
	opts.LabelSelector = sparkcluster.SelectorForUser(user)
	err := handler.client.List(context.TODO(), opts, sc)

	if err != nil {
		log.Warningf("failed to list spark cluster: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(SparkClusterList{SparkClusters: sc.Items}, w, http.StatusOK)
	}
}

func (handler *APIHandler) CreateSparkCluster(w http.ResponseWriter, r *http.Request) {
	user := w.Header().Get("User")

	sc := new(sparkclusterv1alpha1.SparkCluster)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &sc); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}

	if len(sc.Namespace) == 0 {
		sc.Namespace = Namespace
	}
	// workspace.AddUserLabel(ws, user)
	sparkcluster.AddUserLabel(sc, user)
	err = handler.client.Create(context.TODO(), sc)
	if err != nil {
		log.Warningf("Failed to create spark cluster %v: %v", sc.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(sc, w, http.StatusCreated)
	}
}
