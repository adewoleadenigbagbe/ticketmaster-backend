package services

import (
	"fmt"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	paginate "github.com/Wolechacho/ticketmaster-backend/helpers/pagination"
	"github.com/Wolechacho/ticketmaster-backend/helpers/utilities"
)

type GetMoviesRequest struct {
	Page       int    `query:"page" default:"1"`
	PageLength int    `query:"pageLength" default:"10"`
	SortBy     string `query:"sortBy" default:"Title"`
	Order      string `query:"order" default:"asc"`
}

type GetMoviesResponse struct {
	Page                int                 `json:"page"`
	PerPage             int                 `json:"perPage"`
	TotalResults        int64               `json:"totalResults"`
	RequestedPageLength int                 `json:"requestedPageLength"`
	Movies              []MovieDataResponse `json:"movies"`
}

type MovieDataResponse struct {
	Id          string                     `json:"id"`
	Title       string                     `json:"title"`
	Description utilities.Nullable[string] `json:"description"`
	Language    string                     `json:"language"`
	ReleaseDate time.Time                  `json:"releaseDate"`
	Genre       int                        `json:"genre"`
	Popularity  float32                    `json:"popularity"`
	VoteCount   int                        `json:"voteCount"`
	Duration    utilities.Nullable[int]    `json:"duration"`
}

func (movieService MovieService) GetMovies(request GetMoviesRequest) (GetMoviesResponse, error) {
	err := utilities.SetDefaults[GetMoviesRequest](&request)
	if err != nil {
		return GetMoviesResponse{}, err
	}

	//Filter
	filterClause := paginate.FilterFields(&entities.Movie{IsDeprecated: false})

	//paginate
	movies := new([]entities.Movie)
	paginateClause := paginate.Paginate(request.Page, request.PageLength)

	//order by
	sortandorder := fmt.Sprintf("%s %s", request.SortBy, request.Order)
	orderByClause := paginate.OrderBy(sortandorder)

	//this uses functional scope pattern in golang
	movieService.DB.Scopes(filterClause, orderByClause, paginateClause).Find(&movies)

	var countResult int64
	paginate.GetEntityCount(movieService.DB, new(entities.Movie), &countResult)

	resp := new(GetMoviesResponse)
	resp.Page = request.Page
	resp.RequestedPageLength = request.PageLength
	resp.PerPage = len(*movies)
	resp.TotalResults = countResult

	for _, movie := range *movies {
		movieData := MovieDataResponse{
			Id:          movie.Id,
			Title:       movie.Title,
			Language:    movie.Language,
			Description: movie.Description,
			ReleaseDate: movie.ReleaseDate,
			Genre:       movie.Genre,
			Popularity:  movie.Popularity,
			VoteCount:   movie.VoteCount,
		}
		resp.Movies = append(resp.Movies, movieData)
	}

	return *resp, nil
}
