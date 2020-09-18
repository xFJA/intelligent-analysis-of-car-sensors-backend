package models

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
)

// Dataset represents the entity that store all datasets.
type Dataset struct {
	ID                    uint       `gorm:"primary_key" json:"id"`
	Date                  int64      `json:"date"`
	Name                  string     `json:"name"`
	Logs                  []Log      `gorm:"foreignkey:DatasetID" json:"logs"`
	RowsNumber            int        `json:"rowsNumber"`
	ColumnNames           string     `json:"columnNames"`
	ClassificationApplied bool       `json:"classificationApplied"`
	KMeansResult          Kmeans     `gorm:"foreignkey:id" json:"kmeansResult"`
	SVMResult             SVM        `gorm:"foreignkey:id" json:"svmResult"`
	Prediction            Prediction `gorm:"foreignkey:id" json:"prediction"`
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

// CreateCSVFromDatasetEntitySVM returns a csv file created from a Dataset entity plus
// the dataset given.
func CreateCSVFromDatasetEntitySVM(dataset *Dataset, givenDataset *multipart.FileHeader) (string, error) {
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

	// Add dataset given
	// Open file
	src, err := givenDataset.Open()
	if err != nil {
		return "", err
	}

	// Read csv
	reader := csv.NewReader(bufio.NewReader(src))
	header, err := reader.Read()
	if err != nil {
		return "", err
	}
	// Add headers from given dataset
	csvBuilder = append(csvBuilder, headers)

	for index, value := range header {
		header[index] = cleanString(value)
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		newValues := []string{}

		for _, value := range line[1 : len(line)-1] {
			parsedValue, err := parseFloat(value)
			if err != nil {
				return "", err
			}
			newValues = append(newValues, fmt.Sprintf("%f", parsedValue))
		}

		csvBuilder = append(csvBuilder, newValues)
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

func cleanString(value string) string {
	// Remove parentheses
	res := strings.Trim(value, "()")

	// Remove blank spaces
	res = strings.TrimSpace(res)

	return res
}

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(value), 64)
}
