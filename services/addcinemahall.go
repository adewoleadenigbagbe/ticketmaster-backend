package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CinemaHallRequest struct {
	CinemaId string            `json:"cinemaId"`
	Halls    []CinemaHallModel `json:"halls"`
}

type CinemaHallResponse struct {
	CinemaId string `json:"CinemaId"`
}

func (service CinemaService) AddCinemaHall(request CinemaHallRequest) (CinemaHallResponse, []error) {
	var err error
	var errors []error

	//validate request
	errors = validateCinemaHallRequiredFields(request)
	if len(errors) > 0 {
		return CinemaHallResponse{}, errors
	}

	//check for duplicate hall names
	hallNames := []string{}
	for _, hall := range request.Halls {
		hallNames = append(hallNames, hall.Name)
	}

	duplicateHallNames := lo.FindDuplicates(hallNames)
	if len(duplicateHallNames) > 0 {
		errors = append(errors, fmt.Errorf("should not have duplicate hall name : %s", strings.Join(duplicateHallNames, ",")))
		return CinemaHallResponse{}, errors
	}

	//check for duplicate seat number
	seatNumbers := []int{}
	for _, hall := range request.Halls {
		for _, seat := range hall.CinemaSeats {
			seatNumbers = append(seatNumbers, seat.SeatNumber)
		}
	}

	duplicateSeatNumbers := lo.FindDuplicates(seatNumbers)
	if len(duplicateSeatNumbers) > 0 {
		duplicateStrs := lo.Map(duplicateSeatNumbers, func(x int, index int) string {
			return strconv.FormatInt(int64(x), 10)
		})
		errors = append(errors, fmt.Errorf("should not have duplicate seat numbers : %s", strings.Join(duplicateStrs, ",")))
		return CinemaHallResponse{}, errors
	}

	cinema := &entities.Cinema{
		Id:           request.CinemaId,
		IsDeprecated: false,
	}

	result := service.DB.First(cinema)
	if result.Error != nil {
		errors = append(errors, result.Error)
		return CinemaHallResponse{}, errors
	}

	//check if the not duplicate names in the DB
	var countResult int64
	service.DB.Find(&entities.CinemaHall{IsDeprecated: false}, duplicateHallNames).Count(&countResult)
	if countResult > 0 {
		errors = append(errors, fmt.Errorf(("cinemaHall Name Already exist")))
		return CinemaHallResponse{}, errors
	}

	//check if the not duplicate names in the DB
	service.DB.Find(&entities.CinemaSeat{IsDeprecated: false}, duplicateSeatNumbers).Count(&countResult)
	if countResult > 0 {
		errors = append(errors, fmt.Errorf(("cinemaseat number already exist")))
		return CinemaHallResponse{}, errors
	}

	err = service.DB.Transaction(func(tx *gorm.DB) error {
		if len(request.Halls) != 0 {
			for _, hall := range request.Halls {
				cinemaHall := entities.CinemaHall{
					Id:           sequentialguid.New().String(),
					Name:         hall.Name,
					TotalSeat:    hall.TotalSeat,
					CinemaId:     request.CinemaId,
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
							Type:         int(seat.Type),
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
		errors = append(errors, result.Error)
		return CinemaHallResponse{}, errors
	}

	resp := CinemaHallResponse{
		CinemaId: request.CinemaId,
	}
	return resp, nil
}

func validateCinemaHallRequiredFields(request CinemaHallRequest) []error {
	validationErrors := []error{}
	if len(request.CinemaId) == 0 || len(request.CinemaId) < 36 {
		validationErrors = append(validationErrors, errors.New("cinemaId is a required field  with 36 characters"))
	}

	if request.CinemaId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("cinemaId should have a valid UUID"))
	}

	for i, hall := range request.Halls {
		if hall.Name == "" {
			validationErrors = append(validationErrors, fmt.Errorf("cinemahall[%d].Name should be supplied", i))
		}

		if hall.TotalSeat <= 0 {
			validationErrors = append(validationErrors, errors.New("cinemahall[%d].TotalSeat should not be lesss than or equal to zero"))
		}

		if hall.TotalSeat > len(hall.CinemaSeats) {
			validationErrors = append(validationErrors, fmt.Errorf("total Seat : %d is greater than the length of seats : %d to be inserted", hall.TotalSeat, len(hall.CinemaSeats)))
		}

		for j, seat := range hall.CinemaSeats {
			if seat.SeatNumber <= 0 {
				validationErrors = append(validationErrors, fmt.Errorf("CinemaSeat[%d].SeatNumber should not have a negative number", j))
			}

			if seat.Type <= 0 {
				validationErrors = append(validationErrors, fmt.Errorf("CinemaSeat[%d].Type should not have a negative number", j))
			}
		}
	}

	return validationErrors
}
