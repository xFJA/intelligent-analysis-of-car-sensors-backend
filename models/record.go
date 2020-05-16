package models

import "github.com/jinzhu/gorm"

// Record represents the entity that store the value of a sensor.
type Record struct {
	gorm.Model
	LogID  uint
	Value  float64
	Sensor Sensor
}
