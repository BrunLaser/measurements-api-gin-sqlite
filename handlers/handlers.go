package handlers

import (
	"Go-Check24/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) MeasurementPost(c *gin.Context) {
	//check for json header (not need as shouldbindjson should reject)
	/*if c.ContentType() != "application/json" {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Kein JSON Header"})
	    return
	}*/
	newPoint := &database.Measurement{}
	//json to struct
	if err := c.ShouldBindJSON(newPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Data"})
		return
	}
	//struct to database
	if err := h.db.InsertPoint(newPoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database INSERT error"})
		return
	}
	location := fmt.Sprintf("/measurements/%d", newPoint.ID) //The ID is set to the actual ID in the database
	c.Header("Location", location)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Point created",
		"location": location,
		"data":     *newPoint,
	})
}

func (h *Handler) MeasurementGetAll(c *gin.Context) {
	//w.Header().Set("Content-Type", "application/json") //gin does the header when I do json stuff
	measurements, err := h.db.GetAllPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	c.JSON(http.StatusOK, measurements) //write json data back
}
