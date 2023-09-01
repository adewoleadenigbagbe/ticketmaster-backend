package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/labstack/echo/v4"
)

type ShowController struct {
}

func (showController *ShowController) CreateShow(showContext echo.Context) error {
	request := new(createMovieRequest)
	err := showContext.Bind(request)

	if err != nil {
		fmt.Println(err)
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	validationErrors := validateShow(*request)
	if len(validationErrors) != 0 {
		return showContext.JSON(http.StatusBadRequest, validationErrors)
	}

	start := request.StartDateTime.Unix()
	end := request.EndDateTime.Unix()

	show := &entities.Show{
		Id:                 sequentialguid.New().String(),
		Date:               request.StartDateTime,
		StartTime:          start,
		EndTime:            end,
		MovieId:            request.MovieId,
		CinemaHallId:       request.CinemaHallId,
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
	StartDateTime time.Time `json:"startDate"`
	EndDateTime   time.Time `json:"endDate"`
	CinemaHallId  string    `json:"cinemaHallId"`
	MovieId       string    `json:"movieId"`
}

type createMovieResponse struct {
	ShowId string `json:"showId"`
}

func validateShow(request createMovieRequest) []string {
	errors := []string{}

	defaultTime, _ := time.Parse("2006-01-02", entities.MIN_DATE)
	today := time.Now().Local()

	if request.StartDateTime == defaultTime {
		errors = append(errors, "startDate is a required field")
	}

	if request.EndDateTime == defaultTime {
		errors = append(errors, "endDate is a required field")
	}

	if request.StartDateTime != defaultTime && request.StartDateTime.Before(today) {
		errors = append(errors, "show startdate can be added for future dates")
	}

	if request.EndDateTime != defaultTime && request.EndDateTime.Before(today) {
		errors = append(errors, "show enddate can be added for future dates")
	}

	if request.StartDateTime.Equal(request.EndDateTime) {
		errors = append(errors, "show cannot start and end at the same time")
	}

	if request.StartDateTime.After(request.EndDateTime) {
		errorMessage := fmt.Sprintf("show time end time %s must be greater than the start time %s", request.EndDateTime.Format("2006-01-02T15:04:05"), request.StartDateTime.Format("2006-01-02T15:04:05"))
		errors = append(errors, errorMessage)
	}

	var minShowTime float64 = 1
	var maxShowTime float64 = 4
	hours := request.EndDateTime.Sub(request.StartDateTime).Hours()
	if hours < minShowTime || hours > maxShowTime {
		errors = append(errors, fmt.Sprintf("Show time must be between %.0fhrs and %.0fhrs", minShowTime, maxShowTime))
	}

	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		errors = append(errors, "cinemaHallId is a required field  with 36 characters")
	}

	if request.CinemaHallId == entities.DEFAULT_UUID {
		errors = append(errors, "cinemaHallId should have a valid UUID")
	}

	if len(request.MovieId) == 0 || len(request.MovieId) < 36 {
		errors = append(errors, "movieId is a required field with 36 characters")
	}

	if request.MovieId == entities.DEFAULT_UUID {
		errors = append(errors, "movieId should have a valid UUID")
	}

	return errors
}
