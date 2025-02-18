package router

import (
	"Go-Check24/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handler) {
	r.GET("/measurements", h.HandleMeasurementGetAll)
	r.POST("/measurements", h.HandleMeasurementPost)
	r.GET("/measurements/:id", h.HandleMeasurementGetById)
	r.DELETE("/measurements/:id", h.HandleMeasurementDelete)
	r.PUT("/measurements/:id", h.HandleMeasurementUpdate)

	r.GET("experiments/:exp/measurements", h.HandleGetMeasurementsByExperiment)
	//r.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hallo!") })
}
