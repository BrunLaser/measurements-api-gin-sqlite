package router

import (
	"Go-Check24/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handler) {
	r.GET("/measurements", h.MeasurementGetAll)
	r.GET("/measurements/:id", h.MeasurementGetById)
	r.POST("/measurements", h.MeasurementPost)
	r.DELETE("/measurements/:id", h.MeasurementDelete)
	r.PUT("/measurements/:id", h.MeasurementUpdate)

	//r.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hallo!") })
}
