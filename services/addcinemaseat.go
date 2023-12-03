package services

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateCinemaSeatRequest struct {
	Id           string            `param:"Id"`
	CinemaHallId string            `param:"cinemaHallId"`
	Seats        []CinemaSeatModel `json:"cinemaSeats"`
}

type CreateCinemaSeatResponse struct {
}

type CinemaHallDTO struct {
	CinemaId     string
	CinemaHallId string
	TotalSeat    int
}

type CinemaSeatDTO struct {
	CinemaHallId string
	SeatNumber   int
}

func (cinemaService CinemaService) AddCinemaSeat(request CreateCinemaSeatRequest) (CreateCinemaSeatResponse, models.ErrorResponse) {
	var err error
	cinemaService.Logger.Info().Interface("createCinemaSeatRequest", request).Msg("request")
	validationErrors := validateCinemSeatRequiredFields(request)
	if len(validationErrors) > 0 {
		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: validationErrors, StatusCode: http.StatusBadRequest}
	}

	seatNumbers := []int{}
	if len(request.Seats) > 1 {
		for _, seat := range request.Seats {
			seatNumbers = append(seatNumbers, seat.SeatNumber)
		}

		duplicateSeatNumbers := lo.FindDuplicates(seatNumbers)
		if len(duplicateSeatNumbers) > 0 {
			err = errors.New("seat number in the request contains duplicates")
			cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
			return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
		}
	}

	cinemaQuery, err := cinemaService.DB.Table("cinemas").
		Where("cinemas.Id = ?", request.Id).
		Where("cinemas.IsDeprecated = ?", false).
		Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
		Where("cinemaHalls.Id = ?", request.CinemaHallId).
		Where("cinemaHalls.IsDeprecated = ?", false).
		Select("cinemas.Id AS CinemaId, cinemaHalls.Id AS CinemaHallId, cinemaHalls.TotalSeat AS TotalSeat").
		Rows()

	if err != nil {
		cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
	}

	defer cinemaQuery.Close()

	var existingHalls CinemaHallDTO
	i := 0
	for cinemaQuery.Next() {
		if i > 1 {
			break
		}
		err = cinemaQuery.Scan(&existingHalls.CinemaId, &existingHalls.CinemaHallId, &existingHalls.TotalSeat)
		if err != nil {
			cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
			return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
		}
		i++
	}

	if reflect.ValueOf(existingHalls).IsZero() {
		err = errors.New("cinema info not found")
		cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
	}

	if existingHalls.TotalSeat < len(request.Seats) {
		err = fmt.Errorf("total number of cinema seats in the system is less that the new seats to add")
		cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")

		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err},
			StatusCode: http.StatusBadRequest}
	}

	cinemaHallQuery, err := cinemaService.DB.Table("cinemaHalls").
		Where("cinemaHalls.Id = ?", existingHalls.CinemaHallId).
		Where("cinemaHalls.IsDeprecated = ?", false).
		Joins("join cinemaSeats on cinemaHalls.Id = cinemaSeats.CinemaHallId").
		Where("cinemaSeats.IsDeprecated = ?", false).
		Select("cinemaHalls.Id AS CinemaHallId, cinemaSeats.SeatNumber AS SeatNumber").
		Rows()

	if err != nil {
		cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
	}

	defer cinemaHallQuery.Close()

	var existingSeats []CinemaSeatDTO
	for cinemaHallQuery.Next() {
		var cinemaSeatDTO CinemaSeatDTO
		err = cinemaHallQuery.Scan(&cinemaSeatDTO.CinemaHallId, &cinemaSeatDTO.SeatNumber)
		if err != nil {
			return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
		}
		existingSeats = append(existingSeats, cinemaSeatDTO)
	}

	if len(existingSeats)+len(request.Seats) > existingHalls.TotalSeat {
		err = fmt.Errorf("total number of cinema seats in the system is less that the new seats to add")
		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err},
			StatusCode: http.StatusBadRequest}
	}

	//check for duplicates
	for _, exSeat := range existingSeats {
		for _, seat := range request.Seats {
			if exSeat.SeatNumber == seat.SeatNumber {
				err = errors.New("seat number already exist in the system")
				cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
				return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
			}
		}
	}

	//add new seats
	err = cinemaService.DB.Transaction(func(tx *gorm.DB) error {
		if len(request.Seats) > 0 {
			for _, seat := range request.Seats {
				cinemaSeat := entities.CinemaSeat{
					Id:           sequentialguid.New().String(),
					SeatNumber:   seat.SeatNumber,
					Type:         seat.Type,
					IsDeprecated: false,
					CinemaHallId: request.CinemaHallId,
				}

				if err = tx.Create(&cinemaSeat).Error; err != nil {
					// return any error will rollback
					return err
				}
			}

		}
		return nil
	})

	if err != nil {
		cinemaService.Logger.Info().Interface("createCinemaSeatResponse", err.Error()).Msg("response")
		return CreateCinemaSeatResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
	}

	resp := CreateCinemaSeatResponse{}
	cinemaService.Logger.Info().Interface("createCinemaSeatResponse", resp).Msg("response")
	return resp, models.ErrorResponse{}
}

func validateCinemSeatRequiredFields(request CreateCinemaSeatRequest) []error {
	validationErrors := []error{}
	if len(request.Id) == 0 || len(request.Id) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "cinemaId"))
	}

	if request.Id == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "cinemaId"))
	}

	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "cinemahallId"))
	}

	if request.CinemaHallId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "cinemahallId"))
	}

	for i, seat := range request.Seats {
		if seat.SeatNumber <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("CinemaSeat[%d].SeatNumber should not have be zero or negative number", i))
		}

		if seat.Type <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("CinemaSeat[%d].Type should not have be zero or negative number", i))
		}

	}

	return validationErrors
}
