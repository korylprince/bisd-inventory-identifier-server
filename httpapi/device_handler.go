package httpapi

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-inventory-identifier-server/api"
)

//GET /devices/:id
func handleReadDevice(chromeSvc *api.ChromebookService) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		id := mux.Vars(r)["id"]

		device, err := api.ReadDeviceByGoogleID(r.Context(), chromeSvc, id)

		if resp := checkAPIError(err); resp != nil {
			return resp
		}
		if device == nil {
			return handleError(http.StatusNotFound, errors.New("Could not find device"))
		}

		return &handlerResponse{Code: http.StatusOK, Body: device}
	}
}
