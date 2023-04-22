package endpoints

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Metrics struct {
}

func (m Metrics) Route(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
