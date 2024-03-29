package services

import (
	"errors"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"gorm.io/gorm"
)

type GetCityByIdRequest struct {
	Id string `param:"id"`
}

type CityModelResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	IsDeprecated bool   `json:"isDeprecated"`
}
type GetCityByIdResponse struct {
	City       CityModelResponse
	StatusCode int
}

func (cityService CityService) GetCityById(request GetCityByIdRequest) (GetCityByIdResponse, error) {
	var err error
	city := entities.City{}

	result := cityService.DB.Where("Id = ? AND IsDeprecated = ?", request.Id, false).First(&city)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("city record not found")
		return GetCityByIdResponse{StatusCode: http.StatusNotFound}, err
	}

	cityResp := CityModelResponse{
		Id:           city.Id,
		Name:         city.Name,
		State:        city.State,
		IsDeprecated: false,
	}
	return GetCityByIdResponse{City: cityResp, StatusCode: http.StatusOK}, nil
}
