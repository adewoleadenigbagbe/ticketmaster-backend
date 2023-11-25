package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"gorm.io/gorm"
)

type UserLocationRequest struct {
	UserId    string  `query:"id"`
	CityId    string  `json:"cityId"`
	Address   string  `json:"address"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type UserLocationResponse struct {
	AddressId string `json:"addressId"`
}

func (userService UserService) UpdateUserLocation(request UserLocationRequest) (UserLocationResponse, models.ErrorResponse) {
	requiredFieldErrors := validateUserLocation(request)
	if len(requiredFieldErrors) > 0 {
		return UserLocationResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: requiredFieldErrors}
	}

	var err error
	user := entities.User{IsDeprecated: false, Id: request.UserId}

	err = userService.DB.First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return UserLocationResponse{}, models.ErrorResponse{StatusCode: http.StatusNotFound, Errors: []error{errors.New("user not found")}}
	}

	address := entities.Address{
		Id:          sequentialguid.New().String(),
		EntityId:    user.Id,
		AddressType: enums.User,
		AddressLine: request.Address,
		CityId:      request.CityId,
		Coordinates: entities.Coordinate{
			Longitude: request.Longitude,
			Latitude:  request.Latitude,
		},
		IsDeprecated: false,
		IsCurrent:    true,
	}

	err = userService.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Table("addresses").Where("EntityId = ? AND AddressType = ? AND IsCurrent = ?", user.Id, enums.User, true).Update("IsCurrent", false).Error; err != nil {
			return err
		}

		//add the new address
		if err = tx.Create(&address).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return UserLocationResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}

	return UserLocationResponse{AddressId: address.Id}, models.ErrorResponse{}
}

func validateUserLocation(request UserLocationRequest) []error {
	validationErrors := []error{}

	if request.UserId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("userId should have a valid UUID"))
	}

	if len(request.UserId) == 0 || len(request.UserId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("userId is a required field with 36 characters"))
	}

	if request.CityId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("cityId should have a valid UUID"))
	}

	if len(request.CityId) == 0 || len(request.CityId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("cityId is a required field with 36 characters"))
	}

	if request.Address == "" {
		validationErrors = append(validationErrors, fmt.Errorf("address is a required field"))
	}

	return validationErrors
}
