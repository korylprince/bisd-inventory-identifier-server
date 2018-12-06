package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//API is the current API version
const API = "2.0"
const apiPath = "/api/" + API

func notFound(w http.ResponseWriter, r *http.Request) {
	jsonResponse(http.StatusNotFound, nil).ServeHTTP(w, r)
}

type addRouteHelperFunc func(action, method, path string, f func(*http.Request, *sql.Tx) (int, interface{}))

func addRouteHelper(api *mux.Router, s *Server) addRouteHelperFunc {
	return func(action, method, path string, f func(*http.Request, *sql.Tx) (int, interface{})) {
		api.Methods(method).Path(path).Handler(
			//log request
			logRequest(s.output, action,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					//check auth
					if code, err := authenticateRequest(r, s.secret); err != nil {
						jsonResponse(code, err).ServeHTTP(w, r)
						return
					}

					//create db tx
					tx, err := s.db.Begin()
					if err != nil {
						jsonResponse(
							http.StatusInternalServerError,
							fmt.Errorf("Unable to start database transaction: %v", err),
						).ServeHTTP(w, r)
						return
					}

					//run handler
					code, body := f(r, tx)

					//check if not OK
					if code != http.StatusOK {
						if err = tx.Rollback(); err != nil {
							jsonResponse(
								http.StatusInternalServerError,
								fmt.Errorf("Unable to rollback database transaction: %v", err),
							).ServeHTTP(w, r)
						}
						jsonResponse(code, body).ServeHTTP(w, r)
						return
					}

					//close tx
					if err = tx.Commit(); err != nil {
						jsonResponse(
							http.StatusInternalServerError,
							fmt.Errorf("Unable to commit database transaction: %v", err),
						).ServeHTTP(w, r)
						return
					}

					//respond with handler body
					jsonResponse(code, body).ServeHTTP(w, r)
				}),
			))
	}
}

//Router returns a new router for the given Server
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix(apiPath).Subrouter()

	api.NotFoundHandler = http.HandlerFunc(notFound)

	addRoute := addRouteHelper(api, s)

	addRoute("ReadDevice", "GET", "/device/{id}", readDevice)

	return r
}
