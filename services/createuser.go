package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

const (
	EmailRegex       = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
	PhoneNumberRegex = "\\+[1-9]{1}[0-9]{0,2}-[2-9]{1}[0-9]{2}-[2-9]{1}[0-9]{2}-[0-9]{4}$"
)

type CreateUserRequest struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phoneNumber"`
	CityId      string  `json:"cityId"`
	Address     string  `json:"address"`
	Longitude   float32 `json:"longitude"`
	Latitude    float32 `json:"latitude"`
}

type CreateUserResponse struct {
	UserId     string  `json:"userId"`
	Errors     []error `json:"errors"`
	StatusCode int     `json:"statusCode"`
}

func validateUser(request CreateUserRequest) []error {
	var validationErrors []error
	if request.FirstName == "" {
		validationErrors = append(validationErrors, fmt.Errorf("firstName is a required field"))
	}

	if request.LastName == "" {
		validationErrors = append(validationErrors, fmt.Errorf("lastName is a required field"))
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, fmt.Errorf("password is a required field"))
	}

	if request.Address == "" {
		validationErrors = append(validationErrors, fmt.Errorf("address is a required field"))
	}

	if len(request.CityId) == 0 || len(request.CityId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("cityId is a required field  with 36 characters"))
	}

	if request.CityId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("cityId should have a valid UUID"))
	}

	isEmailValid, _ := regexp.MatchString(EmailRegex, request.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, fmt.Errorf("email supplied is invalid"))
	}

	isPhoneValid, _ := regexp.MatchString(PhoneNumberRegex, request.PhoneNumber)

	if !isPhoneValid {
		validationErrors = append(validationErrors, fmt.Errorf("phone number supplied is invalid"))
	}

	return validationErrors
}

func (userService UserService) CreateUser(request CreateUserRequest) CreateUserResponse {
	var err error
	fieldsErrors := validateUser(request)
	if len(fieldsErrors) != 0 {
		return CreateUserResponse{Errors: fieldsErrors, StatusCode: http.StatusBadRequest}
	}

	user := entities.User{
		Id:           sequentialguid.New().String(),
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Email:        request.Email,
		Password:     request.Password,
		PhoneNumber:  sql.NullString{String: request.PhoneNumber, Valid: true},
		IsDeprecated: false,
	}

	err = userService.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			// return any error will rollback
			return err
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
		}

		if err := tx.Create(&address).Error; err != nil {
			// return any error will rollback
			return err
		}

		return nil
	})

	if err != nil {
		return CreateUserResponse{Errors: []error{err}, StatusCode: http.StatusBadRequest}
	}

	return CreateUserResponse{UserId: user.Id, Errors: []error{}, StatusCode: http.StatusOK}
}
