package controllers

import (
	"fmt"
	"net/http"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CinemaController struct {
}

func (cinemaController CinemaController) CreateCinema(cinemaContext echo.Context) error {
	var err error

	request := new(createCinemaRequest)

	err = cinemaContext.Bind(request)
	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	fieldErrors := validateCinema(*request)
	if len(fieldErrors) != 0 {
		errors := []string{}
		for _, err = range fieldErrors {
			errors = append(errors, err.Error())
		}
		return cinemaContext.JSON(http.StatusBadRequest, errors)
	}

	cinema := entities.Cinema{
		Id:                sequentialguid.New().String(),
		Name:              request.Name,
		TotalCinemalHalls: request.TotalCinemalHalls,
		CityId:            request.CityId,
		IsDeprecated:      false,
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&cinema).Error; err != nil {
			// return any error will rollback
			return err
		}

		address := entities.Address{
			Id:          sequentialguid.New().String(),
			EntityId:    cinema.Id,
			AddressType: enums.Cinema,
			AddressLine: request.Address,
			CityId:      request.CityId,
			Coordinates: entities.Coordinate{
				Longitude: request.Longitude,
				Latitude:  request.Latitude,
			},
			IsDeprecated: false,
		}

		if err := tx.Create(&address).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil
	})

	if err != nil {
		return cinemaContext.JSON(http.StatusBadRequest, err.Error())
	}

	response := new(createCinemaResponse)
	response.CinemaId = cinema.Id
	return cinemaContext.JSON(http.StatusOK, response)
}

type createCinemaRequest struct {
	Name              string  `json:"name"`
	CityId            string  `json:"cityId"`
	TotalCinemalHalls int     `json:"totalCinemalHalls"`
	Address           string  `json:"address"`
	Longitude         float32 `json:"longitude"`
	Latitude          float32 `json:"latitude"`
}

type createCinemaResponse struct {
	CinemaId string `json:"CinemaId"`
}

func validateCinema(request createCinemaRequest) []error {
	var validationErrors []error

	if len(request.Name) == 0 {
		validationErrors = append(validationErrors, fmt.Errorf("name is a required field"))
	}

	if request.CityId == entities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("cityId should have a valid UUID"))
	}

	if len(request.CityId) == 0 || len(request.CityId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("cityId is a required field with 36 characters"))
	}

	if request.TotalCinemalHalls <= 0 {
		validationErrors = append(validationErrors, fmt.Errorf("totalCinemalHalls cannot be less than or equal to zero"))
	}

	return validationErrors
}
