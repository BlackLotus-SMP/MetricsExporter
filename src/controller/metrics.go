package controller

import (
	"github.com/gin-gonic/gin"
	"metrics-exporter/src/minecraft"
	"net/http"
)

type Controller struct {
	Metrics minecraft.MCMetrics
}

func (c Controller) Prom(ctx *gin.Context) {
	m := c.Metrics.GetMetrics()
	if m == nil {
		ctx.JSON(http.StatusNotFound, "")
		return
	}
	ctx.JSON(http.StatusOK, m)
}
