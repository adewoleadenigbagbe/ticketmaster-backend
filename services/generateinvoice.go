package services

import (
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/models"
)

type GeneratePdfRequest struct {
	BookingId string `json:"bookingId"`
	UserId    string `json:"userId"`
}
type GeneratePdfResponse struct {
	PdfBytes []byte
}

type BookingModel struct {
	CinemaName   string
	MovieTitle   string
	ShowTime     string
	TicketNumber int
	SeatInfos    []SeatInfo
	SubTotal     float64
	Tax          float32
	Total        float64
}

type SeatInfo struct {
	SeatNumber int
	SeatType   int
	Price      float64
}

func (bookingService BookingService) GenerateInvoicePDF(request GeneratePdfRequest) (GeneratePdfResponse, models.ErrorResponse) {
	data := BookingModel{
		CinemaName:   "Davis Cinema",
		MovieTitle:   "Trannsformer",
		ShowTime:     time.Now().String(),
		TicketNumber: 1,
		SeatInfos: []SeatInfo{
			{
				SeatNumber: 2,
				SeatType:   3,
				Price:      45.09,
			},
			{
				SeatNumber: 3,
				SeatType:   4,
				Price:      55.09,
			},
			{
				SeatNumber: 5,
				SeatType:   6,
				Price:      65.09,
			},
		},

		SubTotal: 30.009,
		Tax:      1.08,
		Total:    60.89,
	}
	pdfBytes, err := bookingService.PDFService.GeneratePDF(data)
	if err != nil {
		return GeneratePdfResponse{}, models.ErrorResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
	}

	return GeneratePdfResponse{PdfBytes: pdfBytes}, models.ErrorResponse{}
}
