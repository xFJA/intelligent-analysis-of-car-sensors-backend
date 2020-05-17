package models

import "github.com/jinzhu/gorm"

// Sensor represents the entity that store the information about a car sensor.
type Sensor struct {
	gorm.Model
	PID         string `gorm:"unique_index"`
	Description string
	MeasureUnit string
	Record      Record `gorm:"foreignkey:SensorPID"`
}
