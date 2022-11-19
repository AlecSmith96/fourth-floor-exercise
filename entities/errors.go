package entities

import "net/http"

// ResponseError struct for returning error from Twitch API
type ResponseError struct {
	error
	Code int
	PresentableError string
}

func NewNotFoundError() *ResponseError {
	return &ResponseError{
		Code: http.StatusNotFound,
		PresentableError: "a not found errror occured",
	}
}

func NewUnauthorizedError() *ResponseError {
	return &ResponseError{
		Code: http.StatusUnauthorized,
		PresentableError: "an unauthorized error occurred",
	}
}

func NewBadRequestError() *ResponseError {
	return &ResponseError{
		Code: http.StatusBadRequest,
		PresentableError: "a bad request error occurred",
	}
}
