package httpapi

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-inventory-identifier-server/api"
)

//NewRouter returns an HTTP router for the HTTP API
func NewRouter(w io.Writer, chromeSvc *api.ChromebookService, db *sql.DB) http.Handler {

	//construct middleware
	var m = func(h returnHandler) http.Handler {
		return logMiddleware(jsonMiddleware(txMiddleware(h, db)), w)
	}

	r := mux.NewRouter()

	r.Path("/devices/{id}").Methods("GET").Handler(m(handleReadDevice(chromeSvc)))

	r.NotFoundHandler = m(notFoundHandler)

	return http.StripPrefix("/api/1.0", r)
}
