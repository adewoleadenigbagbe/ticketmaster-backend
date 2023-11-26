package services

import (
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"gorm.io/gorm"
)

var (
	TIME_OVERLAP_ERROR   = "time overlap between the show Start and End Time"
	ErrInvalidUUID       = "%s should have a valid UUID"
	ErrRequiredUUIDField = "%s is a required field with 36 characters"
	ErrRequiredField     = "%s is a required field"
	ErrInValidField      = "%s supplied is invalid"
)

type CinemaHallModel struct {
	Name        string            `json:"name"`
	TotalSeat   int               `json:"totalSeat"`
	CinemaSeats []CinemaSeatModel `json:"cinemaSeats"`
}

type CinemaSeatModel struct {
	SeatNumber int            `json:"seatNumber"`
	Type       enums.SeatType `json:"type"`
}

type CreateCinemaRequest struct {
	Name              string            `json:"name"`
	CityId            string            `json:"cityId"`
	TotalCinemalHalls int               `json:"totalCinemalHalls"`
	Address           string            `json:"address"`
	Longitude         float32           `json:"longitude"`
	Latitude          float32           `json:"latitude"`
	Halls             []CinemaHallModel `json:"halls"`
}

type CreateCinemaResponse struct {
	CinemaId string `json:"CinemaId"`
}

func (cinemaService CinemaService) CreateCinema(request CreateCinemaRequest) (CreateCinemaResponse, models.ErrorResponse) {
	var err error
	fieldErrors := validateCinema(request)
	if len(fieldErrors) != 0 {
		return CreateCinemaResponse{}, models.ErrorResponse{Errors: fieldErrors, StatusCode: http.StatusBadRequest}
	}

	cinema := entities.Cinema{
		Id:                sequentialguid.New().String(),
		Name:              request.Name,
		TotalCinemalHalls: request.TotalCinemalHalls,
		CityId:            request.CityId,
		IsDeprecated:      false,
	}

	err = cinemaService.DB.Transaction(func(tx *gorm.DB) error {
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

		if len(request.Halls) != 0 {
			for _, hall := range request.Halls {
				cinemaHall := entities.CinemaHall{
					Id:           sequentialguid.New().String(),
					Name:         hall.Name,
					TotalSeat:    hall.TotalSeat,
					CinemaId:     cinema.Id,
					IsDeprecated: false,
				}

				if err := tx.Create(&cinemaHall).Error; err != nil {
					// return any error will rollback
					return err
				}

				if len(hall.CinemaSeats) != 0 {
					for _, seat := range hall.CinemaSeats {
						seat := entities.CinemaSeat{
							Id:           sequentialguid.New().String(),
							SeatNumber:   seat.SeatNumber,
							Type:         seat.Type,
							CinemaHallId: cinemaHall.Id,
							IsDeprecated: false,
						}

						if err := tx.Create(&seat).Error; err != nil {
							// return any error will rollback
							return err
						}
					}
				}

			}
		}

		return nil
	})

	if err != nil {
		return CreateCinemaResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}

	return CreateCinemaResponse{CinemaId: cinema.Id}, models.ErrorResponse{}
}

func validateCinema(request CreateCinemaRequest) []error {
	var validationErrors []error

	if len(request.Name) == 0 {
		validationErrors = append(validationErrors, fmt.Errorf("name is a required field"))
	}

	if request.CityId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("cityId should have a valid UUID"))
	}

	if len(request.CityId) == 0 || len(request.CityId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("cityId is a required field with 36 characters"))
	}

	if request.TotalCinemalHalls <= 0 {
		validationErrors = append(validationErrors, fmt.Errorf("totalCinemalHalls cannot be less than or equal to zero"))
	}

	if len(request.Halls) > request.TotalCinemalHalls {
		validationErrors = append(validationErrors, fmt.Errorf("length of the halls to add cannot be greater than the fixed hall in total"))
	}

	for i, hall := range request.Halls {
		if hall.TotalSeat <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("Halls[%d].TotalSeats cannot be less or equal to zero", i))
		}

		if len(hall.CinemaSeats) > hall.TotalSeat {
			validationErrors = append(validationErrors, fmt.Errorf("length of the cinema hall seats to add cannot be greater than the fixed cinema hall seats in total"))
		}
	}
	return validationErrors
}
