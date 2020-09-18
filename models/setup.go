package models

import (
	"intelligent-analysis-of-car-sensors-backend/utils"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
)

// TODO: load db string from config file
var dbstring string = "host=localhost user=postgres dbname=postgres sslmode=disable password=postgres"

// SetupModels create a connection with the database and migrate the models schema.
func SetupModels() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbstring)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&Dataset{},
		&Log{},
		&Record{},
		&Sensor{},
		&Kmeans{},
		&SVM{},
		&Prediction{})

	// TODO: Investigate official way of setting seeds
	// TODO: load path from config file
	// Add seeds for sensors information
	sensorsInformation, err := utils.NewSensorsInformation("utils/sensors_information.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, sensor := range sensorsInformation.SensorList {
		db.Create(&Sensor{
			PID:         sensor.PID,
			Description: sensor.Description,
			MeasureUnit: sensor.MeasureUnit,
		})
	}

	return db, nil
}
