package controllers

import (
	"fmt"
	"net/http"

	db "github.com/Wolechacho/ticketmaster-backend/database"
	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	paginate "github.com/Wolechacho/ticketmaster-backend/helpers/pagination"
	"github.com/labstack/echo/v4"
)

type MovieController struct {
}

func (movieController *MovieController) GetMovies(movieContext echo.Context) error {
	req := new(getMoviesRequest)
	err := movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	switch {
	case req.PageLength > 100:
		req.PageLength = 100
	case req.PageLength <= 0:
		req.PageLength = 10
	}

	//this uses functional scope pattern in golang
	movies := new([]entities.Movie)
	paginateFunc := paginate.Paginate(req.Page, req.PageLength)
	db.DB.Scopes(paginateFunc).Find(&movies)

	var countResult int64
	paginate.GetEntityCount(db.DB, new(entities.Movie), &countResult)

	//you can pass in the deprecated to this function
	fmt.Println("count result", countResult)

	resp := new(getMoviesResponse)
	resp.Page = req.Page
	resp.RequestedPageLength = req.PageLength
	resp.PerPage = len(*movies)
	resp.TotalResults = countResult
	resp.Movies = *movies
	return movieContext.JSON(http.StatusOK, resp)
}

type getMoviesRequest struct {
	Page       int    `query:"page"`
	PageLength int    `query:"pageLength"`
	SortBy     string `query:"sortBy"`
	Order      string `query:"order"`
}

type getMoviesResponse struct {
	Page                int              `json:"page"`
	PerPage             int              `json:"perPage"`
	TotalResults        int64            `json:"totalResults"`
	RequestedPageLength int              `json:"requestedPageLength"`
	Movies              []entities.Movie `json:"movies"`
}

// type MovieDataResponse struct {
// 	Id           string `gorm:"primaryKey;size:36"`
// 	Title        string `gorm:"not null"`
// 	Description  sql.NullString
// 	Language     string    `gorm:"not null"`
// 	ReleaseDate  time.Time `gorm:"not null"`
// 	Duration     sql.NullInt32
// 	Genre        int     `gorm:"not null"`
// 	Popularity   float32 `gorm:"not null"`
// 	VoteCount    int     `gorm:"not null"`
// 	IsDeprecated bool
// }
