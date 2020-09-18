package controllers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"intelligent-analysis-of-car-sensors-backend/ai"
	"intelligent-analysis-of-car-sensors-backend/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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

// PredictionRequest is the entity that store the parameters used in the prediction.
type PredictionRequest struct {
	Feature string `form:"feature"`
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

// ClassifySVM classify a given dataset using the dataset which id is given as training data.
func (p *AICtrl) ClassifySVM(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	dataset := models.Dataset{}
	db.Preload("Logs.Records").Preload("KMeansResult").First(&dataset, c.PostForm("id"))

	// Get csv file from POST form
	csvFile, err := c.FormFile("csv")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could not be found on POST form :: %w", err))
		return
	}

	client := ai.NewClient("http://localhost:5000")

	result, err := client.ClassifySVM(&ai.ClassifySVMRequest{
		Dataset:           &dataset,
		DatasetToClassify: csvFile,
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request  to ai service failed :: %w", err))
		return
	}
	// Parse classification list to string array
	classificationList := strings.FieldsFunc(result.ClassificationList, split)

	// Add classification list as LABEL column for the dataset given
	csvBuilder := [][]string{}
	src, err := csvFile.Open()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could not be opened :: %w", err))
		return
	}

	reader := csv.NewReader(bufio.NewReader(src))
	headers, err := reader.Read()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could no be readed :: %w", err))
		return
	}
	headers = append(headers, "LABEL (SVM)")

	// Add headers from given dataset
	csvBuilder = append(csvBuilder, headers)

	i := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could no be readed :: %w", err))
			return
		}

		line = append(line, classificationList[i])

		csvBuilder = append(csvBuilder, line)
		i++
	}

	// Create csv writer
	fileNameComplete := csvFile.Filename + ".csv"
	file, err := os.Create(fileNameComplete)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file could not be created :: %w", err))
		return
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(csvBuilder)
	csvWriter.Flush()

	csvFileContent, err := ioutil.ReadFile(fileNameComplete)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("CSV file content could not be obtained :: %w", err))
		return
	}
	defer os.Remove(fileNameComplete)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileNameComplete))

	c.Data(http.StatusOK, "text/csv", csvFileContent)
}

// Predict apply predictions using LSTM neuronal networks over a dataset.
func (p *AICtrl) Predict(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	dataset := models.Dataset{}
	db.Preload("Logs.Records").First(&dataset, c.Param("id"))

	var request PredictionRequest
	err := c.Bind(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request could no be binded :: %w", err))
		return
	}

	client := ai.NewClient("http://localhost:5000")

	result, err := client.Prediction(&ai.PredictionRequest{
		Dataset: &dataset,
		Feature: request.Feature,
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request  to ai service failed :: %w", err))
		return
	}

	dataset.Prediction = models.Prediction{
		LearningCurvePlot: result.LearningCurvePlot,
		PredictionPlot:    result.PredictionPlot,
		RMSE:              result.RMSE,
		Time:              result.Time,
		Feature:           request.Feature,
	}
	dataset.PredictionApplied = true
	db.Save(&dataset)

	c.JSON(http.StatusOK, dataset)
}

func split(r rune) bool {
	return r == '[' || r == ']' || r == '"' || r == ','
}
