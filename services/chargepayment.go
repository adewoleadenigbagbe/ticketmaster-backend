package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type CreatePaymentRequest struct {
	UserId    string `json:"userId"`
	BookingId string `json:"bookingId"`
}

type CreatePaymentResponse struct {
	StatusCode int
	Errors     []error
}

func (bookService BookService) ChargeBooking(request CreatePaymentRequest) (CreatePaymentResponse, []error) {
	fieldErrors := validatePayment(request)
	if len(fieldErrors) != 0 {
		return CreatePaymentResponse{StatusCode: http.StatusBadRequest}, fieldErrors
	}

	// Set Stripe API key
	stripe.Key = "YOUR_API_KEY"

	// Create a new charge
	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(2000),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("Test Charge"),
	}
	params.SetSource("tok_visa") // use a test card token provided by Stripe
	ch, err := charge.New(params)

	// Check for errors
	if err != nil {
		log.Fatal(err)
	}

	// Print charge details
	//if sucessful, update the show seat and booking rows in the table
	fmt.Printf("Charge ID: %s\n", ch.ID)
	fmt.Printf("Amount: %d\n", ch.Amount)
	fmt.Printf("Description: %s\n", ch.Description)

	return CreatePaymentResponse{StatusCode: http.StatusOK}, nil
}

func validatePayment(request CreatePaymentRequest) []error {
	var validationErrors []error
	if request.UserId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("userId should have a valid UUID"))
	}

	if len(request.UserId) == 0 || len(request.UserId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("userId is a required field with 36 characters"))
	}

	if request.BookingId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("bookingId should have a valid UUID"))
	}

	if len(request.BookingId) == 0 || len(request.BookingId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("bookingId is a required field with 36 characters"))
	}

	return validationErrors

}
