package services

import (
	"errors"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
)

type GetMovieByIdRequest struct {
	Id string `param:"id" validate:"required"`
}

type GetMovieByIdResponse struct {
	StatusCode int
	Movie      MovieDataResponse `json:"movie"`
}

func (movieService MovieService) GetMovieById(request GetMovieByIdRequest) (GetMovieByIdResponse, error) {
	movie := &entities.Movie{
		Id: request.Id,
	}

	result := movieService.DB.First(movie)
	if result.Error != nil {
		return GetMovieByIdResponse{StatusCode: http.StatusNotFound}, errors.New("Movie Record not found")
	}

	movieDataResp := MovieDataResponse{
		Id:          movie.Id,
		Title:       movie.Title,
		Language:    movie.Language,
		Description: movie.Description.String,
		ReleaseDate: movie.ReleaseDate,
		Genre:       movie.Genre,
		Popularity:  movie.Popularity,
		VoteCount:   movie.VoteCount,
	}

	return GetMovieByIdResponse{StatusCode: http.StatusOK, Movie: movieDataResp}, nil
}
