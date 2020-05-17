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

	// TODO: load csv from HTTP requests
	csvStore := store.NewCSVStore(db)
	err = csvStore.Load("live1_short.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Setup endpoints
	r.GET("/datasets", controllers.GetDatasets)
	r.GET("/datasets/:id", controllers.GetDataset)
	r.DELETE("/datasets/:id", controllers.DeleteDataset)

	r.Run()
}
