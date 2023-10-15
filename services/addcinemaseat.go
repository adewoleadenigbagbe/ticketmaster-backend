package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreateCinemaSeatRequest struct {
	CinemaId     string            `param:"cinemaId"`
	CinemaHallId string            `param:"cinemaHallId"`
	Seats        []CinemaSeatModel `json:"cinemaSeats"`
}

type CreateCinemaSeatResponse struct {
	StatusCode int
}

type CinemaDTO struct {
	CinemaId     string
	CinemaHallId string
	TotalSeat    int
}

func (cinemaService CinemaService) AddCinemaSeat(request CreateCinemaSeatRequest) (CreateCinemaSeatResponse, error) {

	seatNumbers := []int{}
	if len(request.Seats) > 1 {
		for _, seat := range request.Seats {
			seatNumbers = append(seatNumbers, seat.SeatNumber)
		}

		duplicateSeatNumbers := lo.FindDuplicates(seatNumbers)
		if len(duplicateSeatNumbers) > 0 {
			return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, errors.New("seat number in the request contains duplicates")
		}
	}

	//errors.Is(result.Error, gorm.ErrRecordNotFound)
	var err error
	cinemaQuery, err := cinemaService.DB.Table("cinemas").
		Where("cinemas.Id = ?", request.CinemaId).
		Where("cinemas.IsDeprecated = ?", false).
		Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
		Where("cinemaHalls.Id = ?", request.CinemaHallId).
		Where("cinemaHalls.IsDeprecated = ?", false).
		Select("cinemas.Id AS CinemaId, cinemaHalls.Id AS cinemaHallId, cinemaHalls.TotalSeat").
		Rows()

	if err != nil {
		return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, err
	}

	defer cinemaQuery.Close()

	var cinemaDTO CinemaDTO
	i := 0
	for cinemaQuery.Next() {
		if i > 1 {
			break
		}
		err = cinemaQuery.Scan(&cinemaDTO.CinemaId, &cinemaDTO.CinemaHallId, &cinemaDTO.TotalSeat)
		if err != nil {
			return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, err
		}
		i++
	}

	//check for duplicates seat number in the DB
	var countResult int64
	cinemaService.DB.Table("cinemaHalls").
		Where("cinemaHalls.Id = ?", cinemaDTO.CinemaHallId).
		Where("cinemaHalls.IsDeprecated = ?", false).
		Joins("join cinemaSeats on cinemaHalls.Id = cinemaSeats.CinemaHallId").
		Where("cinemaHalls.Id = ?", request.CinemaHallId).
		Where("cinemaSeats.IsDeprecated = ?", false).
		Count(&countResult)

	if countResult > 0 {
		return CreateCinemaSeatResponse{StatusCode: http.StatusBadRequest}, fmt.Errorf(("seat number already exist in the DB"))
	}

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
		return CreateCinemaSeatResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return CreateCinemaSeatResponse{StatusCode: http.StatusOK}, nil
}
