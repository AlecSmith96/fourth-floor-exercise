package entities

import "net/http"

// ResponseError error struct for returning error from Twitch API
type ResponseError struct {
	error
	Code             int
	PresentableError string
}

// NewNotFoundError returns an error equivalent to a 404 Not Found HTTP status
func NewNotFoundError() *ResponseError {
	return &ResponseError{
		Code:             http.StatusNotFound,
		PresentableError: "a not found errror occured",
	}
}

// NewUnauthorizedError returns an error equivalent to a 401 Unauthorized HTTP status
func NewUnauthorizedError() *ResponseError {
	return &ResponseError{
		Code:             http.StatusUnauthorized,
		PresentableError: "an unauthorized error occurred",
	}
}

// NewBadRequestError returns an error equivalent to a 400 Bad Request HTTP status
func NewBadRequestError() *ResponseError {
	return &ResponseError{
		Code:             http.StatusBadRequest,
		PresentableError: "a bad request error occurred",
	}
}
