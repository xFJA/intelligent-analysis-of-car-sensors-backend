package models

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// Dataset represents the entity that store all datasets.
type Dataset struct {
	ID                    uint   `gorm:"primary_key" json:"id"`
	Date                  int64  `json:"date"`
	Name                  string `json:"name"`
	Logs                  []Log  `gorm:"foreignkey:DatasetID" json:"logs"`
	RowsNumber            int    `json:"rowsNumber"`
	ColumnNames           string `json:"columnNames"`
	ClassificationApplied bool   `json:"classificationApplied"`
	KMeansResult          Kmeans `gorm:"foreignkey:id" json:"kmeansResult"`
	SVMResult             SVM    `gorm:"foreignkey:id" json:"svmResult"`
}

// CreateCSVFromDatasetEntity returns a csv file created from a Dataset entity.
func CreateCSVFromDatasetEntity(dataset *Dataset) (string, error) {
	csvBuilder := [][]string{}

	// Create headers column with each feature name
	headers := []string{}
	for _, record := range dataset.Logs[0].Records {
		headers = append(headers, record.SensorPID)
	}

	var labelList []int
	if dataset.KMeansResult.ClusterList != "" {
		headers = append(headers, "LABEL")
		err := json.Unmarshal([]byte(dataset.KMeansResult.ClusterList), &labelList)
		if err != nil {
			return "", err
		}
	}

	csvBuilder = append(csvBuilder, headers)

	// Add each feature value
	for index, log := range dataset.Logs {
		values := []string{}
		for _, record := range log.Records {
			values = append(values, fmt.Sprintf("%f", record.Value))
		}

		if labelList != nil {
			values = append(values, fmt.Sprintf("%d", labelList[index]))
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
