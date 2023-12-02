package services

import (
	"github.com/Wolechacho/ticketmaster-backend/tools"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type CinemaService struct {
	DB *gorm.DB
}

type CityService struct {
	DB *gorm.DB
}

type MovieService struct {
	DB *gorm.DB
}

type ShowService struct {
	DB *gorm.DB
}

type UserService struct {
	DB *gorm.DB
}

type BookService struct {
	DB         *gorm.DB
	Nc         *nats.Conn
	PDFService tools.PDFService
}

type AuthService struct {
	DB *gorm.DB
}
