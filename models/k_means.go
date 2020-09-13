package models

// Kmeans represents the entity that store the results from k-means.
type Kmeans struct {
	ID                                   uint   `gorm:"primary_key" json:"id"`
	TwoFirstComponentsPlot               string `json:"twoFirstComponentsPlot"`
	ComponentsAndFeaturesPlot            string `json:"componentsAndFeaturesPlot"`
	ExplainedVarianceRatio               string `json:"explainedVarianceRatio"`
	WCSSPlot                             string `json:"wcssPlot"`
	CumulativeExplainedVarianceRatioPlot string `json:"cumulativeExplainedVarianceRatioPlot"`
	ClusterList                          string `gorm:"DEFAULT:false" json:"clusterList"`
	MoreImportantFeatures                string `json:"moreImportantFeatures"`
}
