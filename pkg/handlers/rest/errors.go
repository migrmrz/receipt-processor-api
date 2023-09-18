package rest

import "net/http"

// ErrBadRequest represent a bad request error
type ErrBadRequest struct {
	message string
}

// Error implements the error interface
func (e ErrBadRequest) Error() string {
	return e.message
}

// StatusCode implements the StatusCode interface
func (e ErrBadRequest) StatusCode() int {
	return http.StatusBadRequest
}
