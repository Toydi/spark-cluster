package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	sparkclusterv1alpha1 "github.com/spark-cluster/pkg/apis/spark-cluster/v1alpha1"
	"github.com/gorilla/mux"
	"github.com/spark-cluster/pkg/controller/sparkcluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListSparkCluster(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("User")
	shareflag := "true"
	sc := &sparkclusterv1alpha1.SparkClusterList{}
	sc_share := &sparkclusterv1alpha1.SparkClusterList{}
	//TODO() add label selector
	opts := &client.ListOptions{}
	opts_share := &client.ListOptions{}
	// opts.SetLabelSelector(fmt.Sprintf("type=%s", "spark-cluster"))
	// opts.InNamespace(Namespace)
	opts.LabelSelector = sparkcluster.SelectorForUser(user)
	opts_share.LabelSelector = sparkcluster.SelectorForShare(shareflag)
	//shared cluster
	err := handler.client.List(context.TODO(), opts, sc)
	if err != nil {
		log.Warningf("failed to list spark cluster: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	err = handler.client.List(context.TODO(),opts_share,sc_share)
	if err != nil {
		log.Warningf("failed to list spark cluster: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	sc_len:=len(sc_share.Items)
	for i := 0; i < sc_len; i++ {
		exists := false
		for j := 0; j < len(sc.Items); j++ {
			if sc_share.Items[i].Name==sc.Items[j].Name{
				exists = true
				break
			}
		}
		if !exists{
			sc.Items=append(sc.Items,sc_share.Items[i])
		}
	}
	gitRepos:=make([]string,len(sc.Items))
	for i := 0; i < len(sc.Items); i++ {
		gitRepos[i]=sc.Items[i].Spec.GitRepo
	}

	responseJSON(SparkClusterList{SparkClusters: sc.Items,GitRepos:gitRepos}, w, http.StatusOK)
}

func (handler *APIHandler) CreateSparkCluster(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("User")
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

	shared:=sc.Spec.Shared
	if len(sc.Name)==0{
		// sc.Name=sc.Spec.ClusterPrefix+"-cluster"
		sc.Name=fmt.Sprintf("%s-%s", sc.Spec.ClusterPrefix,"cluster")
	}

	if len(sc.Namespace) == 0 {
		sc.Namespace = ResourceNamespace
	}

	if len(sc.Spec.Shared)!=0{
		sparkcluster.AddUserLabel(sc, shared)
	}
	// workspace.AddUserLabel(sc, user)
	sparkcluster.AddUserLabel(sc, user)
	err = handler.client.Create(context.TODO(), sc)
	if err != nil {
		log.Warningf("Failed to create spark cluster %v: %v clusterprefix:%v", sc.Name, err,sc.Spec.ClusterPrefix)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(sc, w, http.StatusCreated)
	}
}

func (handler *APIHandler) DeleteSparkCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["sparkcluster"]

	sc := new(sparkclusterv1alpha1.SparkCluster)
	sc.Name = name
	sc.Namespace = ResourceNamespace

	err := handler.client.Delete(context.TODO(), sc)
	if err != nil {
		log.Warningf("failed to delete sparkcluster %v under namespace %v: %v", name, ResourceNamespace, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("", w, http.StatusOK)
	}
}

func (handler *APIHandler) UpdateSparkCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["sparkcluster"]

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

	if sc.ObjectMeta.Name != name {
		err := fmt.Errorf("sparkcluster name in path is not the same as that in json.")
		responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		return
	}

	err = handler.client.Update(context.TODO(), sc)
	if err != nil {
		log.Warningf("Failed to create sparkcluster %v: %v", sc.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)

	} else {
		responseJSON(sc, w, http.StatusOK)
	}
}