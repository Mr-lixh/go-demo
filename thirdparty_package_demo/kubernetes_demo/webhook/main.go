package main

import (
	flag "github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	goflag "flag"
	"fmt"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	// TODO: 其它 flag 数据库信息
	version              = flag.BoolP("version", "v", false, "Version.")
	metricsAddr          = flag.String("metrics-addr", ":8080", "The address the metrics endpoint binds to.")
	enableLeaderElection = flag.Bool("enable-leader-election", false, "Enable leader election for controller manager.")
	certDir              = flag.String("tls-cert-dir", "", "The directory thar contains the server key and certificate."+
		"The server key and certificate must be named tls.key and tls.crt, respectively.")
	port = flag.Int("port", 9443, "Port of webhook service.")
	dsn  = flag.String("dsn", "", "MySQL dsn, format as \"username:password@tcp(host:port)/database\".")
)

func main() {
	// Setup logger and parse flag.
	var loggerOpts = &zap.Options{}
	var goFlagSet goflag.FlagSet
	loggerOpts.BindFlags(&goFlagSet)
	flag.CommandLine.AddGoFlagSet(&goFlagSet)
	flag.Parse()

	if *version {
		fmt.Println("version: 1.0.0")
		os.Exit(0)
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(loggerOpts)))

	entryLog := ctrl.Log.WithName("entrypoint")

	// Setup a Manager.
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		MetricsBindAddress:      *metricsAddr,
		CertDir:                 *certDir,
		Port:                    *port,
		LeaderElection:          *enableLeaderElection,
		LeaderElectionNamespace: "default",
		LeaderElectionID:        "hybrid-cloud-webhook",
	})
	if err != nil {
		entryLog.Error(err, "Unable to create manager")
		os.Exit(1)
	}

	// Setup webhooks.
	entryLog.Info("Setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("Registering webhooks to the webhook server")
	hookServer.Register("/v1/hybrid-cloud-pod", &webhook.Admission{Handler: &HybridCloud{Client: mgr.GetClient()}})

	entryLog.Info("Starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "Unable to run manager")
		os.Exit(1)
	}
}
