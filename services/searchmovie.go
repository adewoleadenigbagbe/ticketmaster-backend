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

func (movieService MovieService) SearchMovie(request GetSearchRequest) (GetSearchResponse, models.ErrrorResponse) {
	if len(request.Term) == 0 {
		return GetSearchResponse{}, models.ErrrorResponse{Errors: []error{errors.New("enter a search term")}, StatusCode: http.StatusBadRequest}
	}

	var movieResult []MovieDataResponse
	sqlQuery := fmt.Sprintf("SELECT Id,Title,Description,ReleaseDate,Genre,Popularity,VoteCount FROM movies WHERE MATCH (movies.Title,movies.Description) AGAINST ('%s')", request.Term)
	dbResult := movieService.DB.Raw(sqlQuery).Scan(&movieResult)

	if dbResult.Error != nil {
		return GetSearchResponse{}, models.ErrrorResponse{Errors: []error{dbResult.Error}, StatusCode: http.StatusInternalServerError}
	}
	return GetSearchResponse{Result: movieResult}, models.ErrrorResponse{}
}
