package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	paginate "github.com/Wolechacho/ticketmaster-backend/helpers/pagination"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	TIME_OVERLAP_ERROR = "time overlap between the show Start and End Time"
	INVALID_UUID_ERROR = "%s should have a valid UUID"
)

type ShowController struct {
}

func (showController *ShowController) CreateShow(showContext echo.Context) error {
	var err error
	request := new(createShowRequest)
	err = showContext.Bind(request)
	fmt.Println(request)
	if err != nil {
		fmt.Println(err)
		return showContext.JSON(http.StatusBadRequest, err.Error())
	}

	fieldErrors := validateRequiredFields(*request)
	if len(fieldErrors) != 0 {
		errors := []string{}
		for _, fieldError := range fieldErrors {
			errors = append(errors, fieldError.Error())
		}
		return showContext.JSON(http.StatusBadRequest, errors)
	}

	showTimeErrors := validateShowTime(*request)
	if len(showTimeErrors) != 0 {
		errors := []string{}
		for _, showTimeError := range showTimeErrors {
			errors = append(errors, showTimeError.Error())
		}
		return showContext.JSON(http.StatusBadRequest, errors)
	}

	response := new(createShowResponse)

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
	request := new(GetShowsByLocationRequest)

	//get show that are not deprecated nor cancelled
	//sort them by the earliest show date and time

	// use the cityId to get the : UserLocation -> cinema -> cinemalHall -> Show -> Movie

	//rows, err := db.DB.Table("userLocations").Select("userLocations.CityId").Joins("join cinemas on userLocations.CityId = cinemas.CityId").Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").Joins("join shows on cinemaHalls.Id = shows.CinemaHallId").Select()
	//.Rows()

	//db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})

	// rows, err := db.DB.Table("userLocations").
	// 	Joins("join cinemas on userLocations.CityId = cinemas.CityId").
	// 	Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
	// 	Joins("join shows on cinemaHalls.Id = shows.CinemaHallId").
	// 	Joins("join movies on shows.MovieId = movies.Id").
	// 	Where("userLocations.UserId = ? AND cinemas.IsDeprecated = ? AND cinemaHalls.IsDeprecated = ? AND shows.IsDeprecated = ? AND movies.IsDeprecated = ?", request.UserId, false, false, false, false).
	// 	Select("shows.Id,shows.Date,shows.StartTime,shows.EndTime,movies.Id,movies.Title,movies.Description,movies.Language,movies.Genre,userLocations.CityId,cinemas.Id").
	// 	Rows()

	//filter
	filterClause := query(*request)

	//paginate
	paginateClause := paginate.Paginate(request.Page, request.PageLength)

	//orderby
	sortandorder := fmt.Sprintf("%s %s", request.SortBy, request.Order)
	orderByClause := paginate.OrderBy(sortandorder)

	rows, err := db.DB.Scopes(filterClause, paginateClause, orderByClause).Rows()

	if err != nil {
		log.Fatalln(err)
	}

	shows := []GetShowsByLocationDTO{}
	for rows.Next() {
		var show GetShowsByLocationDTO
		err = rows.Scan(&show)

		if err != nil {
			log.Fatalln(err)
		}
		shows = append(shows, show)
	}

	response := new(GetShowsByLocationResponse)
	response.Results = shows

	return showContext.JSON(http.StatusOK, response)
}

type GetShowsByLocationRequest struct {
	UserId     string `json:"userId"`
	Page       int    `query:"page"`
	PageLength int    `query:"pageLength"`
	SortBy     string `query:"sortBy"`
	Order      string `query:"order"`
}

type GetShowsByLocationResponse struct {
	Results []GetShowsByLocationDTO
}

type GetShowsByLocationDTO struct {
	Id          string    `json:"id"`
	Date        time.Time `json:"showDate"`
	StartTime   int64     `json:"showStartTime"`
	EndTime     int64     `json:"showEndTime"`
	MovieId     string    `json:"movieId"`
	Title       string    `json:"movieTitle"`
	Description string    `json:"movieDescription"`
	Language    string    `json:"language"`
	Genre       int       `json:"genre"`
}

func query(filter GetShowsByLocationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("userLocations").
			Joins("join cinemas on userLocations.CityId = cinemas.CityId").
			Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
			Joins("join shows on cinemaHalls.Id = shows.CinemaHallId").
			Joins("join movies on shows.MovieId = movies.Id").
			Where("userLocations.UserId = ? AND cinemas.IsDeprecated = ? AND cinemaHalls.IsDeprecated = ? AND shows.IsDeprecated = ? AND shows.IsCancelled = ? AND movies.IsDeprecated = ?",
				filter.UserId, false, false, false, false, false).
			Select("shows.Id,shows.Date,shows.StartTime,shows.EndTime,movies.Id,movies.Title,movies.Description,movies.Language,movies.Genre,userLocations.CityId,cinemas.Id")
	}
}

func validateRequiredFields(request createShowRequest) []error {
	var validationErrors []error
	defaultTime, _ := time.Parse("2006-01-02", entities.MIN_DATE)

	//validate the cinemaHallId and movieId
	if len(request.CinemaHallId) == 0 || len(request.CinemaHallId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("cinemaHallId is a required field  with 36 characters"))
	}

	if request.CinemaHallId == entities.DEFAULT_UUID {
		validationErrors = append(validationErrors, fmt.Errorf("cinemaHallId should have a valid UUID"))
	}

	if len(request.MovieId) == 0 || len(request.MovieId) < 36 {
		validationErrors = append(validationErrors, fmt.Errorf("movieId is a required field with 36 characters"))
	}

	if request.MovieId == entities.DEFAULT_UUID {
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
	defaultTime, _ := time.Parse("2006-01-02", entities.MIN_DATE)

	timeOverlap := false
	tempStartDate, _ := time.Parse("2006-01-02", entities.MAX_DATE)
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
			errorMessage := fmt.Sprintf("showTimes[%d].endDate: %s must be greater than showTimes[%d].startDate: %s", i, showTime.EndDateTime.Format("2006-01-02T15:04:05"), i, showTime.StartDateTime.Format("2006-01-02T15:04:05"))
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
