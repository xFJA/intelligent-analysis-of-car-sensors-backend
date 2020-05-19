package store

import (
	"encoding/csv"
	"intelligent-analysis-of-car-sensors-backend/models"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// CSVStore represents the entity that manage all csv operations.
type CSVStore struct {
	db *gorm.DB
}

// NewCSVStore returns a new instance of CSVStore.
func NewCSVStore(db *gorm.DB) *CSVStore {
	return &CSVStore{db: db}
}

// Load reads a csv reader and store the data in the database.
func (s *CSVStore) Load(data *csv.Reader) (*models.Dataset, error) {
	// Store data
	header, err := data.Read()
	if err != nil {
		return nil, err
	}

	for index, value := range header {
		header[index] = cleanString(value)
	}

	// Create dataset entity
	dataset := models.Dataset{
		Date: time.Now().Unix(),
		Logs: []models.Log{},
	}

	for {
		line, err := data.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// Create log
		time, err := parseFloat(line[0])
		if err != nil {
			return nil, err
		}
		log := models.Log{
			Time:    time,
			Records: []models.Record{}}

		for index, value := range line[1 : len(line)-1] {
			// Create record entity
			parsedValue, err := parseFloat(value)
			if err != nil {
				return nil, err
			}
			record := models.Record{
				Value:     parsedValue,
				SensorPID: header[index+1],
			}

			// Add record to the log
			log.Records = append(log.Records, record)
		}

		dataset.Logs = append(dataset.Logs, log)
	}

	// Store dataset
	s.db.Create(&dataset)

	return &dataset, nil
}

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(value), 64)
}

func cleanString(value string) string {
	// Remove parentheses
	res := strings.Trim(value, "()")

	// Remove blank spaces
	res = strings.TrimSpace(res)

	return res
}

// LoadFromFile reads a csv file and store the data in the database.
func (s *CSVStore) LoadFromFile(path string) error {
	// Open file
	csvfile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	// Parse file
	data := csv.NewReader(csvfile)

	_, err = s.Load(data)

	return err
}
