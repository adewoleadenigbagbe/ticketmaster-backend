package db

import (
	"github.com/fatih/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase() (*gorm.DB, error) {
	var err error
	dsn := "root:P@ssw0r1d@tcp(127.0.0.1:3306)/ticketmasterDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	color.Blue("Connected to the Database")

	return db, nil
}
