package services

import (
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"gorm.io/gorm"
)

type CinemaRateRequest struct {
	CinemaId   string                  `param:"id"`
	BaseFee    float32                 `json:"baseFee"`
	Discount   utilities.JsonNullFloat `json:"discount"`
	IsSpecials utilities.JsonNullBool  `json:"isSpecials"`
}

type CinemaRateResponse struct {
	Id string `json:"Id"`
}

func (cinemaService CinemaService) AddCinemaRate(request CinemaRateRequest) (CinemaRateResponse, models.ErrorResponse) {
	var err error
	requireFieldErrors := validateCinemaRate(request)
	if len(requireFieldErrors) > 0 {
		return CinemaRateResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: requireFieldErrors}
	}

	cinemaRate := entities.CinemaRate{
		Id:         sequentialguid.New().String(),
		CinemaId:   request.CinemaId,
		BaseFee:    request.BaseFee,
		Discount:   request.Discount.NullFloat64,
		IsSpecials: request.IsSpecials.NullBool,
		IsActive:   true,
	}
	err = cinemaService.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Table("cinemarates").Where("CinemaId = ? AND IsActive = ?", request.CinemaId, true).Update("IsActive", false).Error; err != nil {
			return err
		}

		//add the new rate
		if err = tx.Create(&cinemaRate).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return CinemaRateResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}
	return CinemaRateResponse{Id: cinemaRate.Id}, models.ErrorResponse{}
}

func validateCinemaRate(request CinemaRateRequest) []error {
	validationErrors := []error{}

	if request.CinemaId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf(ErrInvalidUUID, "cinemaId"))
	}

	if len(request.CinemaId) == 0 || len(request.CinemaId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredUUIDField, "cinemaId"))
	}

	if request.BaseFee <= 0 {
		validationErrors = append(validationErrors, fmt.Errorf("base fee is negative"))
	}

	return validationErrors
}
