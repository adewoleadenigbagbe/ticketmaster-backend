package services

import (
	"errors"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/shared/database/entities"
	"github.com/Wolechacho/ticketmaster-backend/shared/models"
	"gorm.io/gorm"
)

type GetMovieByIdRequest struct {
	Id string `param:"id" validate:"required"`
}

type GetMovieByIdResponse struct {
	Movie MovieDataResponse `json:"movie"`
}

func (movieService MovieService) GetMovieById(request GetMovieByIdRequest) (GetMovieByIdResponse, models.ErrorResponse) {
	movie := &entities.Movie{}
	result := movieService.DB.Where("Id = ? AND IsDeprecated = ?", request.Id, false).First(&movie)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return GetMovieByIdResponse{}, models.ErrorResponse{StatusCode: http.StatusNotFound, Errors: []error{errors.New("movie Record not found")}}
	}

	movieDataResp := MovieDataResponse{
		Id:          movie.Id,
		Title:       movie.Title,
		Language:    movie.Language,
		Description: movie.Description,
		ReleaseDate: movie.ReleaseDate,
		Genre:       movie.Genre,
		Popularity:  movie.Popularity,
		VoteCount:   movie.VoteCount,
		Duration:    movie.Duration,
	}

	return GetMovieByIdResponse{Movie: movieDataResp}, models.ErrorResponse{}
}
