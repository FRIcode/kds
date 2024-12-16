package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MetricDeploymentsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "deployments_total",
		Help: "The total number of deployments",
	}, []string{"name"})

	MetricDeploymentsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "deployments_failed",
		Help: "The total number of failed deployments",
	}, []string{"name"})

	MetricDeploymentsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "deployments_success",
		Help: "The total number of successful deployments",
	}, []string{"name"})

	MetricDeploymentsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "deployments_duration_seconds",
		Help:    "The duration of deployments",
		Buckets: prometheus.DefBuckets,
	}, []string{"name"})

	MetricRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of requests",
	}, []string{"name"})

	MetricRequestsAuthFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_auth_failed",
		Help: "The total number of failed authentication requests",
	}, []string{"name"})

	MetricRequestsAuthSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_auth_success",
		Help: "The total number of successful authentication requests",
	}, []string{"name"})
)
