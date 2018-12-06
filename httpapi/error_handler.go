package httpapi

import (
	"errors"
	"net/http"
)

//ErrorResponse represents an HTTP error
type ErrorResponse struct {
	Code        int    `json:"code"`
	Error       string `json:"error"`
	Description string `json:"-"`
}

//handleError returns a handlerResponse response for the given code
func handleError(code int, err error) *handlerResponse {
	return &handlerResponse{Code: code, Body: &ErrorResponse{Code: code, Error: http.StatusText(code)}}
}

//notFoundHandler returns a 401 handlerResponse
func notFoundHandler(w http.ResponseWriter, r *http.Request) *handlerResponse {
	return handleError(http.StatusNotFound, errors.New("Could not find handler"))
}
