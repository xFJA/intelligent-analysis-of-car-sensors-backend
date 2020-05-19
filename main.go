package main

import (
	"intelligent-analysis-of-car-sensors-backend/controllers"
	"intelligent-analysis-of-car-sensors-backend/models"
	"intelligent-analysis-of-car-sensors-backend/store"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Setup the database
	db, err := models.SetupModels()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Pass db instance to controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// TODO: Investigate official way of setting seeds
	// TODO: load path from config file
	// Add seeds for sensors information
	sensorsInformation, err := store.NewSensorsInformation("store/sensors_information.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, sensor := range sensorsInformation.SensorList {
		db.Create(&models.Sensor{
			PID:         sensor.PID,
			Description: sensor.Description,
			MeasureUnit: sensor.MeasureUnit,
		})
	}

	// Setup controllers
	csvStore := store.NewCSVStore(db)
	datasetsController := controllers.NewDatasetsCtrl(csvStore)

	// Setup endpoints
	r.GET("/datasets", datasetsController.GetDatasets)
	r.POST("/datasets", datasetsController.AddDataset)
	r.GET("/datasets/:id", datasetsController.GetDataset)
	r.DELETE("/datasets/:id", datasetsController.DeleteDataset)

	r.Run()
}
