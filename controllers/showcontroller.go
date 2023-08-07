package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ShowController struct {
}

func (showController *ShowController) CreateShow(showContext echo.Context) error {
	request := new(createMovieRequest)
	err := showContext.Bind(request)
	if err != nil {
		return showContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	fmt.Println(request)

	//validate the request
	validate := validator.New()
	err = validate.Struct(request)
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return showContext.JSON(http.StatusBadRequest, validationErrors)
	}

	start := request.StartDateTime.Unix()
	end := request.EndDateTime.Unix()

	//custom validation to check the duration of the show using the start date and the end date
	//custom validation to check whether the end date is greater than the start date

	show := &entities.Show{
		Id:                 sequentialguid.New().String(),
		Date:               request.StartDateTime,
		StartTime:          start,
		EndTime:            end,
		MovieId:            request.MovieId,
		CinemaHallId:       request.CinemallHallId,
		IsDeprecated:       false,
		IsCancelled:        false,
		CancellationReason: sql.NullString{Valid: false},
	}

	result := db.DB.Create(&show)
	if result.Error != nil || result.RowsAffected < 1 {
		return showContext.JSON(http.StatusBadRequest, result.Error)
	}

	response := new(createMovieResponse)
	response.ShowId = show.Id
	return showContext.JSON(http.StatusOK, response)
}

type createMovieRequest struct {
	StartDateTime  time.Time `json:"startDate" validate:"required;isdefault;datetime=2006-01-02;gt"`
	EndDateTime    time.Time `json:"endDate" validate:"required;isdefault;datetime=2006-01-02;gt;"`
	CinemallHallId string    `json:"cinemaHallId" validate:"required;uuid;isdefault"`
	MovieId        string    `json:"movieId" validate:"required;uuid;isdefault"`
}

type createMovieResponse struct {
	ShowId string `json:"showId"`
}
