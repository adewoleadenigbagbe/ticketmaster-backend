package services

import (
	"errors"
	"fmt"
	"net/http"
)

type GetSearchRequest struct {
	Term string `query:"term"`
}

type GetSearchResponse struct {
	StatusCode int
	Result     []MovieDataResponse `json:"results"`
}

func (movieService MovieService) SearchMovie(request GetSearchRequest) (GetSearchResponse, error) {
	if len(request.Term) == 0 {
		return GetSearchResponse{StatusCode: http.StatusBadRequest}, errors.New("enter a search term")
	}

	var movieResult []MovieDataResponse
	sqlQuery := fmt.Sprintf("SELECT Id,Title,Description,ReleaseDate,Genre,Popularity,VoteCount FROM movies WHERE MATCH (movies.Title,movies.Description) AGAINST ('%s')", request.Term)
	dbResult := movieService.DB.Raw(sqlQuery).Scan(&movieResult)

	if dbResult.Error != nil {
		return GetSearchResponse{StatusCode: http.StatusInternalServerError}, errors.New(dbResult.Error.Error())
	}
	return GetSearchResponse{Result: movieResult, StatusCode: http.StatusOK}, nil
}
