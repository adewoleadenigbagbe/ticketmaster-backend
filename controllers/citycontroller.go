package controllers

import (
	"net/http"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/labstack/echo/v4"
)

type CityController struct {
}

func (cityController *CityController) GetCityById(cityContext echo.Context) error {
	req := new(getCityByIdRequest)

	err := cityContext.Bind(req)
	if err != nil {
		return cityContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	city := &entities.City{
		Id:           req.Id,
		IsDeprecated: false,
	}
	result := db.DB.First(city)
	if result.Error != nil {
		return cityContext.JSON(http.StatusNotFound, "Movie Record not found")
	}

	response := new(getCityByIdResponse)
	response.Id = city.Id
	response.Name = city.Name
	response.State = city.State
	response.IsDeprecated = city.IsDeprecated

	return cityContext.JSON(http.StatusOK, response)
}

type getCityByIdRequest struct {
	Id string `param:"id"`
}

type getCityByIdResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	IsDeprecated bool   `json:"isDeprecated"`
}
