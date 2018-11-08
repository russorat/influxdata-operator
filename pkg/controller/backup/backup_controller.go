package backup

import (
	"context"
	"fmt"
	"log"

	influxdatav1alpha1 "bitbucket.org/AhmedDev9/influxdataoperator/pkg/apis/influxdata/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
)

const (
	backupDir  = "/var/lib/influxdb/backup"
)

func Add(mgr manager.Manager) error {
	log.Printf("Adding to manager...")
	return add(mgr, newReconciler(mgr))
}

type ReconcileInfluxdbBackup struct {
	// TODO: Clarify the split client
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileInfluxdbBackup{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	log.Println("in add func")
	// Create a new controller
	c, err := controller.New("backup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}


	log.Println("about to watch for backup primary")
	// Watch for changes to primary resource Backup
	err = c.Watch(&source.Kind{Type: &influxdatav1alpha1.Backup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}


	log.Println("about to watch for secondary")
	// Watch for changes to secondary resource Pods and requeue the owner Backup
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &influxdatav1alpha1.Backup{},
	})
	if err != nil {
		return err
	}

	log.Println("Watches done")
	return nil
}

var _ reconcile.Reconciler = &ReconcileInfluxdbBackup{}

// This gets called when a Backup resource is created... I think.
func (r *ReconcileInfluxdbBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Println("Starting Influxdb Backup")

	// Fetch the Influxdb Backup instance
	backup := &influxdatav1alpha1.Backup{}
	err := r.client.Get(context.TODO(), request.NamespacedName, backup)
	log.Println(err)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Printf("Influxdb Backup %s/%s not found. Ignoring since object must be deleted\n", request.Namespace, request.Name)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Printf("Failed to get Influxdb Backup: %v", err)
		return reconcile.Result{}, err
	}

	cmdOpts := []string{
		"influxd",
		"backup",
		"-portable",
		"-database " + backup.Spec.Database,
		backupDir,
	}

	cmd := strings.Join(cmdOpts, " ")
	log.Printf("Command: %q", cmd)

	podName := "influxdb-0"
	containerName := "influxdb"

	output, stderr, err := ExecToPodThroughAPI(cmd, containerName, podName,	request.Namespace, nil)

	if len(stderr) != 0 {
		log.Println("STDERR:", stderr)
	}

	if err != nil {
		log.Printf("Error occured while `exec`ing to the Pod %q, namespace %q, command %q. Error: %+v\n", podName, request.Namespace, cmd, err)
		return reconcile.Result{}, err
	} else {
		log.Println("Backup Output:")
		fmt.Println(output)
		return reconcile.Result{}, nil
	}
}
