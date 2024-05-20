package services

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/shared/common"
	"github.com/Wolechacho/ticketmaster-backend/shared/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/shared/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/shared/helpers"
	"github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/shared/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type BookRequest struct {
	UserId        string               `json:"userId"`
	ShowId        string               `json:"showId"`
	CinemaSeatIds []string             `json:"cinemaSeatIds"`
	Status        enums.ShowSeatStatus `json:"status"`
}

type BookResponse struct {
	BookingId string `json:"bookingId"`
}

type RateModel struct {
	BaseFee    float64
	Discount   utilities.Nullable[float64]
	IsSpecials utilities.Nullable[bool]
}

type SeatTypeModel struct {
	Id   string
	Type enums.SeatType
}

func (bookService BookService) BookShow(request BookRequest) (BookResponse, models.ErrorResponse) {
	var err error
	validationError := validateShowBook(request)
	if len(validationError) != 0 {
		return BookResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: validationError}
	}

	//check if the seat already existed
	var count int64
	bookService.DB.Table("showseats").
		Where("CinemaSeatId IN ?", request.CinemaSeatIds).
		Where("Status != ?", enums.ExpiredSeat).
		Count(&count)

	if count > 0 {
		return BookResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{fmt.Errorf("show seat already reserved or booked")}}
	}
	//get seat type
	seatQuery, err := bookService.DB.Table("cinemaseats").
		Where("Id IN ?", request.CinemaSeatIds).
		Where("IsDeprecated = ?", false).
		Select("Id", "Type").
		Rows()

	if err != nil {
		return BookResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
	}

	defer seatQuery.Close()
	var seatTypes []SeatTypeModel
	for seatQuery.Next() {
		seatType := SeatTypeModel{}
		err = seatQuery.Scan(&seatType.Id, &seatType.Type)
		if err != nil {
			return BookResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
		}
		seatTypes = append(seatTypes, seatType)
	}

	groupedSeatTypeById := lo.KeyBy(seatTypes, func(item SeatTypeModel) string {
		return item.Id
	})

	// get the cinema rate
	var rateModel RateModel
	bookService.DB.Table("shows").
		Where("shows.Id = ?", request.ShowId).
		Where("shows.IsDeprecated = ?", false).
		Joins("join cinemahalls on shows.CinemaHallId = cinemahalls.Id").
		Where("cinemahalls.IsDeprecated = ?", false).
		Joins("join cinemas on cinemahalls.CinemaId = cinemas.Id").
		Where("cinemas.IsDeprecated = ?", false).
		Joins("join cinemarates on cinemas.Id = cinemarates.CinemaId").
		Where("cinemarates.IsActive = ?", true).
		Select("cinemarates.BaseFee", "cinemarates.Discount", "cinemarates.IsSpecials").
		Scan(&rateModel)

	var rate float64
	if rateModel.IsSpecials.Valid && rateModel.Discount.Valid {
		rate = (rateModel.BaseFee - (rateModel.BaseFee * rateModel.Discount.Val))
	}

	today := time.Now()
	booking := entities.Booking{
		Id:            sequentialguid.New().String(),
		BookDateTime:  today,
		NumberOfSeats: len(request.CinemaSeatIds),
		Status:        enums.PendingBook,
		UserId:        request.UserId,
		ShowId:        request.ShowId,
		IsDeprecated:  false,
	}
	err = bookService.DB.Transaction(func(tx *gorm.DB) error {
		if err = bookService.DB.Create(&booking).Error; err != nil {
			// return any error will rollback
			return err
		}

		showSeats := []*entities.ShowSeat{}
		for _, cinemaSeatId := range request.CinemaSeatIds {
			seatType := groupedSeatTypeById[cinemaSeatId]
			var price float64
			if seatType.Type == enums.Standard {
				price = rate * 2
			} else if seatType.Type == enums.Gold {
				price = rate * 3
			} else {
				price = rate
			}

			showseat := entities.ShowSeat{
				Id:           sequentialguid.New().String(),
				Status:       request.Status,
				Price:        price,
				CinemaSeatId: cinemaSeatId,
				ShowId:       request.ShowId,
				BookingId:    utilities.NewNullable[string](booking.Id, true),
				IsDeprecated: false,
			}
			showSeats = append(showSeats, &showseat)
		}
		if err = bookService.DB.Create(showSeats).Error; err != nil {
			// return any error will rollback
			return err
		}
		return nil
	})

	if err != nil {
		return BookResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}

	//push the message to a queue and return immediately to user response
	bk := common.BookingMessage{
		UserId:          request.UserId,
		CinemaSeatIds:   request.CinemaSeatIds,
		ShowId:          request.ShowId,
		Status:          request.Status,
		BookingDateTime: today,
		ExpiryDateTime:  today.Add(5 * time.Minute),
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(bk)
	if err != nil {
		//TODO: log to file, cmd , or elastic search
		return BookResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
	}

	if err := bookService.Nc.Publish(common.EventSubject, buf.Bytes()); err != nil {
		//TODO: log to file, cmd , or elastic search
		return BookResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusInternalServerError}
	}

	bookService.Nc.Flush()
	return BookResponse{BookingId: booking.Id}, models.ErrorResponse{}
}

func validateShowBook(request BookRequest) []error {
	vErrors := []error{}
	if len(request.ShowId) == 0 || len(request.ShowId) < 36 {
		vErrors = append(vErrors, errors.New("showId is a required field  with 36 characters"))
	}

	if request.ShowId == utilities.DEFAULT_UUID {
		vErrors = append(vErrors, errors.New("showId should have a valid UUID"))
	}

	if request.Status != enums.Reserved && request.Status != enums.PendingAssignment {
		vErrors = append(vErrors, errors.New("seat can only be reserved or booked"))
	}

	for i, cinemaSeatId := range request.CinemaSeatIds {
		if len(cinemaSeatId) == 0 || len(cinemaSeatId) < 36 {
			vErrors = append(vErrors, fmt.Errorf("CinemaSeatIds[%d] is a required field  with 36 characters", i))
		}

		if cinemaSeatId == utilities.DEFAULT_UUID {
			vErrors = append(vErrors, fmt.Errorf("CinemaSeatIds[%d] should have a valid UUID", i))
		}
	}

	return vErrors
}
