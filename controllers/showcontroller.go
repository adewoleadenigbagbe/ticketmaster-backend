package controllers

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/core"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"github.com/labstack/echo/v4"
)

const (
	TIME_OVERLAP_ERROR = "time overlap between the show Start and End Time"
	INVALID_UUID_ERROR = "%s should have a valid UUID"
)

type ShowController struct {
	App *core.BaseApp
}

func (showController *ShowController) CreateShow(showContext echo.Context) error {
	// var err error
	// request := new(createShowRequest)
	// err = showContext.Bind(request)
	// fmt.Println(request)
	// if err != nil {
	// 	return showContext.JSON(http.StatusBadRequest, err.Error())
	// }

	// fieldErrors := validateRequiredFields(*request)
	// if len(fieldErrors) != 0 {
	// 	errors := []string{}
	// 	for _, fieldError := range fieldErrors {
	// 		errors = append(errors, fieldError.Error())
	// 	}
	// 	return showContext.JSON(http.StatusBadRequest, errors)
	// }

	// showTimeErrors := validateShowTime(*request)
	// if len(showTimeErrors) != 0 {
	// 	errors := []string{}
	// 	for _, showTimeError := range showTimeErrors {
	// 		errors = append(errors, showTimeError.Error())
	// 	}
	// 	return showContext.JSON(http.StatusBadRequest, errors)
	// }

	response := new(createShowResponse)

	// err = db.DB.Transaction(func(tx *gorm.DB) error {
	// 	for _, showTime := range request.ShowTimes {
	// 		show := &entities.Show{
	// 			Id:                 sequentialguid.New().String(),
	// 			Date:               showTime.StartDateTime,
	// 			StartTime:          showTime.StartDateTime.Unix(),
	// 			EndTime:            showTime.EndDateTime.Unix(),
	// 			MovieId:            request.MovieId,
	// 			CinemaHallId:       request.CinemaHallId,
	// 			IsDeprecated:       false,
	// 			IsCancelled:        false,
	// 			CancellationReason: sql.NullString{Valid: false},
	// 		}

	// 		result := tx.Create(&show)
	// 		if result.Error != nil || result.RowsAffected < 1 {
	// 			return result.Error
	// 		}

	// 		response.ShowIds = append(response.ShowIds, show.Id)
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	return showContext.JSON(http.StatusBadRequest, err)
	// }

	return showContext.JSON(http.StatusOK, response)
}

type createShowRequest struct {
	ShowTimes    []ShowDateTime `json:"showTimes"`
	CinemaHallId string         `json:"cinemaHallId"`
	MovieId      string         `json:"movieId"`
}
type ShowDateTime struct {
	StartDateTime time.Time `json:"startDate"`
	EndDateTime   time.Time `json:"endDate"`
}

type createShowResponse struct {
	ShowIds []string `json:"showIds"`
}

func (showController ShowController) GetShowsByUserLocation(showContext echo.Context) error {
	// var err error
	// request := new(GetShowsByLocationRequest)

	// err = showContext.Bind(request)
	// if err != nil {
	// 	return showContext.JSON(http.StatusBadRequest, err.Error())
	// }

	// //get show that are not deprecated nor cancelled
	// //sort them by the earliest show date and time

	// userQuery, err := db.DB.Table("users").
	// 	Where("users.Id = ?", request.UserId).
	// 	Where("users.IsDeprecated = ?", false).
	// 	Joins("join addresses on users.Id = addresses.EntityId").
	// 	Where("addresses.EntityId = ?", request.UserId).
	// 	Where("addresses.IsDeprecated = ?", false).
	// 	Where("addresses.AddressType = ?", enums.User).
	// 	Select("users.Id AS UserId, users.IsDeprecated, addresses.CityId,addresses.Coordinates").
	// 	Rows()

	// defer userQuery.Close()
	// if err != nil {
	// 	showContext.JSON(http.StatusInternalServerError, err.Error())
	// }

	// var user UserDTO
	// i := 0
	// for userQuery.Next() {
	// 	if i > 1 {
	// 		break
	// 	}
	// 	err = userQuery.Scan(&user.UserId, &user.IsDeprecated, &user.CityId, &user.Coordinates)
	// 	if err != nil {
	// 		return showContext.JSON(http.StatusInternalServerError, err.Error())
	// 	}
	// 	i++
	// }

	// showQuery, err := db.DB.Table("addresses").
	// 	Where("addresses.CityId = ?", user.CityId).
	// 	Where("addresses.IsDeprecated = ?", false).
	// 	Where("addresses.AddressType = ?", enums.Cinema).
	// 	Joins("join cinemas on addresses.EntityId = cinemas.Id").
	// 	Where("cinemas.IsDeprecated = ?", false).
	// 	Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
	// 	Where("cinemaHalls.IsDeprecated = ?", false).
	// 	Joins("join shows on cinemaHalls.Id = shows.CinemaHallId").
	// 	Where("shows.IsDeprecated = ?", false).
	// 	Where("shows.IsCancelled = ?", false).
	// 	Joins("join movies on shows.MovieId = movies.Id").
	// 	Where("movies.IsDeprecated = ?", false).
	// 	Select("shows.Id AS ShowId, shows.Date, shows.StartTime, shows.EndTime,movies.Id AS MovieId, movies.Title, movies.Description, movies.Language, movies.Genre,addresses.AddressLine, addresses.Coordinates").
	// 	Rows()

	// defer showQuery.Close()
	// if err != nil {
	// 	return showContext.JSON(http.StatusInternalServerError, err.Error())
	// }

	// shows := []ShowsDTO{}
	// for showQuery.Next() {
	// 	show := &ShowsDTO{}
	// 	err = showQuery.Scan(&show.ShowId, &show.Date, &show.startTime, &show.endTime, &show.MovieId, &show.Title,
	// 		&show.Description, &show.Language, &show.Genre, &show.AddressLine, &show.Coordinates)

	// 	if err != nil {
	// 		return showContext.JSON(http.StatusInternalServerError, err.Error())
	// 	}

	// 	show.StartTime = time.Unix(show.startTime, 0)
	// 	show.EndTime = time.Unix(show.endTime, 0)

	// 	show.Distance = distance(user.Coordinates, show.Coordinates)

	// 	shows = append(shows, *show)
	// }

	response := new(GetShowsByLocationResponse)
	//response.Results = shows

	//fmt.Println(shows)
	return showContext.JSON(http.StatusOK, response)
}

type GetShowsByLocationRequest struct {
	UserId string `json:"userId"`
}

type GetShowsByLocationResponse struct {
	Results []ShowsDTO
}

type ShowsDTO struct {
	ShowId      string    `json:"showId"`
	Date        time.Time `json:"showDate"`
	startTime   int64
	endTime     int64
	MovieId     string              `json:"movieId"`
	Title       string              `json:"movieTitle"`
	Description string              `json:"movieDescription"`
	Language    string              `json:"language"`
	Genre       int                 `json:"genre"`
	StartTime   time.Time           `json:"showStartTime"`
	EndTime     time.Time           `json:"showEndTime"`
	AddressLine string              `json:"address"`
	Coordinates entities.Coordinate `json:"coordinates"`
	Distance    float64             `json:"distance"`
}

type UserDTO struct {
	UserId       string
	IsDeprecated bool
	Coordinates  entities.Coordinate
	CityId       string
}

func validateRequiredFields(request createShowRequest) []error {
	var validationErrors []error
	defaultTime, _ := time.Parse(time.RFC3339, utilities.MIN_DATE)

	//validate the cinemaHallId and movieId
	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("cinemaHallId is a required field  with 36 characters"))
	}

	if request.CinemaHallId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("cinemaHallId should have a valid UUID"))
	}

	if len(request.MovieId) == 0 || len(request.MovieId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("movieId is a required field with 36 characters"))
	}

	if request.MovieId == utilities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("movieId should have a valid UUID"))
	}

	if len(request.ShowTimes) == 0 {
		validationErrors = append(validationErrors, fmt.Errorf("show times is empty"))
	}

	//Validate the show time
	for i, showTime := range request.ShowTimes {
		if showTime.StartDateTime == defaultTime {
			validationErrors = append(validationErrors, fmt.Errorf("showTimes[%d].startDate is a required field", i))
		}

		if showTime.EndDateTime == defaultTime {
			validationErrors = append(validationErrors, fmt.Errorf("showTimes[%d].endDate is a required field", i))
		}
	}
	return validationErrors
}

func validateShowTime(request createShowRequest) []error {

	var minShowTime float64 = 1
	var maxShowTime float64 = 4
	var validationErrors []error

	today := time.Now().Local()
	defaultTime, _ := time.Parse(time.RFC3339, utilities.MIN_DATE)

	timeOverlap := false
	tempStartDate, _ := time.Parse(time.RFC3339, utilities.MAX_DATE)
	tempEndDate := tempStartDate

	//Validate the show time
	for i, showTime := range request.ShowTimes {
		if showTime.StartDateTime != defaultTime && showTime.StartDateTime.Before(today) {
			validationErrors = append(validationErrors, fmt.Errorf("showTimes[%d].startDate can only be added for future dates", i))
		}

		if showTime.EndDateTime != defaultTime && showTime.EndDateTime.Before(today) {
			validationErrors = append(validationErrors, fmt.Errorf("showTimes[%d].endDate can only be added for future dates", i))
		}

		if showTime.StartDateTime.Equal(showTime.EndDateTime) {
			validationErrors = append(validationErrors, fmt.Errorf("showTimes[%d].startDate and showTimes[%d].endDate cannot start and end at the same time", i, i))
		}

		if showTime.StartDateTime.After(showTime.EndDateTime) {
			errorMessage := fmt.Sprintf("showTimes[%d].endDate: %s must be greater than showTimes[%d].startDate: %s", i, showTime.EndDateTime.Format(time.RFC3339), i, showTime.StartDateTime.Format(time.RFC3339))
			validationErrors = append(validationErrors, fmt.Errorf(errorMessage))
		}

		hours := showTime.EndDateTime.Sub(showTime.StartDateTime).Hours()
		if hours < minShowTime || hours > maxShowTime {
			validationErrors = append(validationErrors, fmt.Errorf("showTimes[%d] must be between %.0fhrs and %.0fhrs", i, minShowTime, maxShowTime))
		}

		//check for date overlap
		if timeOverlap {
			continue
		}

		timeOverlap = tempStartDate.Before(showTime.EndDateTime) && tempEndDate.After(showTime.StartDateTime)
		tempStartDate = showTime.StartDateTime
		tempEndDate = showTime.EndDateTime
	}

	//if there is overlap, add error
	fmt.Println("Overlap : ", timeOverlap)
	if timeOverlap {
		validationErrors = append(validationErrors, fmt.Errorf(TIME_OVERLAP_ERROR))
	}

	return validationErrors
}

func distance(coordinate1, coordinate2 entities.Coordinate, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * float64(coordinate1.Latitude) / 180)
	radlat2 := float64(PI * float64(coordinate2.Latitude) / 180)

	theta := float64(float64(coordinate1.Longitude) - float64(coordinate2.Longitude))
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}
	return dist
}
