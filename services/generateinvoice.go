package services

import (
	"errors"
	"net/http"
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
	Tax           float32
	Total         float64
}

type BookingModel struct {
	CinemaName    string
	MovieTitle    string
	ShowStartTime int
	ShowEndTime   int
	TicketNumber  string
	SeatNumber    int
	SeatType      int
	Price         float64
}

type SeatInfo struct {
	SeatNumber int
	SeatType   int
	Price      float64
}

func (bookingService BookingService) GenerateInvoicePDF(request GeneratePdfRequest) (GeneratePdfResponse, models.ErrorResponse) {
	var err error
	requiredFieldErrors := validateInvoice(request)
	if len(requiredFieldErrors) > 0 {
		return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: requiredFieldErrors}
	}

	query, err := bookingService.DB.Table("bookings").
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
			"cinemaseats.SeatType AS SeatType",
			"showseats.Price AS Price").
		Rows()

	if err != nil {
		return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
	}

	defer query.Close()

	bookings := []BookingModel{}
	for query.Next() {
		var booking BookingModel
		err = query.Scan(&booking.TicketNumber, &booking.CinemaName, &booking.MovieTitle, &booking.ShowStartTime, &booking.ShowEndTime, &booking.SeatNumber, &booking.SeatType, &booking.Price)
		if err != nil {
			return GeneratePdfResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
		}

		bookings = append(bookings, booking)
	}

	var ticketNumber string
	pdfModel := PdfModel{}

	for _, booking := range bookings {
		if ticketNumber == "" {
			pdfModel.TicketNumber = booking.TicketNumber
			pdfModel.CinemaName = booking.CinemaName
			pdfModel.MovieTitle = booking.MovieTitle
			pdfModel.ShowStartTime = time.Unix(int64(booking.ShowStartTime), 0).String()
			pdfModel.ShowEndTime = time.Unix(int64(booking.ShowEndTime), 0).String()
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

	pdfModel.Total = pdfModel.SubTotal + float64(pdfModel.Tax)

	pdfBytes, err := bookingService.PDFService.GeneratePDF(pdfModel)
	if err != nil {
		return GeneratePdfResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
	}

	return GeneratePdfResponse{PdfBytes: pdfBytes}, models.ErrorResponse{}
}

func validateInvoice(request GeneratePdfRequest) []error {
	validationErrors := []error{}
	if len(request.BookingId) == 0 || len(request.BookingId) < 36 {
		validationErrors = append(validationErrors, errors.New("bookingId is a required field  with 36 characters"))
	}

	if request.BookingId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("bookingId should have a valid UUID"))
	}

	if len(request.UserId) == 0 || len(request.UserId) < 36 {
		validationErrors = append(validationErrors, errors.New("userId is a required field  with 36 characters"))
	}

	if request.UserId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, errors.New("userId should have a valid UUID"))
	}

	return validationErrors
}
