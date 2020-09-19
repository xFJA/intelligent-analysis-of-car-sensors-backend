package models

// Prediction represents the entity that store all the results after the LSTM prediction.
type Prediction struct {
	ID                        uint   `gorm:"primary_key" json:"id"`
	LearningCurvePlot         string `json:"learningCurvePlot"`
	PredictionPlot            string `json:"predictionPlot"`
	Time                      string `json:"time"`
	RMSE                      string `json:"rmse"`
	Feature                   string `json:"feature"`
	Epochs                    int    `json:"epochs"`
	PredictionFeaturesType    string `json:"predictionFeaturesType"`
	PrincipalComponentsNumber int    `json:"principalComponentsNumber"`
}
