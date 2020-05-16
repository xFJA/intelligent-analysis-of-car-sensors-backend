package main

import (
	"intelligent-analysis-of-car-sensors-backend/models"
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

	// Pass db instance to controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Run()
}
