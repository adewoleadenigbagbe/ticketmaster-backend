package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/enums"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
)

type GetAvailableSeatRequest struct {
	Id           string `param:"id"`
	CinemaHallId string `json:"cinemaHallId"`
}

type GetAvailableSeatResponse struct {
	Id                 string             `json:"id"`
	AvailableShowSeats []ShowSeatResponse `json:"availableShowSeats"`
	ReservedShowSeats  []ShowSeatResponse `json:"reservedShowSeats"`
	BookedShowSeats    []ShowSeatResponse `json:"bookedShowSeats"`
	StatusCode         int
}

type ShowSeatResponse struct {
	SeatId       string               `json:"seatId"`
	CinemaSeatId string               `json:"cinemaSeatId"`
	BookingId    string               `json:"bookingId"`
	SeatNumber   int                  `json:"seatNumber"`
	SeatType     enums.SeatType       `json:"seatType"`
	Status       enums.ShowSeatStatus `json:"status"`
	Price        float64              `json:"price"`
}

type SeatDTO struct {
	SeatId       string
	CinemaSeatId string
	BookingId    sql.NullString
	SeatNumber   int
	SeatType     int
	Status       int
	Price        float64
}

func (showService ShowService) GetAvailableShowSeat(request GetAvailableSeatRequest) (GetAvailableSeatResponse, []error) {
	errs := validateAvailableSeat(request)
	if len(errs) != 0 {
		return GetAvailableSeatResponse{StatusCode: http.StatusBadRequest}, errs
	}

	var err error
	seatQuery, err := showService.DB.Table("cinemaseats").
		Where("cinemaseats.CinemaHallId = ?", request.CinemaHallId).
		Where("cinemaseats.IsDeprecated = ?", false).
		Joins("left join showseats on cinemaseats.Id = showseats.CinemaSeatId").
		Where("showseats.ShowId = ?", request.Id).
		Where("showseats.IsDeprecated = ?", false).
		Select("cinemaseats.SeatNumber AS SeatNumber, cinemaseats.Type AS SeatType, showseats.Id AS SeatId, showseats.Status AS Status,showseats.Price AS SeatPrice, showseats.CinemaSeatId AS CinemaSeatId, showseats.BookingId AS BookingId").
		Group("Status").
		Rows()

	if err != nil {
		return GetAvailableSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
	}

	defer seatQuery.Close()

	var seats []SeatDTO
	for seatQuery.Next() {
		seatDTO := SeatDTO{}
		err = seatQuery.Scan(&seatDTO.SeatNumber, &seatDTO.BookingId, &seatDTO.CinemaSeatId, &seatDTO.Price, &seatDTO.SeatId, &seatDTO.SeatType, &seatDTO.Status)
		if err != nil {
			return GetAvailableSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
		}
		seats = append(seats, seatDTO)
	}

	return nil
}

func validateAvailableSeat(request GetAvailableSeatRequest) []error {
	validationErrors := []error{}

	if len(request.Id) == 0 || len(request.Id) < 36 {
		validationErrors = append(validationErrors, errors.New("showId is a required field  with 36 characters"))
	}

	if request.Id == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("showId should have a valid UUID"))
	}

	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		validationErrors = append(validationErrors, errors.New("cinemahallId is a required field  with 36 characters"))
	}

	if request.CinemaHallId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("cinemahallId should have a valid UUID"))
	}

	return validationErrors
}
