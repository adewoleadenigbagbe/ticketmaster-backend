package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/labstack/echo/v4"
)

const (
	EmailRegex       = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
	PhoneNumberRegex = "\\+[1-9]{1}[0-9]{0,2}-[2-9]{1}[0-9]{2}-[2-9]{1}[0-9]{2}-[0-9]{4}$"
)

type UserController struct {
}

func (userController UserController) CreateUser(userContext echo.Context) error {
	var err error

	request := new(createUserRequest)

	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	fieldsErrors := validateUser(*request)
	if len(fieldsErrors) != 0 {
		errors := []string{}
		for _, err = range fieldsErrors {
			errors = append(errors, err.Error())
		}
		return userContext.JSON(http.StatusBadRequest, errors)
	}

	user := entities.User{
		Id:           sequentialguid.New().String(),
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Email:        request.Email,
		PhoneNumber:  sql.NullString{String: request.PhoneNumber, Valid: false},
		IsDeprecated: false,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		return userContext.JSON(http.StatusBadRequest, result.Error.Error())
	}

	response := new(createUserResponse)
	response.UserId = user.Id
	return userContext.JSON(http.StatusOK, response)
}

type createUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type createUserResponse struct {
	UserId string `json:"userId"`
}

func validateUser(request createUserRequest) []error {
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

	var err error
	_, err = regexp.MatchString(EmailRegex, request.Email)
	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("email supplied is invalid"))
	}

	_, err = regexp.MatchString(PhoneNumberRegex, request.PhoneNumber)
	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("phone number supplied is invalid"))
	}

	return validationErrors
}
