package endpoints

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Metrics struct {
	registry *prometheus.Registry
}

func NewMetricsEndpoint(promReg *prometheus.Registry) *Metrics {
	m := new(Metrics)
	m.registry = promReg
	return m
}

func (m Metrics) Route(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))
}
