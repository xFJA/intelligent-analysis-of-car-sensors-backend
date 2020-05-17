package models

// Dataset represents the entity that store all datasets.
type Dataset struct {
	ID   uint  `gorm:"primary_key" json:"id"`
	Date int64 `json:"date"`
	Logs []Log `gorm:"foreignkey:DatasetID" json:"logs"`
}
