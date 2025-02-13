package router

import (
	"Go-Check24/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handler) {
	r.GET("/measurements", h.MeasurementGetAll)
	r.POST("/measurements", h.MeasurementPost)
	r.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hallo!") })
}
