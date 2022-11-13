package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	// 测试指标
	ErrPodResource = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "err_pod_resource",
			Help: "Resource of err write pod",
		}, []string{"namespace", "name", "resource_type"})
)

func init() {
	metrics.Registry.MustRegister(ErrPodResource)
}
