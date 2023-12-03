package services

import (
	"errors"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/models"
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
	City CityModelResponse
}

func (cityService CityService) GetCityById(request GetCityByIdRequest) (GetCityByIdResponse, models.ErrorResponse) {
	cityService.Logger.Info().Interface("getCityByIdRequest", request).Msg("request")
	var err error
	city := entities.City{}

	result := cityService.DB.Where("Id = ? AND IsDeprecated = ?", request.Id, false).First(&city)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("city record not found")
		cityService.Logger.Info().Interface("getCityByIdResponse", err.Error()).Msg("response")
		return GetCityByIdResponse{}, models.ErrorResponse{StatusCode: http.StatusNotFound, Errors: []error{err}}
	}

	cityResp := CityModelResponse{
		Id:           city.Id,
		Name:         city.Name,
		State:        city.State,
		IsDeprecated: false,
	}
	cityService.Logger.Info().Interface("getCityByIdResponse", cityResp).Msg("response")
	return GetCityByIdResponse{City: cityResp}, models.ErrorResponse{}
}
