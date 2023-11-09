package services

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/common"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

type BookRequest struct {
	UserId        string               `json:"userId"`
	ShowId        string               `json:"showId"`
	CinemaSeatIds []string             `json:"cinemaSeatIds"`
	Status        enums.ShowSeatStatus `json:"status"`
	Price         float64              `json:"price"`
}

type BookResponse struct {
	BookingId  string `json:"bookingId"`
	StatusCode int
}

func (bookService BookService) BookShow(request BookRequest) (BookResponse, []error) {
	var err error
	validationError := validateShowBook(request)
	if len(validationError) != 0 {
		return BookResponse{StatusCode: http.StatusBadRequest}, validationError
	}

	//check if the seat already existed
	var count int64
	bookService.DB.Table("showseats").Where("CinemaSeatId IN ?", request.CinemaSeatIds).Count(&count)
	if count > 0 {
		return BookResponse{StatusCode: http.StatusBadRequest}, []error{fmt.Errorf("show seat already reserved or booked")}
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
			showseat := entities.ShowSeat{
				Id:           sequentialguid.New().String(),
				Status:       request.Status,
				Price:        request.Price,
				CinemaSeatId: cinemaSeatId,
				ShowId:       request.ShowId,
				BookingId:    sql.NullString{String: booking.Id, Valid: true},
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
		return BookResponse{StatusCode: http.StatusBadRequest}, []error{err}
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
		log.Fatal("encode error:", err)
	}

	if err := bookService.Nc.Publish(common.EventSubject, buf.Bytes()); err != nil {
		log.Fatal(err)
	}

	bookService.Nc.Flush()
	return BookResponse{BookingId: booking.Id, StatusCode: http.StatusOK}, nil
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

	if int(request.Price) <= 0 {
		vErrors = append(vErrors, errors.New("price should not be less than or equal to zero"))
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
