package models

import "github.com/jinzhu/gorm"

// Record represents the entity that store the value of a sensor.
type Record struct {
	gorm.Model
	Value     float64
	LogID     uint
	SensorPID string
}
