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
	userService.Logger.Info().Interface("request", request)
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
		errResponse := models.ErrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
		userService.Logger.Info().Interface("response", errResponse)
		return CreateRoleResponse{}, errResponse
	}

	resp := CreateRoleResponse{UserRoleId: userRole.Id}
	userService.Logger.Info().Interface("response", resp)
	return resp, models.ErrorResponse{}
}

func validateRole(request CreateRoleRequest) []error {
	vErrors := []error{}
	if request.Name == "" {
		vErrors = append(vErrors, fmt.Errorf(ErrRequiredField, "name"))
	}

	if request.Description == "" {
		vErrors = append(vErrors, fmt.Errorf(ErrRequiredField, "description"))
	}

	if request.Role <= 0 {
		vErrors = append(vErrors, fmt.Errorf(ErrRequiredField, "role"))
	}

	return vErrors
}
