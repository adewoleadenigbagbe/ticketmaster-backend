package models

type ErrorResponse struct {
	Errors     []error
	StatusCode int
}
