package src

import (
	"metrics-exporter/src/endpoints"
	"net/http"
)

type Router interface {
	Route(engine *http.ServeMux)
}

type Loader struct {
}

func (l *Loader) Load() []Router {
	healthcheck := new(endpoints.HealthCheck)
	metrics := &endpoints.Metrics{}
	return []Router{
		healthcheck,
		metrics,
	}
}
