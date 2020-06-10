package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// SensorsInformation represents all the information about all the sensors.
type SensorsInformation struct {
	SensorList []Sensor `json:"sensors"`
}

// Sensor represents the information of a sensor.
type Sensor struct {
	PID         string `json:"pid"`
	Description string `json:"description"`
	MeasureUnit string `json:"measurement_unit"`
}

// NewSensorsInformation returns a new instance of SensorInformation.
func NewSensorsInformation(path string) (*SensorsInformation, error) {
	var sensorsInformation SensorsInformation
	err := loadSensorsInformation(path, &sensorsInformation)
	if err != nil {
		return nil, err
	}
	return &sensorsInformation, nil
}

// loadSensorsInformation loads all the sensor information from a JSON file.
func loadSensorsInformation(path string, sensorsInformation *SensorsInformation) error {
	// Open file
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Parse data
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteValue, sensorsInformation)
}
