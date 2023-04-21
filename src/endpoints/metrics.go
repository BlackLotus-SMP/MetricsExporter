package endpoints

import (
	"github.com/gin-gonic/gin"
	"metrics-exporter/src/controller"
	"metrics-exporter/src/minecraft"
)

type Metrics struct {
	Metrics minecraft.MCMetrics
}

func (m Metrics) Route(engine *gin.Engine) {
	c := controller.Controller{Metrics: m.Metrics}
	engine.GET("/metrics", c.Prom)
}
