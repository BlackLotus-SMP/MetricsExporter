package src

import (
	"github.com/gin-gonic/gin"
	"metrics-exporter/src/endpoints"
	"metrics-exporter/src/minecraft"
)

type Router interface {
	Route(engine *gin.Engine)
}

type Loader struct {
	Metrics minecraft.MCMetrics
}

func (l *Loader) Load() []Router {
	healthcheck := new(endpoints.HealthCheck)
	metrics := &endpoints.Metrics{Metrics: l.Metrics}
	return []Router{
		healthcheck,
		metrics,
	}
}
