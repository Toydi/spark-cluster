package dataset

import (
	"context"

	log "github.com/sirupsen/logrus"
	datasetv1alpha1 "github.com/spark-cluster/pkg/apis/dataset/v1alpha1"
	"github.com/spark-cluster/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName = "dataset-controller"

// Add creates a new Dataset Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDataset{client: mgr.GetClient(), scheme: mgr.GetScheme(), recorder: mgr.GetRecorder(controllerName)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dataset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Dataset
	err = c.Watch(&source.Kind{Type: &datasetv1alpha1.Dataset{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	subresources := []runtime.Object{
		&v1.Pod{},
		&v1.PersistentVolumeClaim{},
		&v1.Service{},
	}

	for _, subresource := range subresources {
		err = c.Watch(&source.Kind{Type: subresource}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &datasetv1alpha1.Dataset{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDataset{}

// ReconcileDataset reconciles a Dataset object
type ReconcileDataset struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a Dataset object and makes changes based on the state read
// and what is in the Dataset.Spec
func (r *ReconcileDataset) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconciling Dataset %s/%s\n", request.Namespace, request.Name)

	// Fetch the Dataset instance
	dataset := &datasetv1alpha1.Dataset{}
	err := r.client.Get(context.TODO(), request.NamespacedName, dataset)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// if the dataset is terminating, stop reconcile
	if dataset.ObjectMeta.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	syncers := []syncer.Interface{
		NewPersistentVolumeClaimSyncer(dataset, r.client, r.scheme),
		NewPodSyncer(dataset, r.client, r.scheme),
		NewServiceSyncer(dataset, r.client, r.scheme),
	}
	if err := r.sync(syncers); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(dataset)
}

func (r *ReconcileDataset) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			return err
		}
	}
	return nil
}
