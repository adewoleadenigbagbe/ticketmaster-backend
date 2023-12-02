package services

import (
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/models"
)

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Role        int    `json:"role"`
	Description string `json:"description"`
}

type CreateRoleResponse struct {
	UserRoleId string `json:"userRoleId"`
}

func (userService UserService) AddRole(request CreateRoleRequest) (CreateRoleResponse, models.ErrorResponse) {
	requiredFieldErrors := validateRole(request)
	if len(requiredFieldErrors) > 0 {
		return CreateRoleResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: requiredFieldErrors}
	}

	userRole := entities.UserRole{
		Id:          sequentialguid.New().String(),
		Name:        request.Name,
		Description: request.Description,
		Role:        enums.Role(request.Role),
	}

	rowsAffected := userService.DB.Where(entities.UserRole{Role: enums.Role(request.Role)}).FirstOrCreate(&userRole).RowsAffected

	if rowsAffected < 1 {
		err := fmt.Errorf("role already exist")
		return CreateRoleResponse{}, models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}

	return CreateRoleResponse{UserRoleId: userRole.Id}, models.ErrorResponse{}
}

func validateRole(request CreateRoleRequest) []error {
	var validationErrors []error
	if request.Name == "" {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredField, "name"))
	}

	if request.Description == "" {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredField, "description"))
	}

	if request.Role <= 0 {
		validationErrors = append(validationErrors, fmt.Errorf(ErrRequiredField, "role"))
	}

	return validationErrors
}
