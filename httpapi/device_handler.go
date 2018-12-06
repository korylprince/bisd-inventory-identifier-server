package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-inventory-identifier-server/api"
)

//GET /devices/:id
func handleReadDevice() returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		serial := mux.Vars(r)["serial"]

		device, err := api.ReadDeviceBySerialNumber(r.Context(), serial)

		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Unable to query device: %v", err))
		}

		if device == nil {
			return handleError(http.StatusNotFound, errors.New("Could not find device"))
		}

		return &handlerResponse{Code: http.StatusOK, Body: device}
	}
}
