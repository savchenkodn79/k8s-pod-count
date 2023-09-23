package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ReplicaSetReconciler struct {
	client.Client
}

func (a *ReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rs := &appsv1.ReplicaSet{}
	err := a.Get(ctx, req.NamespacedName, rs)
	if err != nil {
		return ctrl.Result{}, err
	}

	pods := &corev1.PodList{}
	err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
	if err != nil {
		return ctrl.Result{}, err
	}

	rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
	err = a.Update(context.TODO(), rs)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func main() {
	var log = ctrl.Log.WithName("builder-examples")

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	err = ctrl.NewControllerManagedBy(manager).For(&appsv1.ReplicaSet{}).Owns(&corev1.Pod{}).Complete(&ReplicaSetReconciler{Client: manager.GetClient()})
	if err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}
