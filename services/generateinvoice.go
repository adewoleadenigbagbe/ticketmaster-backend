package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/enums"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
)

type GeneratePdfRequest struct {
	BookingId string `json:"bookingId"`
	UserId    string `json:"userId"`
}
type GeneratePdfResponse struct {
	PdfBytes []byte
}

type PdfModel struct {
	CinemaName    string
	MovieTitle    string
	ShowStartTime string
	ShowEndTime   string
	TicketNumber  string
	SeatInfos     []SeatInfo
	SubTotal      float64
	Tax           float64
	Total         float64
	HallName      string
}

type BookingModel struct {
	CinemaName    string
	MovieTitle    string
	ShowStartTime int64
	ShowEndTime   int64
	TicketNumber  string
	SeatNumber    int
	SeatType      int
	Price         float64
	HallName      string
}

type SeatInfo struct {
	SeatNumber int
	SeatType   int
	Price      float64
}

func (bookService BookService) GenerateInvoicePDF(request GeneratePdfRequest) (GeneratePdfResponse, models.ErrorResponse) {
	var err error
	requiredFieldErrors := validateInvoice(request)
	if len(requiredFieldErrors) > 0 {
		return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: requiredFieldErrors}
	}

	query, err := bookService.DB.Table("bookings").
		Where("bookings.Id = ?", request.BookingId).
		Where("bookings.UserId = ?", request.UserId).
		Where("bookings.Status = ?", enums.Booked).
		Where("bookings.IsDeprecated = ? ", false).
		Joins("join showseats on bookings.Id = showseats.BookingId").
		Where("showseats.Status = ?", enums.Assigned).
		Where("showseats.IsDeprecated = ? ", false).
		Joins("join cinemaseats on showseats.CinemaSeatId = cinemaseats.Id").
		Where("cinemaseats.IsDeprecated = ? ", false).
		Joins("join shows on bookings.ShowId = shows.Id").
		Where("shows.IsDeprecated = ? ", false).
		Joins("join movies on shows.MovieId = movies.Id").
		Where("movies.IsDeprecated = ? ", false).
		Joins("join cinemahalls on cinemaseats.CinemaHallId = cinemahalls.Id").
		Where("cinemahalls.IsDeprecated = ? ", false).
		Joins("join cinemas on cinemahalls.CinemaId = cinemas.Id").
		Where("cinemas.IsDeprecated = ? ", false).
		Select("bookings.Id AS TicketNumber",
			"cinemas.Name AS CinemaName",
			"movies.Title AS MovieTitle",
			"shows.StartTime AS ShowStartTime",
			"shows.EndTime AS ShowEndTime",
			"cinemaseats.SeatNumber AS SeatNumber",
			"cinemaseats.Type AS SeatType",
			"showseats.Price AS Price",
			"cinemahalls.Name AS HallName").
		Rows()

	if err != nil {
		return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
	}

	defer query.Close()

	bookings := []BookingModel{}
	for query.Next() {
		var booking BookingModel
		err = query.Scan(&booking.TicketNumber, &booking.CinemaName, &booking.MovieTitle, &booking.ShowStartTime, &booking.ShowEndTime, &booking.SeatNumber, &booking.SeatType, &booking.Price, &booking.HallName)
		if err != nil {
			return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
		}

		bookings = append(bookings, booking)
	}

	if len(bookings) == 0 {
		return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusNotFound, Errors: []error{errors.New("booking not found")}}
	}

	var ticketNumber string
	pdfModel := PdfModel{}

	for _, booking := range bookings {
		if ticketNumber == "" {
			pdfModel.TicketNumber = strings.Replace(booking.TicketNumber, "-", "", -1)
			pdfModel.CinemaName = booking.CinemaName
			pdfModel.HallName = booking.HallName
			pdfModel.MovieTitle = booking.MovieTitle
			pdfModel.ShowStartTime = time.Unix(booking.ShowStartTime, 0).Format(time.DateTime)
			pdfModel.ShowEndTime = time.Unix(booking.ShowEndTime, 0).Format(time.DateTime)
			pdfModel.Tax = 1.02
		}

		var seat SeatInfo
		seat.SeatNumber = booking.SeatNumber
		seat.SeatType = booking.SeatType
		seat.Price = booking.Price
		pdfModel.SeatInfos = append(pdfModel.SeatInfos, seat)
		pdfModel.SubTotal += booking.Price

		ticketNumber = booking.TicketNumber
	}

	pdfModel.Total = pdfModel.SubTotal + pdfModel.Tax

	pdfBytes, err := bookService.PDFService.GeneratePDF(pdfModel)
	if err != nil {
		return GeneratePdfResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
	}

	return GeneratePdfResponse{PdfBytes: pdfBytes}, models.ErrorResponse{}
}

func validateInvoice(request GeneratePdfRequest) []error {
	validationErrors := []error{}
	if len(request.BookingId) == 0 || len(request.BookingId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "bookingId"))
	}

	if request.BookingId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "bookingId"))
	}

	if len(request.UserId) == 0 || len(request.UserId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "userId"))
	}

	if request.UserId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "userId"))
	}

	return validationErrors
}
