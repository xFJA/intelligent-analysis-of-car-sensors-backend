package models

// SVM represents the entity that store the results from SVM.
type SVM struct {
	ID                     uint   `gorm:"primary_key" json:"id"`
	TwoFirstComponentsPlot string `json:"twoFirstComponentsPlot"`
}
