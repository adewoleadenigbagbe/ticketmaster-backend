package db

import (
	"os"

	"github.com/Wolechacho/ticketmaster-backend/shared/models"
	"github.com/fatih/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase(path string) (*gorm.DB, error) {
	var err error

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	dbConfig, err := models.CreateDbConfig(content)
	if err != nil {
		return nil, err
	}

	dsn := dbConfig.GetDsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	color.Blue("Connected to the Database")

	return db, nil
}
