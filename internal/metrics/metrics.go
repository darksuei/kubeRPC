package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RegisteredServices = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "kuberpc_registered_services",
		Help: "Current number of registered services",
	})

	MethodRegistrations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "kuberpc_method_registrations_total",
		Help: "Total number of method registrations",
	})

	MethodLookups = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kuberpc_method_lookups_total",
		Help: "Number of method resolution lookups by service and method",
	}, []string{"service", "method"})

	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kuberpc_http_requests_total",
		Help: "Total HTTP requests handled by kubeRPC core",
	}, []string{"method", "path", "status"})
)
