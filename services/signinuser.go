package services

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	jwtauth "github.com/Wolechacho/ticketmaster-backend/helpers/utilities/auth"
	"github.com/Wolechacho/ticketmaster-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"access_token"`
}

func (authService AuthService) SignIn(request SignInRequest) (SignInResponse, models.ErrorResponse) {
	var err error

	validationErrors := validateSignInCredentials(request)
	if len(validationErrors) > 0 {
		return SignInResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: validationErrors}
	}

	var user *entities.User
	err = authService.DB.Where("Email = ?", request.Email).
		Where("IsDeprecated = ?", false).
		Preload("UserRole").
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return SignInResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{fmt.Errorf("email or password not found")}}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return SignInResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{fmt.Errorf("email or password not found")}}
	}

	token, err := jwtauth.GenerateJWT(*user)
	if err != nil {
		return SignInResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}

	return SignInResponse{Token: token}, models.ErrorResponse{}
}

func validateSignInCredentials(request SignInRequest) []error {
	var validationErrors []error
	if request.Password == "" {
		validationErrors = append(validationErrors, fmt.Errorf("password is a required field"))
	}

	isEmailValid, _ := regexp.MatchString(EmailRegex, request.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, fmt.Errorf("email supplied is invalid"))
	}

	return validationErrors
}
