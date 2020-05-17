package models

import (
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

	db.DropTableIfExists(
		&Dataset{},
		&Log{},
		&Record{},
		&Sensor{})

	db.AutoMigrate(
		&Dataset{},
		&Log{},
		&Record{},
		&Sensor{})

	return db, nil
}
