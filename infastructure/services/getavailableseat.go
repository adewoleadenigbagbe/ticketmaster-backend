package services

import (
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/shared/enums"
	"github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities"
)

type GetAvailableSeatRequest struct {
	Id           string `param:"id"`
	CinemaHallId string `json:"cinemaHallId"`
	SortBy       string `json:"sortBy" default:"SeatNumber"`
	Order        string `json:"order" default:"asc"`
}

type GetAvailableSeatResponse struct {
	Id                 string             `json:"id"`
	AvailableShowSeats []ShowSeatResponse `json:"availableShowSeats"`
	ReservedShowSeats  []ShowSeatResponse `json:"reservedShowSeats"`
	AssignedShowSeats  []ShowSeatResponse `json:"assignedShowSeats"`
	StatusCode         int
}

type ShowSeatResponse struct {
	SeatId       utilities.Nullable[string]  `json:"seatId"`
	CinemaSeatId utilities.Nullable[string]  `json:"cinemaSeatId"`
	BookingId    utilities.Nullable[string]  `json:"bookingId"`
	SeatNumber   int                         `json:"seatNumber"`
	SeatType     enums.SeatType              `json:"seatType"`
	Status       enums.ShowSeatStatus        `json:"status"`
	Price        utilities.Nullable[float64] `json:"price"`
}

type SeatDTO struct {
	SeatId       utilities.Nullable[string]
	CinemaSeatId utilities.Nullable[string]
	BookingId    utilities.Nullable[string]
	SeatNumber   int
	SeatType     int
	Status       utilities.Nullable[int]
	Price        utilities.Nullable[float64]
}

func (showService ShowService) GetAvailableShowSeat(request GetAvailableSeatRequest) (GetAvailableSeatResponse, []error) {
	var err error
	err = utilities.SetDefaults[GetAvailableSeatRequest](&request)
	if err != nil {
		return GetAvailableSeatResponse{StatusCode: http.StatusBadRequest}, []error{err}
	}

	errs := validateAvailableSeat(request)
	if len(errs) != 0 {
		return GetAvailableSeatResponse{StatusCode: http.StatusBadRequest}, errs
	}

	sortOrder := fmt.Sprint(request.SortBy, " ", request.Order)

	seatQuery, err := showService.DB.Table("cinemaseats").
		Where("cinemaseats.CinemaHallId = ?", request.CinemaHallId).
		Where("cinemaseats.IsDeprecated = ?", false).
		Joins("left join showseats on cinemaseats.Id = showseats.CinemaSeatId").
		Where("(showseats.ShowId = ? AND showseats.Status != ?) OR showseats.ShowId IS NULL", request.Id, enums.ExpiredSeat).
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
		AssignedShowSeats:  []ShowSeatResponse{},
	}

	for _, seatDTO := range seatsDTO {
		seat := ShowSeatResponse{
			SeatId:       seatDTO.SeatId,
			CinemaSeatId: seatDTO.CinemaSeatId,
			SeatNumber:   seatDTO.SeatNumber,
			SeatType:     enums.SeatType(seatDTO.SeatType),
			Price:        seatDTO.Price,
			BookingId:    seatDTO.BookingId,
		}

		if !seatDTO.SeatId.Valid && seat.SeatId.Val == "" {
			seat.Status = enums.Available
			resp.AvailableShowSeats = append(resp.AvailableShowSeats, seat)
		} else {
			if int(seatDTO.Status.Val) == int(enums.Assigned) {
				seat.Status = enums.Assigned
				resp.AssignedShowSeats = append(resp.AssignedShowSeats, seat)
			} else if int(seatDTO.Status.Val) == int(enums.Reserved) || int(seatDTO.Status.Val) == int(enums.PendingAssignment) {
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

		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredField, "showId"))
	}

	if request.Id == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "showId"))
	}

	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredField, "cinemahallId"))
	}

	if request.CinemaHallId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "cinemahallId"))
	}

	return validationErrors
}
