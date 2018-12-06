package httpapi

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

//NewRouter returns an HTTP router for the HTTP API
func NewRouter(w io.Writer, db *sql.DB) http.Handler {

	//construct middleware
	var m = func(h returnHandler) http.Handler {
		return logMiddleware(jsonMiddleware(txMiddleware(h, db)), w)
	}

	r := mux.NewRouter()

	r.Path("/devices/{serial}").Methods("GET").Handler(m(handleReadDevice()))

	r.NotFoundHandler = m(notFoundHandler)

	return http.StripPrefix("/api/2.0", r)
}
