package controllers

import (
	"time"

	"github.com/labstack/echo/v4"
)

type ShowController struct {
}

func (showController ShowController) GetShowsByUserLocation(showContext echo.Context) error {
	//get show that are not deprecated nor cancelled
	//sort them by the earlist show date and time
	// UserId -> User -> CityId -> Cinema ->

	// use the cityId to get the : UserLocation -> cinema -> cinemalHall -> Show -> Movie

	return nil
}

type GetShowsRequest struct {
	UserId     string `json:"userId"`
	Page       int    `query:"page"`
	PageLength int    `query:"pageLength"`
	SortBy     string `query:"sortBy"`
	Order      string `query:"order"`
}

type GetShowsResponse struct {
	Id           string    `json:"id"`
	Date         time.Time `json:"showDate"`
	StartTime    int64     `json:"showStartTime"`
	EndTime      int64     `json:"showEndTime"`
	MovieId      string    `json:"movieId"`
	Title        string    `json:"movieTitle"`
	Description  string    `json:"movieDescription"`
	Language     string    `json:"language"`
	Genre        int       `json:"genre"`
	IsDeprecated bool
}
