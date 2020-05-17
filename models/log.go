package models

// Log represents the entity that store sensors records for a certain instance of time.
type Log struct {
	ID        uint     `gorm:"primary_key" json:"id"`
	Time      float64  `json:"time"`
	Records   []Record `gorm:"foreignkey:LogID" json:"records"`
	DatasetID uint     `json:"-"`
}
