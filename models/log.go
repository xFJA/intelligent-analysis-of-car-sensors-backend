package models

import "github.com/jinzhu/gorm"

// Log represents the entity that store sensors records for a certain instance of time.
type Log struct {
	gorm.Model
	Time    float64
	Records []Record
}
