package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateCinemaSeatRequest struct {
	Id           string            `param:"Id"`
	CinemaHallId string            `param:"cinemaHallId"`
	Seats        []CinemaSeatModel `json:"cinemaSeats"`
}

type CreateCinemaSeatResponse struct {
	StatusCode int
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

func (cinemaService CinemaService) AddCinemaSeat(request CreateCinemaSeatRequest) (CreateCinemaSeatResponse, []error) {
	validationErrors := validateCinemSeatRequiredFields(request)
	if len(validationErrors) > 0 {
		return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, validationErrors
	}

	seatNumbers := []int{}
	if len(request.Seats) > 1 {
		for _, seat := range request.Seats {
			seatNumbers = append(seatNumbers, seat.SeatNumber)
		}

		duplicateSeatNumbers := lo.FindDuplicates(seatNumbers)
		if len(duplicateSeatNumbers) > 0 {
			return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, []error{errors.New("seat number in the request contains duplicates")}
		}
	}

	var err error
	cinemaQuery, err := cinemaService.DB.Table("cinemas").
		Where("cinemas.Id = ?", request.Id).
		Where("cinemas.IsDeprecated = ?", false).
		Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
		Where("cinemaHalls.Id = ?", request.CinemaHallId).
		Where("cinemaHalls.IsDeprecated = ?", false).
		Select("cinemas.Id AS CinemaId, cinemaHalls.Id AS CinemaHallId, cinemaHalls.TotalSeat AS TotalSeat").
		Rows()

	if err != nil {
		return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
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
			return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
		}
		i++
	}

	numOfSeats := len(request.Seats)
	if existingHalls.TotalSeat < numOfSeats {
		return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, []error{fmt.Errorf(("total number of cinema seats in the system is less that the new seats to add"))}
	}

	cinemaHallQuery, err := cinemaService.DB.Table("cinemaHalls").
		Where("cinemaHalls.Id = ?", existingHalls.CinemaHallId).
		Where("cinemaHalls.IsDeprecated = ?", false).
		Joins("join cinemaSeats on cinemaHalls.Id = cinemaSeats.CinemaHallId").
		Where("cinemaSeats.IsDeprecated = ?", false).
		Select("cinemaHalls.Id AS CinemaHallId, cinemaSeats.SeatNumber AS SeatNumber").
		Rows()

	if err != nil {
		return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
	}

	defer cinemaHallQuery.Close()

	var existingSeats []CinemaSeatDTO
	for cinemaHallQuery.Next() {
		var cinemaSeatDTO CinemaSeatDTO
		err = cinemaHallQuery.Scan(&cinemaSeatDTO.CinemaHallId, &cinemaSeatDTO.SeatNumber)
		if err != nil {
			return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
		}
		existingSeats = append(existingSeats, cinemaSeatDTO)
	}

	if len(existingSeats)+len(request.Seats) > existingHalls.TotalSeat {
		return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, []error{fmt.Errorf(("total number of cinema seats in the system is less that the new seats to add"))}
	}

	//check for duplicates
	for _, exSeat := range existingSeats {
		for _, seat := range request.Seats {
			if exSeat.SeatNumber == seat.SeatNumber {
				return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, []error{errors.New("seat number already exist in the db")}
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
					Type:         int(seat.Type),
					IsDeprecated: false,
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
		return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
	}

	return CreateCinemaSeatResponse{StatusCode: http.StatusOK}, nil
}

func validateCinemSeatRequiredFields(request CreateCinemaSeatRequest) []error {
	validationErrors := []error{}
	if len(request.Id) == 0 || len(request.Id) < 36 {
		validationErrors = append(validationErrors, errors.New("cinemaId is a required field  with 36 characters"))
	}

	if request.Id == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("cinemaId should have a valid UUID"))
	}

	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		validationErrors = append(validationErrors, errors.New("cinemahallId is a required field  with 36 characters"))
	}

	if request.CinemaHallId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("cinemahallId should have a valid UUID"))
	}

	for i, seat := range request.Seats {
		if seat.SeatNumber <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("CinemaSeat[%d].SeatNumber should not have a negative number", i))
		}

		if seat.Type <= 0 {
			validationErrors = append(validationErrors, fmt.Errorf("CinemaSeat[%d].Type should not have a negative number", i))
		}

	}

	return validationErrors
}
