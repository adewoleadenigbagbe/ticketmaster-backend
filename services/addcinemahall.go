package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CinemaHallRequest struct {
	Id    string            `param:"Id"`
	Halls []CinemaHallModel `json:"halls"`
}

type CinemaHallResponse struct {
	CinemaId string `json:"CinemaId"`
}

func (cinemaService CinemaService) AddCinemaHall(request CinemaHallRequest) (CinemaHallResponse, models.ErrorResponse) {
	var err error
	var errs []error

	cinemaService.Logger.Info().Interface("cinemaHallRequest", request).Msg("request")
	//validate request
	errs = validateCinemaHallRequiredFields(request)
	if len(errs) > 0 {
		return CinemaHallResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: errs}
	}

	//check for duplicate hall names
	hallNames := []string{}
	for _, hall := range request.Halls {
		hallNames = append(hallNames, hall.Name)
	}

	duplicateHallNames := lo.FindDuplicates(hallNames)
	if len(duplicateHallNames) > 0 {
		errs = append(errs, fmt.Errorf("should not have duplicate hall names : %s", strings.Join(duplicateHallNames, ",")))
		errResp := models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: errs}
		cinemaService.Logger.Info().Interface("cinemaHallResponse", errResp).Msg("response")
		return CinemaHallResponse{}, errResp
	}

	//check for duplicate seat number
	for _, hall := range request.Halls {
		seatNumbers := []int{}
		for _, seat := range hall.CinemaSeats {
			seatNumbers = append(seatNumbers, seat.SeatNumber)
		}

		duplicateSeatNumbers := lo.FindDuplicates(seatNumbers)
		if len(duplicateSeatNumbers) > 0 {
			errs = append(errs, fmt.Errorf("should not have duplicate seat numbers for hall name: %s", hall.Name))
			errResp := models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: errs}
			cinemaService.Logger.Info().Interface("cinemaHallResponse", errResp).Msg("response")
			return CinemaHallResponse{}, errResp
		}
	}

	cinema := &entities.Cinema{}
	result := cinemaService.DB.Where("Id = ? AND IsDeprecated = ?", request.Id, false).First(cinema)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		errs = append(errs, errors.New("cinema record not found"))
		errResp := models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: errs}
		cinemaService.Logger.Info().Interface("cinemaHallResponse", errResp).Msg("response")
		return CinemaHallResponse{}, errResp
	}

	//check if the not duplicate names in the DB
	var countResult int64
	cinemaService.DB.Model(&entities.CinemaHall{}).Where("Name IN ? AND CinemaId = ? AND IsDeprecated = ?", hallNames, cinema.Id, false).Count(&countResult)
	if countResult > 0 {
		errs = append(errs, fmt.Errorf(("cinemaHall name already exist in system")))
		errResp := models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: errs}
		cinemaService.Logger.Info().Interface("cinemaHallResponse", errResp).Msg("response")
		return CinemaHallResponse{}, errResp
	}

	err = cinemaService.DB.Transaction(func(tx *gorm.DB) error {
		if len(request.Halls) != 0 {
			for _, hall := range request.Halls {
				cinemaHall := entities.CinemaHall{
					Id:           sequentialguid.New().String(),
					Name:         hall.Name,
					TotalSeat:    hall.TotalSeat,
					CinemaId:     request.Id,
					IsDeprecated: false,
				}

				if err = tx.Create(&cinemaHall).Error; err != nil {
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
		errs = append(errs, result.Error)
		errResp := models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: errs}
		cinemaService.Logger.Info().Interface("cinemaHallResponse", err.Error()).Msg("response")

		return CinemaHallResponse{}, errResp
	}

	resp := CinemaHallResponse{
		CinemaId: request.Id,
	}
	cinemaService.Logger.Info().Interface("cinemaHallResponse", resp).Msg("response")
	return resp, models.ErrorResponse{}
}

func validateCinemaHallRequiredFields(request CinemaHallRequest) []error {
	validationErrors := []error{}
	if len(request.Id) == 0 || len(request.Id) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "cinemaId"))
	}

	if request.Id == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "cinemaId"))
	}

	for i, hall := range request.Halls {
		if hall.Name == "" {
			validationErrors = append(validationErrors, fmt.Errorf("cinemahall[%d].Name should be supplied", i))
		}

		if hall.TotalSeat <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("cinemahall[%d].TotalSeat should not be lesss than or equal to zero", i))
		}

		if len(hall.CinemaSeats) > hall.TotalSeat {
			validationErrors = append(validationErrors, fmt.Errorf("total Seat: %d is less than the number of seats for a hall: %d to be inserted", hall.TotalSeat, len(hall.CinemaSeats)))
		}

		for j, seat := range hall.CinemaSeats {
			if seat.SeatNumber <= 0 {
				validationErrors = append(validationErrors, fmt.Errorf("Halls[%d].CinemaSeat[%d].SeatNumber should not have a negative number", i, j))
			}

			if seat.Type <= 0 {
				validationErrors = append(validationErrors, fmt.Errorf("Halls[%d].CinemaSeat[%d].Type should not have a negative number", i, j))
			}
		}
	}

	return validationErrors
}
