package main

import (
	"intelligent-analysis-of-car-sensors-backend/controllers"
	"intelligent-analysis-of-car-sensors-backend/models"
	"intelligent-analysis-of-car-sensors-backend/store"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(ErrorHandler)

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

	// Setup controllers
	csvStore := store.NewCSVStore()
	datasetsController := controllers.NewDatasetsCtrl(csvStore)
	aiController := controllers.NewAICtrl()

	// Setup endpoints
	r.GET("/datasets", datasetsController.GetDatasets)
	r.POST("/datasets", datasetsController.AddDataset)
	r.GET("/datasets/:id", datasetsController.GetDataset)
	r.DELETE("/datasets/:id", datasetsController.DeleteDataset)
	r.GET("/datasets/:id/csv", datasetsController.GetCSVFile)

	r.GET("/classify/:id", aiController.Classify)
	r.POST("/classify-svm", aiController.ClassifySVM)

	r.Run()
}
