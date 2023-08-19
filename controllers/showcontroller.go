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
	"gorm.io/gorm"
)

const (
	TIME_OVERLAP_ERROR = "Time overlap between the show Start and End Time"
	INVALID_UUID_ERROR = "%s should have a valid UUID"
)

type ShowController struct {
}

func (showController *ShowController) CreateShow(showContext echo.Context) error {
	var err error
	request := new(createMovieRequest)
	err = showContext.Bind(request)

	if err != nil {
		fmt.Println(err)
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	validationErrors := validateShow(*request)
	if len(validationErrors) != 0 {
		return showContext.JSON(http.StatusBadRequest, validationErrors)
	}

	response := new(createMovieResponse)

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		for _, showTime := range request.ShowTimes {
			show := &entities.Show{
				Id:                 sequentialguid.New().String(),
				Date:               showTime.StartDateTime,
				StartTime:          showTime.StartDateTime.Unix(),
				EndTime:            showTime.EndDateTime.Unix(),
				MovieId:            request.MovieId,
				CinemaHallId:       request.CinemaHallId,
				IsDeprecated:       false,
				IsCancelled:        false,
				CancellationReason: sql.NullString{Valid: false},
			}

			result := tx.Create(&show)
			if result.Error != nil || result.RowsAffected < 1 {
				return result.Error
			}

			response.ShowIds = append(response.ShowIds, show.Id)
		}

		return nil
	})

	if err != nil {
		return showContext.JSON(http.StatusBadRequest, err)
	}

	return showContext.JSON(http.StatusOK, response)
}

type createMovieRequest struct {
	ShowTimes    []ShowDateTime `json:"showTimes"`
	CinemaHallId string         `json:"cinemaHallId"`
	MovieId      string         `json:"movieId"`
}
type ShowDateTime struct {
	StartDateTime time.Time `json:"startDate"`
	EndDateTime   time.Time `json:"endDate"`
}

type createMovieResponse struct {
	ShowIds []string `json:"showIds"`
}

func validateShow(request createMovieRequest) []string {
	errors := []string{}

	defaultTime, _ := time.Parse("2006-01-02", entities.DEFAULT_TIME)
	today := time.Now().Local()
	var minShowTime float64 = 1
	var maxShowTime float64 = 4

	var timeOverlap = false
	var tempStartDate = today
	var tempEndDate = today

	//Validate the show time
	for i, showTime := range request.ShowTimes {
		if showTime.StartDateTime == defaultTime {
			errors = append(errors, "startDate is a required field")
		}

		if showTime.EndDateTime == defaultTime {
			errors = append(errors, "endDate is a required field")
		}

		if showTime.StartDateTime != defaultTime && showTime.StartDateTime.Before(today) {
			errors = append(errors, "show startdate can be added for future dates")
		}

		if showTime.EndDateTime != defaultTime && showTime.EndDateTime.Before(today) {
			errors = append(errors, "show enddate can be added for future dates")
		}

		if showTime.StartDateTime.Equal(showTime.EndDateTime) {
			errors = append(errors, "show cannot start and end at the same time")
		}

		if showTime.StartDateTime.After(showTime.EndDateTime) {
			errorMessage := fmt.Sprintf("show time end time %s must be greater than the start time %s", showTime.EndDateTime.Format("2006-01-02T15:04:05"), showTime.StartDateTime.Format("2006-01-02T15:04:05"))
			errors = append(errors, errorMessage)
		}

		hours := showTime.EndDateTime.Sub(showTime.StartDateTime).Hours()
		if hours < minShowTime || hours > maxShowTime {
			errors = append(errors, fmt.Sprintf("Show time must be between %.0fhrs and %.0fhrs", minShowTime, maxShowTime))
		}

		//check for date overlap
		if i != 0 {
			timeOverlap = tempStartDate.Before(showTime.EndDateTime) && tempEndDate.After(showTime.StartDateTime)
			if timeOverlap {
				continue
			}
		}

		tempStartDate = showTime.StartDateTime
		tempEndDate = showTime.EndDateTime
	}

	//if there is overlap, add error
	if timeOverlap {
		errors = append(errors, TIME_OVERLAP_ERROR)
	}

	//validate the cinemaHallId and movieId
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
