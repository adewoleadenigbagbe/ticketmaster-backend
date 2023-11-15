package models

type ErrrorResponse struct {
	Errors     []error
	StatusCode int
}