package controllers

import (
	"bufio"
	"encoding/csv"
	"intelligent-analysis-of-car-sensors-backend/models"
	"intelligent-analysis-of-car-sensors-backend/store"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DatasetsCtrl is the entity that manages all datasets requests.
type DatasetsCtrl struct {
	csvStore *store.CSVStore
}

// NewDatasetsCtrl returns a new instance of DatasetsCtrl.
func NewDatasetsCtrl(csvStore *store.CSVStore) *DatasetsCtrl {
	return &DatasetsCtrl{csvStore: csvStore}
}

// GetDatasets returns all datasets.
func (d *DatasetsCtrl) GetDatasets(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datasets []models.Dataset
	db.Preload("Logs.Records").Find(&datasets)

	c.JSON(http.StatusOK, gin.H{"data": datasets})
}

// AddDataset creates a new dataset.
func (d *DatasetsCtrl) AddDataset(c *gin.Context) {
	// Get csv file from POST form
	csvFile, err := c.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file could not be found on POST form"})
		return
	}

	// Open file
	src, err := csvFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file could not be opened"})
		return
	}

	// Read csv
	reader := csv.NewReader(bufio.NewReader(src))
	dataset, err := d.csvStore.Load(reader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file could not be stored"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

// GetDataset returns a single dataset.
func (d *DatasetsCtrl) GetDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataset models.Dataset
	if err := db.Preload("Logs.Records").Where("id = ?", c.Param("id")).First(&dataset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset could not be found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

// DeleteDataset removes a dataset.
func (d *DatasetsCtrl) DeleteDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataset models.Dataset
	if err := db.Where("id = ?", c.Param("id")).First(&dataset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset could not be found"})
		return
	}

	db.Delete(&dataset)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
