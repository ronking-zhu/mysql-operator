/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	dbv1 "my.domain/mysql-operator/api/v1"
)

// MySQLClusterReconciler reconciles a MySQLCluster object
type MySQLClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=db.my.domain,resources=mysqlclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.my.domain,resources=mysqlclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.my.domain,resources=mysqlclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MySQLCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
//
//	func (r *MySQLClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
//		_ = log.FromContext(ctx)
//
//		// TODO(user): your logic here
//
//		return ctrl.Result{}, nil
//	}
//
// +kubebuilder:rbac:groups=db.my.domain,resources=mysqlclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=db.my.domain,resources=mysqlclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=db.my.domain,resources=mysqlclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
func (r *MySQLClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1. Fetch the MySQLCluster instance
	mysql := &dbv1.MySQLCluster{}
	if err := r.Get(ctx, req.NamespacedName, mysql); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Define the desired StatefulSet object
	found := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, found)

	// 3. If StatefulSet doesn't exist, create it
	if err != nil && errors.IsNotFound(err) {
		dep := r.statefulSetForMySQL(mysql)
		log.Info("Creating a new StatefulSet", "Namespace", dep.Namespace, "Name", dep.Name)
		if err := r.Create(ctx, dep); err != nil {
			return ctrl.Result{}, err
		}
		// Requeue to check status again
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// 4. Update Replicas if changed (Self-Healing / Day 2 Ops)
	if *found.Spec.Replicas != mysql.Spec.Replicas {
		found.Spec.Replicas = &mysql.Spec.Replicas
		log.Info("Updating Replicas", "From", *found.Spec.Replicas, "To", mysql.Spec.Replicas)
		if err := r.Update(ctx, found); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MySQLClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.MySQLCluster{}).
		Complete(r)
}
func (r *MySQLClusterReconciler) statefulSetForMySQL(m *dbv1.MySQLCluster) *appsv1.StatefulSet {
	ls := map[string]string{"app": "mysql-operator", "mysql_cluster": m.Name}
	replicas := m.Spec.Replicas

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: ls},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: ls},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "mysql:5.7",
						Name:  "mysql",
						Env: []corev1.EnvVar{{
							Name:  "MYSQL_ROOT_PASSWORD",
							Value: "password", // In production, use Secrets!
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
					}},
				},
			},
		},
	}
	// Set controller reference so deleting the CR deletes the Pods
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
