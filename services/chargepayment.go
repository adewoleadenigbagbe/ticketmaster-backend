package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"gorm.io/gorm"
)

type CreatePaymentRequest struct {
	UserId    string `json:"userId"`
	BookingId string `json:"bookingId"`
}

type CreatePaymentResponse struct {
	PaymentId string `json:"paymentId"`
}

type BookModel struct {
	BookingId     string
	BookingStatus enums.BookingStatus
	SeatStatus    enums.ShowSeatStatus
	Price         float64
}

func (bookService BookService) ChargeBooking(request CreatePaymentRequest) (CreatePaymentResponse, models.ErrorResponse) {
	var err error
	fieldErrors := validatePayment(request)
	if len(fieldErrors) != 0 {
		return CreatePaymentResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: fieldErrors}
	}

	query, err := bookService.DB.Table("bookings").
		Where("bookings.Id = ?", request.BookingId).
		Where("bookings.UserId = ?", request.UserId).
		Where("bookings.Status = ?", enums.PendingBook).
		Where("bookings.IsDeprecated = ?", false).
		Joins("join showseats on bookings.Id = showseats.BookingId").
		Where("showseats.Status = ?", enums.PendingAssignment).
		Where("showseats.IsDeprecated = ?", false).
		Select("bookings.Id AS BookingId",
			"bookings.Status AS BookingStatus",
			"showseats.Status AS SeatStatus",
			"showseats.Price AS Price").
		Rows()

	if err != nil {
		return CreatePaymentResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
	}
	defer query.Close()

	bookings := []BookModel{}
	for query.Next() {
		var booking BookModel
		err = query.Scan(&booking.BookingId, &booking.BookingStatus, &booking.Price, &booking.SeatStatus)
		if err != nil {
			return CreatePaymentResponse{}, models.ErrorResponse{StatusCode: http.StatusInternalServerError, Errors: []error{err}}
		}
		bookings = append(bookings, booking)
	}

	if len(bookings) == 0 {
		return CreatePaymentResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{errors.New("seat booked not found")}}
	}

	paymentId := sequentialguid.New().String()

	bookService.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Table("bookings").
			Where("Id = ? ", request.BookingId).
			Where("BookingId = ?", request.BookingId).
			Where("IsDeprecated = ?", false).
			Update("Status", enums.Booked).Error; err != nil {
			return err
		}

		seatStatuses := lo.Map(bookings, func(item BookModel, index int) enums.ShowSeatStatus {
			return item.SeatStatus
		})

		if err = tx.Table("showseats").
			Where("BookingId = ?", request.BookingId).
			Where("IsDeprecated = ?", false).
			Where("Status IN ?", seatStatuses).
			Update("Status", enums.Assigned).Error; err != nil {
			return err
		}

		amount := lo.SumBy(bookings, func(item BookModel) float64 {
			return item.Price
		})

		// Set Stripe API key
		stripe.Key = "YOUR_API_KEY"

		// Create a new charge
		params := &stripe.ChargeParams{
			Amount:      stripe.Int64(int64(amount)),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Description: stripe.String("Test Charge"),
			Source:      &stripe.SourceParams{Token: stripe.String("tok_visa")}, // use a test card token provided by Stripe
		}

		//ignore the error as this just test
		ch, _ := charge.New(params)
		stripeInfoBytes, _ := json.Marshal(*ch)
		fmt.Println(string(stripeInfoBytes))

		payment := entities.Payment{
			Id:                       paymentId,
			Amount:                   amount,
			PaymentDate:              time.Now(),
			PaymentMethod:            enums.Stripe,
			RemoteTransactionId:      sql.NullString{String: ch.ID, Valid: true},
			BookingId:                request.BookingId,
			ProviderExtraInformation: sql.NullString{String: string(stripeInfoBytes), Valid: true},
		}

		if err = tx.Create(&payment).Error; err != nil {
			return err
		}

		return nil
	})

	return CreatePaymentResponse{PaymentId: paymentId}, models.ErrorResponse{}
}

func validatePayment(request CreatePaymentRequest) []error {
	var validationErrors []error
	if request.UserId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "userId"))
	}

	if len(request.UserId) == 0 || len(request.UserId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "userId"))
	}

	if request.BookingId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "bookingId"))
	}

	if len(request.BookingId) == 0 || len(request.BookingId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "bookingId"))
	}

	return validationErrors
}
