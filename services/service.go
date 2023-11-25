package services

import (
	"github.com/nats-io/nats.go"
	"github.com/Wolechacho/ticketmaster-backend/tools"
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
	DB *gorm.DB
	Nc *nats.Conn
}

type AuthService struct {
	DB *gorm.DB
}

type BookingService struct {
	DB         *gorm.DB
	PDFService tools.PDFService
}
