package controllers

import (
	"fmt"
	"intelligent-analysis-of-car-sensors-backend/ai"
	"intelligent-analysis-of-car-sensors-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AICtrl is the entity that manages PCA, k-means and SVM analysis requests for a dataset.
type AICtrl struct {
}

// NewAICtrl returns a new instance of AICtrl.
func NewAICtrl() *AICtrl {
	return &AICtrl{}
}

// Request is the entity that store the parameters used in PCA, k-means and SVM.
type Request struct {
	ClustersNumber   string `form:"clusters-number"`
	ComponentsNumber string `form:"components-number"`
}

// Classify process a dataset applying principal components analysis and store the results.
func (p *AICtrl) Classify(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	dataset := models.Dataset{}
	db.Preload("Logs.Records").First(&dataset, c.Param("id"))

	var request Request
	err := c.Bind(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request could no be binded :: %w", err))
		return
	}

	// TODO: remove my local adress
	//client := ai.NewClient("http://172.18.0.1:5000")
	client := ai.NewClient("http://localhost:5000")

	result, err := client.Start(&ai.ClientRequest{
		Dataset:          &dataset,
		ClustersNumber:   request.ClustersNumber,
		ComponentsNumber: request.ComponentsNumber})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request  to ai service failed :: %w", err))
		return
	}

	// Save results
	kMeans := models.Kmeans{
		TwoFirstComponentsPlot:               result.TwoFirstComponentsPlot,
		ComponentsAndFeaturesPlot:            result.ComponentsAndFeaturesPlot,
		ExplainedVarianceRatio:               result.ExplainedVarianceRatio,
		WCSSPlot:                             result.WCSSPlot,
		CumulativeExplainedVarianceRatioPlot: result.CumulativeExplainedVarianceRatioPlot,
		ClusterList:                          result.ClusterList,
		MoreImportantFeatures:                result.MoreImportantFeatures}

	svm := models.SVM{TwoFirstComponentsPlot: result.SVMPlot}

	dataset.ClassificationApplied = true
	dataset.KMeansResult = kMeans
	dataset.SVMResult = svm
	db.Save(&dataset)

	c.JSON(http.StatusOK, dataset)
}
