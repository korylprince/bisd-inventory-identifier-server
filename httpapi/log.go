package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Debug will pass debugging information to the client if true
var Debug = false

type logData struct {
	Action   string    `json:"action"`
	ActionID string    `json:"action_id,omitempty"`
	Error    string    `json:"error,omitempty"`
	Time     time.Time `json:"time"`
}

func logRequest(output io.Writer, action string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := &logData{Action: action}

		if id, ok := mux.Vars(r)["id"]; ok {
			l.ActionID = id
		}

		ctx := context.WithValue(r.Context(), contextKeyLogData, l)
		next.ServeHTTP(w, r.WithContext(ctx))

		l.Time = time.Now()
		j, err := json.Marshal(l)
		if err != nil {
			log.Println("Unable to marshal JSON:", err)
		}
		_, err = fmt.Fprintln(output, string(j))
		if err != nil {
			log.Println("Unable to output log:", err)
		}
	})
}
