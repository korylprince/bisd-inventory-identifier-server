package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonResponse(code int, data interface{}) http.Handler {
	type response struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		Debug       string `json:"debug,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err, ok := data.(error); ok || data == nil {
			resp := response{Code: code, Description: http.StatusText(code)}
			if err != nil {
				(r.Context().Value(contextKeyLogData)).(*logData).Error = err.Error()
				if Debug {
					resp.Debug = err.Error()
				}
			}
			data = resp
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)

		e := json.NewEncoder(w)
		eErr := e.Encode(data)

		if eErr != nil {
			log.Println("Error writing JSON response:", eErr)
		}
	})
}
