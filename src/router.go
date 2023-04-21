package src

import (
	"github.com/gin-gonic/gin"
	"metrics-exporter/src/endpoints"
)

type Router interface {
	Route(engine *gin.Engine)
}

type Loader struct {
}

func (l *Loader) Load() []Router {
	healthcheck := new(endpoints.HealthCheck)
	return []Router{
		healthcheck,
	}
}
