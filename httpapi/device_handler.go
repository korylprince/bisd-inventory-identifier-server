package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-inventory-identifier-server/api"
)

func readDevice(r *http.Request, tx *sql.Tx) (int, interface{}) {
	serial := mux.Vars(r)["id"]

	device, err := api.ReadDeviceBySerialNumber(r.Context(), tx, serial)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to query device: %v", err)
	}

	if device == nil {
		return http.StatusNotFound, errors.New("Could not find device")
	}

	return http.StatusOK, device
}
