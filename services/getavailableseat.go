package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/enums"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
)

type GetAvailableSeatRequest struct {
	Id           string `param:"id"`
	CinemaHallId string `json:"cinemaHallId"`
	SortBy       string `json:"sortBy"`
	Order        string `json:"order"`
}

type GetAvailableSeatResponse struct {
	Id                 string             `json:"id"`
	AvailableShowSeats []ShowSeatResponse `json:"availableShowSeats"`
	ReservedShowSeats  []ShowSeatResponse `json:"reservedShowSeats"`
	BookedShowSeats    []ShowSeatResponse `json:"bookedShowSeats"`
	StatusCode         int
}

type ShowSeatResponse struct {
	SeatId       utilities.JsonNullString `json:"seatId"`
	CinemaSeatId utilities.JsonNullString `json:"cinemaSeatId"`
	BookingId    utilities.JsonNullString `json:"bookingId"`
	SeatNumber   int                      `json:"seatNumber"`
	SeatType     enums.SeatType           `json:"seatType"`
	Status       enums.ShowSeatStatus     `json:"status"`
	Price        utilities.JsonNullFloat  `json:"price"`
}

type SeatDTO struct {
	SeatId       sql.NullString
	CinemaSeatId sql.NullString
	BookingId    sql.NullString
	SeatNumber   int
	SeatType     int
	Status       sql.NullInt32
	Price        sql.NullFloat64
}

func (showService ShowService) GetAvailableShowSeat(request GetAvailableSeatRequest) (GetAvailableSeatResponse, []error) {
	errs := validateAvailableSeat(request)
	if len(errs) != 0 {
		return GetAvailableSeatResponse{StatusCode: http.StatusBadRequest}, errs
	}

	if request.SortBy == "" {
		request.SortBy = "SeatNumber"
	}

	if request.Order == "" {
		request.Order = "asc"
	}

	sortOrder := fmt.Sprint(request.SortBy, " ", request.Order)

	var err error
	seatQuery, err := showService.DB.Table("cinemaseats").
		Where("cinemaseats.CinemaHallId = ?", request.CinemaHallId).
		Where("cinemaseats.IsDeprecated = ?", false).
		Joins("left join showseats on cinemaseats.Id = showseats.CinemaSeatId").
		Where("showseats.ShowId = ? OR showseats.ShowId IS NULL", request.Id).
		Where("showseats.IsDeprecated = ? OR showseats.IsDeprecated IS NULL", false).
		Select("cinemaseats.SeatNumber AS SeatNumber, cinemaseats.Type AS SeatType, showseats.Id AS SeatId, showseats.Status AS Status,showseats.Price AS SeatPrice, showseats.CinemaSeatId AS CinemaSeatId, showseats.BookingId AS BookingId").
		Order(sortOrder).
		Rows()

	if err != nil {
		return GetAvailableSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
	}

	defer seatQuery.Close()
	var seatsDTO []SeatDTO
	for seatQuery.Next() {
		seatDTO := SeatDTO{}
		err = seatQuery.Scan(&seatDTO.SeatNumber, &seatDTO.SeatType, &seatDTO.SeatId, &seatDTO.Status, &seatDTO.Price, &seatDTO.CinemaSeatId, &seatDTO.BookingId)
		if err != nil {
			return GetAvailableSeatResponse{StatusCode: http.StatusInternalServerError}, []error{err}
		}
		seatsDTO = append(seatsDTO, seatDTO)
	}

	resp := GetAvailableSeatResponse{StatusCode: http.StatusOK,
		Id:                 request.Id,
		AvailableShowSeats: []ShowSeatResponse{},
		ReservedShowSeats:  []ShowSeatResponse{},
		BookedShowSeats:    []ShowSeatResponse{},
	}

	for _, seatDTO := range seatsDTO {
		seat := ShowSeatResponse{
			SeatId:       utilities.JsonNullString{NullString: seatDTO.SeatId},
			CinemaSeatId: utilities.JsonNullString{NullString: seatDTO.CinemaSeatId},
			SeatNumber:   seatDTO.SeatNumber,
			SeatType:     enums.SeatType(seatDTO.SeatType),
			Price:        utilities.JsonNullFloat{NullFloat64: seatDTO.Price},
			BookingId:    utilities.JsonNullString{NullString: seatDTO.BookingId},
		}

		if !seatDTO.SeatId.Valid && seat.SeatId.String == "" {
			seat.Status = enums.Available
			resp.AvailableShowSeats = append(resp.AvailableShowSeats, seat)
		} else {
			if int(seatDTO.Status.Int32) == int(enums.Booked) {
				seat.Status = enums.Assigned
				resp.BookedShowSeats = append(resp.BookedShowSeats, seat)
			} else if int(seatDTO.Status.Int32) == int(enums.Reserved) {
				seat.Status = enums.Reserved
				resp.ReservedShowSeats = append(resp.ReservedShowSeats, seat)
			} else {
				seat.Status = enums.Available
				resp.AvailableShowSeats = append(resp.AvailableShowSeats, seat)
			}
		}

	}
	return resp, nil
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