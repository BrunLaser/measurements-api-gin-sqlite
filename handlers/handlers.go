package handlers

import (
	"Go-Check24/database"
	"Go-Check24/util"
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

/*
curl -X POST http://localhost:8080/measurements
	 -H "Content-Type: application/json"
	 -d '{"sensor_id": 1, "value": 23.45, "unit": "Temp"}'
*/

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
	if err := h.db.InsertMeasurement(newPoint); err != nil {
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
	measurements, err := h.db.GetAllMeasurements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, measurements) //write json data back
}

func (h *Handler) MeasurementGetById(c *gin.Context) {
	id, err := util.GetParamInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	point, err := h.db.GetMeasurementById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, point)
}

func (h *Handler) MeasurementDelete(c *gin.Context) {
	id, err := util.GetParamInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.db.DeleteMeasurement(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("succesfully deleted measurement %v", id)})
}

/*
curl -X PUT http://localhost:8080/measurements/1 \
     -H "Content-Type: application/json" \
     -d '{"unit": "volt"}'
*/

func (h *Handler) MeasurementUpdate(c *gin.Context) {
	id, err := util.GetParamInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateData map[string]any
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateMeasurement(id, updateData); err != nil {
		if err.Error() == "record not found" {
			// Return 404 if record not found
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			//Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated measurement"})
}
