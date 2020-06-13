package models

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Dataset represents the entity that store all datasets.
type Dataset struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Date int64  `json:"date"`
	Name string `json:"name"`
	Logs []Log  `gorm:"foreignkey:DatasetID" json:"logs"`
	PCAResult
	RowsNumber  int    `json:"rowsNumber"`
	ColumnNames string `json:"columnNames"`
}

// PCAResult represents the analysis result from the PCA service.
type PCAResult struct {
	TwoFirstComponentsPlot    string `json:"twoFirstComponentsPlot"`
	ComponentsAndFeaturesPlot string `json:"componentsAndFeaturesPlot"`
	ExplainedVarianceRatio    string `json:"explainedVarianceRatio"`
	WCSSPlot                  string `json:"wcssPlot"`
}

// CreateCSVFromDatasetEntity returns a csv file created from a Dataset entity.
func CreateCSVFromDatasetEntity(dataset *Dataset) (string, error) {
	csvBuilder := [][]string{}

	// Create headers column with each feature name
	headers := []string{}
	for _, record := range dataset.Logs[0].Records {
		headers = append(headers, record.SensorPID)
	}
	csvBuilder = append(csvBuilder, headers)

	// Add each feature value
	for _, log := range dataset.Logs {
		values := []string{}
		for _, record := range log.Records {
			values = append(values, fmt.Sprintf("%f", record.Value))
		}
		csvBuilder = append(csvBuilder, values)
	}

	// Create csv writer
	fileNameComplete := dataset.Name + ".csv"
	file, err := os.Create(fileNameComplete)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(csvBuilder)
	csvWriter.Flush()

	return fileNameComplete, nil
}
