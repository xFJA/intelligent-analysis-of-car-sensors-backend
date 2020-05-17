package controllers

import (
	"intelligent-analysis-of-car-sensors-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetDatasets returns all datasets.
func GetDatasets(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datasets []models.Dataset
	db.Preload("Logs.Records").Find(&datasets)

	c.JSON(http.StatusOK, gin.H{"data": datasets})
}

// GetDataset returns a single dataset.
func GetDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataset models.Dataset
	if err := db.Preload("Logs.Records").Where("id = ?", c.Param("id")).First(&dataset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset could not be found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

// DeleteDataset removes a dataset.
func DeleteDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataset models.Dataset
	if err := db.Where("id = ?", c.Param("id")).First(&dataset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset could not be found"})
		return
	}

	db.Delete(&dataset)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
