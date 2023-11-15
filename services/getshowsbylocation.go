package services

import (
	"math"
	"net/http"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/enums"
)

type GetShowsByLocationRequest struct {
	UserId string `json:"userId"`
}

type GetShowsByLocationResponse struct {
	Results    []ShowsDTO
	StatusCode int
}

type ShowsDTO struct {
	ShowId      string    `json:"showId"`
	Date        time.Time `json:"showDate"`
	startTime   int64
	endTime     int64
	MovieId     string    `json:"movieId"`
	Title       string    `json:"movieTitle"`
	Description string    `json:"movieDescription"`
	Language    string    `json:"language"`
	Genre       int       `json:"genre"`
	StartTime   time.Time `json:"showStartTime"`
	EndTime     time.Time `json:"showEndTime"`
	AddressLine string    `json:"address"`
	//Coordinates entities.Coordinate `json:"coordinates"`
	Distance float64 `json:"distance"`
}

type UserDTO struct {
	UserId       string
	IsDeprecated bool
	Coordinates  entities.Coordinate
	CityId       string
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

func (showService ShowService) GetShowsByUserLocation(request GetShowsByLocationRequest) (GetShowsByLocationResponse, error) {

	//get show that are not deprecated nor cancelled
	//sort them by the earliest show date and time
	userQuery, err := showService.DB.Table("users").
		Where("users.Id = ?", request.UserId).
		Where("users.IsDeprecated = ?", false).
		Joins("join addresses on users.Id = addresses.EntityId").
		Where("addresses.EntityId = ?", request.UserId).
		Where("addresses.IsDeprecated = ?", false).
		Where("addresses.AddressType = ?", enums.User).
		Select("users.Id AS UserId, users.IsDeprecated, addresses.CityId,addresses.Coordinates").
		Rows()

	if err != nil {
		return GetShowsByLocationResponse{StatusCode: http.StatusInternalServerError}, err
	}

	defer userQuery.Close()

	var user UserDTO
	i := 0
	for userQuery.Next() {
		if i > 1 {
			break
		}
		err = userQuery.Scan(&user.UserId, &user.IsDeprecated, &user.CityId, &user.Coordinates)
		if err != nil {
			return GetShowsByLocationResponse{StatusCode: http.StatusInternalServerError}, err
		}
		i++
	}

	// showQuery, err := showService.DB.Table("addresses").
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

	showQuery, err := showService.DB.Table("addresses").
		Where("addresses.CityId = ?", user.CityId).
		Where("addresses.IsDeprecated = ?", false).
		Where("addresses.AddressType = ?", enums.Cinema).
		Joins("join cinemas on addresses.EntityId = cinemas.Id").
		Where("cinemas.IsDeprecated = ?", false).
		Joins("join cinemaHalls on cinemas.Id = cinemaHalls.CinemaId").
		Where("cinemaHalls.IsDeprecated = ?", false).
		Joins("join shows on cinemaHalls.Id = shows.CinemaHallId").
		Where("shows.IsDeprecated = ?", false).
		Where("shows.IsCancelled = ?", false).
		Joins("join movies on shows.MovieId = movies.Id").
		Where("movies.IsDeprecated = ?", false).
		Select("shows.Id AS ShowId", "shows.Date", "shows.StartTime", "shows.EndTime", "movies.Id AS MovieId", "movies.Title", "movies.Description", "movies.Language", "movies.Genre", "addresses.AddressLine").
		Select("ST_Distance_Sphere(addresses.Coordinates,?) AS Distance", user.Coordinates).
		Rows()
	if err != nil {
		return GetShowsByLocationResponse{StatusCode: http.StatusInternalServerError}, err
	}

	defer showQuery.Close()

	shows := []ShowsDTO{}
	for showQuery.Next() {
		show := &ShowsDTO{}
		err = showQuery.Scan(&show.ShowId, &show.Date, &show.startTime, &show.endTime, &show.MovieId, &show.Title,
			&show.Description, &show.Language, &show.Genre, &show.AddressLine, &show.Distance)

		if err != nil {
			return GetShowsByLocationResponse{StatusCode: http.StatusInternalServerError}, err
		}

		// show.StartTime = time.Unix(show.startTime, 0)
		// show.EndTime = time.Unix(show.endTime, 0)

		// show.Distance = distance(user.Coordinates, show.Coordinates)

		shows = append(shows, *show)
	}

	return GetShowsByLocationResponse{Results: shows, StatusCode: http.StatusOK}, nil
}
