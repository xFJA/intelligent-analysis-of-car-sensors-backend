package store

import (
	"encoding/csv"
	"intelligent-analysis-of-car-sensors-backend/models"
	"io"
	"os"
	"strconv"
	"strings"

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

// Load reads a csv file and store the data in the database.
func (s *CSVStore) Load(path string) error {
	// Open file
	csvfile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	// Parse file
	data := csv.NewReader(csvfile)

	// Store data
	header, err := data.Read()
	if err != nil {
		return err
	}

	for index, value := range header {
		header[index] = cleanString(value)
	}

	for {
		line, err := data.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// Create log
		time, err := parseFloat(line[0])
		if err != nil {
			return err
		}
		log := models.Log{
			Time:    time,
			Records: []models.Record{}}

		s.db.Create(&log)

		for index, value := range line[1 : len(line)-1] {
			// Create record entity
			parsedValue, err := parseFloat(value)
			if err != nil {
				return err
			}
			record := models.Record{
				Value:     parsedValue,
				SensorPID: header[index+1],
			}

			// Add record to the log
			log.Records = append(log.Records, record)
		}

		s.db.Save(&log)
	}

	return nil
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
