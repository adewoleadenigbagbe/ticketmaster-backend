package services

import "gorm.io/gorm"

//TODO: add all these to the service response
type ErrrorResponse struct {
	Errors     []error
	StatusCode int
}

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
