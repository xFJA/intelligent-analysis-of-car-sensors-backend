package models

// Record represents the entity that store the value of a sensor.
type Record struct {
	ID        uint    `gorm:"primary_key" json:"id"`
	Value     float64 `json:"value"`
	LogID     uint    `json:"-"`
	SensorPID string  `json:"sensorPID"`
}
