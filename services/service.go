package services

import (
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

type AuthService struct {
	DB *gorm.DB
}

type BookingService struct {
	DB         *gorm.DB
	PDFService tools.PDFService
}
