package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Wolechacho/ticketmaster-backend/models"
)

type GetSearchRequest struct {
	Term string `query:"term"`
}

type GetSearchResponse struct {
	Result []MovieDataResponse `json:"results"`
}

func (movieService MovieService) SearchMovie(request GetSearchRequest) (GetSearchResponse, models.ErrorResponse) {
	movieService.Logger.Info().Interface("getSearchRequest", request).Msg("request")
	if len(request.Term) == 0 {
		return GetSearchResponse{}, models.ErrorResponse{Errors: []error{errors.New("enter a search term")}, StatusCode: http.StatusBadRequest}
	}

	var movieResult []MovieDataResponse
	sqlQuery := fmt.Sprintf("SELECT Id,Title,Description,ReleaseDate,Genre,Popularity,VoteCount FROM movies WHERE MATCH (movies.Title,movies.Description) AGAINST ('%s')", request.Term)
	dbResult := movieService.DB.Raw(sqlQuery).Scan(&movieResult)

	if dbResult.Error != nil {
		movieService.Logger.Info().Interface("getSearchResponse", dbResult.Error.Error()).Msg("response")
		return GetSearchResponse{}, models.ErrorResponse{Errors: []error{dbResult.Error}, StatusCode: http.StatusInternalServerError}
	}

	resp := GetSearchResponse{Result: movieResult}
	movieService.Logger.Info().Interface("getSearchResponse", resp).Msg("response")
	return resp, models.ErrorResponse{}
}
