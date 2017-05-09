package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/korylprince/bisd-inventory-identifier-server/api"
)

type handlerResponse struct {
	Code int
	Body interface{}
	Err  error
}

type returnHandler func(http.ResponseWriter, *http.Request) *handlerResponse

const logTemplate = "{{.Date}} {{.Method}} {{.Path}}{{if .Query}}?{{.Query}}{{end}} {{.Code}} ({{.Status}}){{if .Err}}, Error: {{.Err}}{{end}}\n"

type logData struct {
	Date   string
	Status string
	Code   int
	Method string
	Path   string
	Query  string
	Err    error
}

func logMiddleware(next returnHandler, writer io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := next(w, r)

		err := template.Must(template.New("log").Parse(logTemplate)).Execute(writer, &logData{
			Date:   time.Now().Format("2006-01-02:15:04:05 -0700"),
			Status: http.StatusText(resp.Code),
			Code:   resp.Code,
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Err:    resp.Err,
		})

		if err != nil {
			panic(err)
		}
	})
}

func jsonMiddleware(next returnHandler) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		var resp *handlerResponse

		if r.Method != "GET" {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				resp = handleError(http.StatusBadRequest, errors.New("Could not parse Content-Type"))
				goto serve
			}
			if mediaType != "application/json" {
				resp = handleError(http.StatusBadRequest, errors.New("Content-Type not application/json"))
				goto serve
			}
		}

		w.Header().Set("Content-Type", "application/json")
		resp = next(w, r)

	serve:
		w.WriteHeader(resp.Code)
		e := json.NewEncoder(w)
		err := e.Encode(resp.Body)
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could encode json: %v", err))
		}
		return resp
	}
}

func txMiddleware(next returnHandler, db *sql.DB) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		tx, err := db.Begin()
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not begin transaction: %v", err))
		}

		ctx := context.WithValue(r.Context(), api.TransactionKey, tx)
		resp := next(w, r.WithContext(ctx))

		if err = tx.Commit(); err != nil {
			if rErr := tx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
				return handleError(http.StatusInternalServerError, fmt.Errorf("Could not rollback transaction: %v", rErr))
			}
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not commit transaction: %v", err))
		}

		return resp
	}
}
