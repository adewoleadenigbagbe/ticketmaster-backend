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

func (userService UserService) AddRole(request CreateRoleRequest) (CreateRoleResponse, models.ErrrorResponse) {
	requiredFieldErrors := validateRole(request)
	if len(requiredFieldErrors) > 0 {
		return CreateRoleResponse{}, models.ErrrorResponse{StatusCode: http.StatusBadRequest, Errors: requiredFieldErrors}
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
		return CreateRoleResponse{}, models.ErrrorResponse{StatusCode: http.StatusBadRequest, Errors: []error{err}}
	}

	return CreateRoleResponse{UserRoleId: userRole.Id}, models.ErrrorResponse{}
}

func validateRole(request CreateRoleRequest) []error {
	vErrors := []error{}
	if request.Name == "" {
		vErrors = append(vErrors, fmt.Errorf("name field is required"))
	}

	if request.Description == "" {
		vErrors = append(vErrors, fmt.Errorf("description field is required"))
	}

	if request.Role <= 0 {
		vErrors = append(vErrors, fmt.Errorf("role field is required"))
	}

	return vErrors
}
