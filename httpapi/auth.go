package httpapi

import (
	"errors"
	"net/http"
	"strings"
)

func authenticateRequest(r *http.Request, secret string) (int, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return http.StatusUnauthorized, errors.New("No Authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return http.StatusBadRequest, errors.New("Invalid Authorization header")
	}

	if auth != "Bearer "+secret {
		return http.StatusUnauthorized, errors.New("Invalid token")
	}

	return http.StatusOK, nil
}
