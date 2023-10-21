package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
	"gorm.io/gorm"
)

const (
	TIME_OVERLAP_ERROR = "time overlap between the show Start and End Time"
	INVALID_UUID_ERROR = "%s should have a valid UUID"
)

type CreateShowRequest struct {
	ShowTimes    []ShowDateTime `json:"showTimes"`
	CinemaHallId string         `json:"cinemaHallId"`
	MovieId      string         `json:"movieId"`
}

type CreateShowResponse struct {
	ShowIds    []string `json:"showIds"`
	StatusCode int
}

type ShowDateTime struct {
	StartDateTime time.Time `json:"startDate"`
	EndDateTime   time.Time `json:"endDate"`
}

func (showService ShowService) CreateShow(request CreateShowRequest) (CreateShowResponse, []error) {
	var err error
	fieldErrors := validateRequiredFields(request)
	if len(fieldErrors) != 0 {
		return CreateShowResponse{StatusCode: http.StatusBadRequest}, fieldErrors
	}

	showTimeErrors := validateShowTime(request)
	if len(showTimeErrors) != 0 {
		return CreateShowResponse{StatusCode: http.StatusBadRequest}, showTimeErrors
	}

	showIds := []string{}
	err = showService.DB.Transaction(func(tx *gorm.DB) error {
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

			showIds = append(showIds, show.Id)
		}

		return nil
	})

	if err != nil {
		return CreateShowResponse{StatusCode: http.StatusBadRequest}, []error{err}
	}
	return CreateShowResponse{ShowIds: showIds, StatusCode: http.StatusOK}, nil
}

func validateRequiredFields(request CreateShowRequest) []error {
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

func validateShowTime(request CreateShowRequest) []error {

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
