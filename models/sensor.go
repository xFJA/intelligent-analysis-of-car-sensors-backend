package models

// Sensor represents the entity that store the information about a car sensor.
type Sensor struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	PID         string `gorm:"unique_index" json:"pid"`
	Description string `json:"description"`
	MeasureUnit string `json:"measureUnit"`
	Record      Record `gorm:"foreignkey:SensorPID" json:"-"`
}
