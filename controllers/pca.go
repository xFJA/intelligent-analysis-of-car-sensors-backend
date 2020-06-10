package controllers

import (
	"fmt"
	"intelligent-analysis-of-car-sensors-backend/models"
	"intelligent-analysis-of-car-sensors-backend/pca"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// PCACtrl is the entity that manages the pincipal components analysis requests for a dataset.
type PCACtrl struct {
}

// NewPCACtrl returns a new instance of PCACtrl.
func NewPCACtrl() *PCACtrl {
	return &PCACtrl{}
}

// PCA process a dataset applying principal components analysis and store the results.
func (p *PCACtrl) PCA(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	dataset := models.Dataset{}
	db.Preload("Logs.Records").First(&dataset, c.Param("id"))

	// TODO: remove my local adress
	//pcaClient := pca.NewClient("http://172.18.0.1:5000")
	pcaClient := pca.NewClient("http://localhost:5000")

	pcaResult, err := pcaClient.PCA(&pca.ClientRequest{Dataset: &dataset})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("PCA request failed :: %w", err))
		return
	}

	// Add PCA results to dataset
	dataset.TwoFirstComponentsPlot = pcaResult.TwoFirstComponentsPlot
	dataset.ComponentsAndFeaturesPlot = pcaResult.ComponentsAndFeaturesPlot
	dataset.ExplainedVarianceRatio = pcaResult.ExplainedVarianceRatio

	db.Save(&dataset)

	c.JSON(http.StatusOK, dataset)
}
