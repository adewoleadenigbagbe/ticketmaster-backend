package services

import (
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/shared/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/shared/helpers"
	"github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/shared/models"
	"gorm.io/gorm"
)

const (
	MaxDiscount = 0.3
)

type CinemaRateRequest struct {
	CinemaId   string                      `param:"id"`
	BaseFee    float32                     `json:"baseFee"`
	Discount   utilities.Nullable[float64] `json:"discount"`
	IsSpecials utilities.Nullable[bool]    `json:"isSpecials"`
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
		Discount:   request.Discount,
		IsSpecials: request.IsSpecials,
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

	if request.Discount.Valid && request.IsSpecials.Valid && request.Discount.Val > MaxDiscount {
		validationErrors = append(validationErrors, fmt.Errorf("discount is invalid.should not be over 30 percent"))
	}

	return validationErrors
}