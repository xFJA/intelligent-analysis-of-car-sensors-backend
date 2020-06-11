package controllers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"intelligent-analysis-of-car-sensors-backend/models"
	"intelligent-analysis-of-car-sensors-backend/store"
	"net/http"
	"strings"

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
	db.Find(&datasets)

	c.JSON(http.StatusOK, gin.H{"data": datasets})
}

// AddDataset creates a new dataset.
func (d *DatasetsCtrl) AddDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get csv file from POST form
	csvFile, err := c.FormFile("csv")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could not be found on POST form :: %w", err))
		return
	}

	// Open file
	src, err := csvFile.Open()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could not be opened :: %w", err))
		return
	}

	// Read csv
	reader := csv.NewReader(bufio.NewReader(src))
	dataset, err := d.csvStore.Load(reader, csvFile.Filename)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could not be stored :: %w", err))
		return
	}

	// Store dataset
	// TODO: check if these values could change
	dataset.RowsNumber = len(dataset.Logs)
	dataset.ColumnNames = getColumnNames(dataset)
	db.Create(&dataset)

	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

// GetDataset returns a single dataset.
func (d *DatasetsCtrl) GetDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataset models.Dataset
	if err := db.Preload("Logs.Records").Where("id = ?", c.Param("id")).First(&dataset).Error; err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Dataset could not be found :: %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

// DeleteDataset removes a dataset.
func (d *DatasetsCtrl) DeleteDataset(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataset models.Dataset
	if err := db.Where("id = ?", c.Param("id")).First(&dataset).Error; err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Dataset could not be found :: %w", err))
		return
	}

	db.Delete(&dataset)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// getColumnNames returns a list of the column names of the dataset.
func getColumnNames(dataset *models.Dataset) string {
	var names string

	// We can get them only checking the first log
	for _, record := range dataset.Logs[0].Records {
		names += record.SensorPID + ", "
	}
	strings.TrimSuffix(names, ",")

	return names
}
