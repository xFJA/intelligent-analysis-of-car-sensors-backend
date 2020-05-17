package controllers

import (
	"intelligent-analysis-of-car-sensors-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetDatasets returns all datasets from the database.
func GetDatasets(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datasets []models.Dataset
	db.Preload("Logs.Records").Find(&datasets)

	c.JSON(http.StatusOK, datasets)
}
