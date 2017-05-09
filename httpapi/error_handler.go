package httpapi

import (
	"errors"
	"net/http"

	"github.com/korylprince/bisd-inventory-identifier-server/api"
	"google.golang.org/api/googleapi"
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

//checkAPIError checks an api.Error and returns a handlerResponse for it, or nil if there was no error
func checkAPIError(err error) *handlerResponse {
	if err == nil {
		return nil
	}

	if e, ok := err.(*api.Error); ok {
		if ge, ok := (e.Err).(*googleapi.Error); ok && ge.Code == 404 {
			return handleError(http.StatusNotFound, err)
		}
	}

	return handleError(http.StatusInternalServerError, err)
}
