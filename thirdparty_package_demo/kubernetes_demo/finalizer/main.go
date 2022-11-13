package main

import (
	flag "github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	goflag "flag"
	"github.com/Mr-lixh/go-demo/thirdparty_package_demo/kubernetes_demo/finalizer/controllers"
	"os"
)

var (
	// TODO: 其它 flag 数据库信息、排除或者关注的 namespace 列表、是否每天统计一次数据
	metricsAddr          = flag.String("metrics-addr", ":8080", "The address the metrics endpoint binds to.")
	enableLeaderElection = flag.Bool("enable-leader-election", false, "Enable leader election for controller manager.")
	namespaces           = flag.String("namespaces", "", "Namespace list.")
)

func main() {
	// Setup logger and arse flag.
	var loggerOpts = &zap.Options{}
	var goFlagSet goflag.FlagSet
	loggerOpts.BindFlags(&goFlagSet)
	flag.CommandLine.AddGoFlagSet(&goFlagSet)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(loggerOpts)))

	entryLog := ctrl.Log.WithName("entrypoint")

	// Setup a Manager.
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		MetricsBindAddress:      *metricsAddr,
		LeaderElection:          *enableLeaderElection,
		LeaderElectionNamespace: "default",
		LeaderElectionID:        "pod-cost-controller",
	})
	if err != nil {
		entryLog.Error(err, "Unable to create manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile Pods.
	if err := (&controllers.PodCostReconciler{
		Client:     mgr.GetClient(),
		Namespaces: *namespaces,
	}).SetupWithManager(mgr); err != nil {
		entryLog.Error(err, "Unable to create controller", "controller", "pod-cost")
		os.Exit(1)
	}

	entryLog.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "Problem running manager")
		os.Exit(1)
	}
}
