package controllers

import (
	"intelligent-analysis-of-car-sensors-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SensorsCtrl is the entity that manages all sensors requests.
type SensorsCtrl struct {
}

// NewSensorsCtrl returns a new instance of SensorsCtrl.
func NewSensorsCtrl() *SensorsCtrl {
	return &SensorsCtrl{}
}

// GetSensors returns all sensors information.
func (s *SensorsCtrl) GetSensors(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var sensors []models.Sensor
	db.Find(&sensors)

	c.JSON(http.StatusOK, gin.H{"data": sensors})
}
