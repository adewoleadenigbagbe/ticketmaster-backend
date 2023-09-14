package controllers

import (
	"fmt"
	"net/http"
	"time"

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

	//Filter
	filterClause := paginate.FilterFields(&entities.Movie{IsDeprecated: false})

	//paginate
	movies := new([]entities.Movie)
	paginateClause := paginate.Paginate(req.Page, req.PageLength)

	//order by
	sortandorder := fmt.Sprintf("%s %s", req.SortBy, req.Order)
	fmt.Println(sortandorder)
	orderByClause := paginate.OrderBy(sortandorder)

	//this uses functional scope pattern in golang
	db.DB.Scopes(filterClause, paginateClause, orderByClause).Find(&movies)

	var countResult int64
	paginate.GetEntityCount(db.DB, new(entities.Movie), &countResult)

	//you can pass in the deprecated to this function
	fmt.Println("count result", countResult)

	resp := new(getMoviesResponse)
	resp.Page = req.Page
	resp.RequestedPageLength = req.PageLength
	resp.PerPage = len(*movies)
	resp.TotalResults = countResult

	for _, movie := range *movies {
		movieData := MovieDataResponse{
			Id:          movie.Id,
			Title:       movie.Title,
			Language:    movie.Language,
			Description: movie.Description.String,
			ReleaseDate: movie.ReleaseDate,
			Genre:       movie.Genre,
			Popularity:  movie.Popularity,
			VoteCount:   movie.VoteCount,
		}
		resp.Movies = append(resp.Movies, movieData)
	}
	return movieContext.JSON(http.StatusOK, resp)
}

type getMoviesRequest struct {
	Page       int    `query:"page"`
	PageLength int    `query:"pageLength"`
	SortBy     string `query:"sortBy"`
	Order      string `query:"order"`
}

type getMoviesResponse struct {
	Page                int                 `json:"page"`
	PerPage             int                 `json:"perPage"`
	TotalResults        int64               `json:"totalResults"`
	RequestedPageLength int                 `json:"requestedPageLength"`
	Movies              []MovieDataResponse `json:"movies"`
}

type MovieDataResponse struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	ReleaseDate time.Time `json:"releaseDate"`
	Genre       int       `json:"genre"`
	Popularity  float32   `json:"popularity"`
	VoteCount   int       `json:"voteCount"`
}

func (movieController *MovieController) GetMovieById(movieContext echo.Context) error {
	req := new(getMovieByIdRequest)

	err := movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	movie := &entities.Movie{
		Id: req.Id,
	}
	result := db.DB.First(movie)

	if result.Error != nil {
		return movieContext.JSON(http.StatusNotFound, "Movie Record not found")
	}

	resp := new(getMovieByIdResponse)
	resp.Movie = MovieDataResponse{
		Id:          movie.Id,
		Title:       movie.Title,
		Language:    movie.Language,
		Description: movie.Description.String,
		ReleaseDate: movie.ReleaseDate,
		Genre:       movie.Genre,
		Popularity:  movie.Popularity,
		VoteCount:   movie.VoteCount,
	}

	return movieContext.JSON(http.StatusOK, resp)
}

type getMovieByIdRequest struct {
	Id string `param:"id" validate:"required"`
}

type getMovieByIdResponse struct {
	Movie MovieDataResponse `json:"movie"`
}

func (movieController MovieController) SearchMovie(movieContext echo.Context) error {
	var err error
	req := new(getSearchRequest)

	err = movieContext.Bind(req)
	if err != nil {
		return movieContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	if len(req.Term) == 0 {
		return movieContext.JSON(http.StatusBadRequest, "enter a search term")
	}

	var movieResult []MovieDataResponse
	sqlQuery := fmt.Sprintf("SELECT Id,Title,Description,ReleaseDate,Genre,Popularity,VoteCount FROM movies WHERE MATCH (movies.Title,movies.Description) AGAINST ('%s')", req.Term)

	dbResult := db.DB.Raw(sqlQuery).Scan(&movieResult)
	if dbResult.Error != nil {
		return movieContext.JSON(http.StatusInternalServerError, dbResult.Error.Error())
	}

	resp := new(getSearchResponse)
	resp.Result = movieResult
	return movieContext.JSON(http.StatusOK, resp)
}

type getSearchRequest struct {
	Term string `query:"term"`
}

type getSearchResponse struct {
	Result []MovieDataResponse
}
