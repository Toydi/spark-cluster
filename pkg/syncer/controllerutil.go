package syncer

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MutateFn func(existing runtime.Object) error

// OperationResult is the action result of a CreateOrUpdate call
type OperationResult string

const ( // They should complete the sentence "Deployment default/foo has been ..."
	// OperationResultNone means that the resource has not been changed
	OperationResultNone OperationResult = "unchanged"
	// OperationResultCreated means that a new resource is created
	OperationResultCreated OperationResult = "created"
	// OperationResultUpdated means that an existing resource is updated
	OperationResultUpdated OperationResult = "updated"
)

// CreateOrUpdate creates or updates the given object obj in the Kubernetes
// cluster. The object's desired state should be reconciled with the existing
// state using the passed in ReconcileFn. obj must be a struct pointer so that
// obj can be updated with the content returned by the Server.
//
// It returns the executed operation and an error.
func CreateOrUpdate(ctx context.Context, c client.Client, obj runtime.Object, f MutateFn) (OperationResult, error) {
	// op is the operation we are going to attempt
	op := OperationResultNone

	// get the existing object meta
	metaObj, ok := obj.(v1.Object)
	if !ok {
		return OperationResultNone, fmt.Errorf("%T does not implement metav1.Object interface", obj)
	}

	// retrieve the existing object
	key := client.ObjectKey{
		Name:      metaObj.GetName(),
		Namespace: metaObj.GetNamespace(),
	}
	err := c.Get(ctx, key, obj)

	// reconcile the existing object
	existing := obj.DeepCopyObject()
	existingObjMeta := existing.(v1.Object)
	existingObjMeta.SetName(metaObj.GetName())
	existingObjMeta.SetNamespace(metaObj.GetNamespace())

	if e := f(obj); e != nil {
		return OperationResultNone, e
	}

	if metaObj.GetName() != existingObjMeta.GetName() {
		return OperationResultNone, fmt.Errorf("ReconcileFn cannot mutate objects name")
	}

	if metaObj.GetNamespace() != existingObjMeta.GetNamespace() {
		return OperationResultNone, fmt.Errorf("ReconcileFn cannot mutate objects namespace")
	}

	if errors.IsNotFound(err) {
		err = c.Create(ctx, obj)
		op = OperationResultCreated
	} else if err == nil {
		if reflect.DeepEqual(existing, obj) {
			return OperationResultNone, nil
		}
		err = c.Update(ctx, obj)
		op = OperationResultUpdated
	} else {
		return OperationResultNone, err
	}

	if err != nil {
		op = OperationResultNone
	}
	return op, err
}
