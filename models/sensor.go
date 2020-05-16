package models

import "github.com/jinzhu/gorm"

// Sensor represents the entity that store the information about a car sensor.
type Sensor struct {
	gorm.Model
	PID         string
	Name        string
	MeasureUnit string
}
