package services

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CinemaService struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

type CityService struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

type MovieService struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

type ShowService struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

type UserService struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

type AuthService struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}
