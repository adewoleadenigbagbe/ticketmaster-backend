package services

import (
	"bytes"
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
)

type BookRequest struct {
	UserId        string              `json:"userId"`
	ShowId        string              `json:"showId"`
	CinemaSeatIds []string            `json:"cinemaSeatIds"`
	BookStatus    enums.BookingStatus `json:"bookStatus"`
	Price         float64             `json:"price"`
}

type BookResponse struct {
	StatusCode int
}

func (bookService BookService) BookShow(request BookRequest) (BookResponse, []error) {

	validationError := validateShowBook(request)
	if len(validationError) != 0 {
		return BookResponse{StatusCode: http.StatusBadRequest}, validationError
	}

	showSeats := []*entities.ShowSeat{}
	for _, cinemaSeatId := range request.CinemaSeatIds {
		showseat := entities.ShowSeat{
			Id:           sequentialguid.New().String(),
			Status:       int(request.BookStatus),
			Price:        request.Price,
			CinemaSeatId: cinemaSeatId,
			ShowId:       request.ShowId,
			IsDeprecated: false,
		}
		showSeats = append(showSeats, &showseat)
	}

	result := bookService.DB.Create(showSeats)

	if result.Error != nil && result.RowsAffected < 0 {
		return BookResponse{StatusCode: http.StatusBadRequest}, []error{result.Error}
	}

	//push the message to a queue and return immediately to user response
	today := time.Now()
	bk := common.BookingMessage{
		CinemaSeatIds:   request.CinemaSeatIds,
		ShowId:          request.ShowId,
		Status:          request.BookStatus,
		BookingDateTime: today,
		ExpiryDateTime:  today.Add(5 * time.Minute),
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(bk)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	if err := bookService.Nc.Publish(common.EventSubject, buf.Bytes()); err != nil {
		log.Fatal(err)
	}
	bookService.Nc.Flush()

	return BookResponse{StatusCode: http.StatusOK}, nil
}

func validateShowBook(request BookRequest) []error {
	vErrors := []error{}
	if len(request.ShowId) == 0 || len(request.ShowId) < 36 {
		vErrors = append(vErrors, errors.New("showId is a required field  with 36 characters"))
	}

	if request.ShowId == utilities.DEFAULT_UUID {
		vErrors = append(vErrors, errors.New("showId should have a valid UUID"))
	}

	if int(request.BookStatus) <= 0 {
		vErrors = append(vErrors, errors.New("bookstatus should not be less than or equal to zero"))
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
