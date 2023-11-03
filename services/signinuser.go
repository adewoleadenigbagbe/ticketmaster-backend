package services

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	jwtauth "github.com/Wolechacho/ticketmaster-backend/helpers/utilities/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	StatusCode int
	Token      string `json:"access_token"`
}

func (authService AuthService) SignIn(request SignInRequest) (SignInResponse, []error) {
	var err error

	validationErrors := validateSignInCredentials(request)
	if len(validationErrors) > 0 {
		return SignInResponse{StatusCode: http.StatusBadRequest}, validationErrors
	}
	var user entities.User
	hashedPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	err = authService.DB.Where("Email = ?", request.Email).
		Where("Password = ?", hashedPassword).
		Where("IsDeprecated = ?", false).
		Preload("userroles").
		First(user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return SignInResponse{StatusCode: http.StatusBadRequest}, []error{fmt.Errorf("email or password not found")}
	}

	token, err := jwtauth.GenerateJWT(user)
	if err != nil {
		return SignInResponse{StatusCode: http.StatusBadRequest}, []error{err}
	}

	return SignInResponse{StatusCode: http.StatusOK, Token: token}, nil
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
