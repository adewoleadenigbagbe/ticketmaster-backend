package services

import (
	"errors"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
)

type GetCityByIdRequest struct {
	Id string `param:"id"`
}

type GetCityByIdResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	IsDeprecated bool   `json:"isDeprecated"`
}

func (cityService CityService) GetCityById(request GetCityByIdRequest) (GetCityByIdResponse, error) {
	var err error
	city := &entities.City{
		Id:           request.Id,
		IsDeprecated: false,
	}

	result := cityService.DB.First(city)
	if result.Error != nil {
		err = errors.New("city record not found")
		return GetCityByIdResponse{}, err
	}

	response := new(GetCityByIdResponse)
	response.Id = city.Id
	response.Name = city.Name
	response.State = city.State
	response.IsDeprecated = city.IsDeprecated

	return *response, nil
}
