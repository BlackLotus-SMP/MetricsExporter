package src

import (
	"github.com/prometheus/client_golang/prometheus"
	"metrics-exporter/src/endpoints"
	"net/http"
)

type Router interface {
	Route(engine *http.ServeMux)
}

type Loader struct {
}

func (l *Loader) Load(promReg *prometheus.Registry) []Router {
	healthcheck := new(endpoints.HealthCheck)
	metrics := endpoints.NewMetricsEndpoint(promReg)
	return []Router{
		healthcheck,
		metrics,
	}
}
