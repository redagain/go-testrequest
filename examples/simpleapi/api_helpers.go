package simpleapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func BearerAuth(req *http.Request) (accessToken string, ok bool) {
	auth := req.Header.Get("Authorization")
	if auth != "" {
		parts := strings.Fields(auth)
		prefix := parts[0]
		if len(parts) == 2 && strings.EqualFold(prefix, "bearer") {
			accessToken = parts[1]
			ok = true
		}
	}
	return
}

func Created(w http.ResponseWriter, resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return nil
}

func Handler(handler RequestHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			// ...
			// Write error to log
			// ...
			statusCode := ErrorStatusCode(err)
			w.WriteHeader(statusCode)
			return
		}
	})
}
