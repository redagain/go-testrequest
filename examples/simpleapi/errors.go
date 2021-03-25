package simpleapi

import (
	"errors"
	"net/http"
)

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrBadRequest           = errors.New("bad request")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)

func ErrorStatusCode(err error) (statusCode int) {
	if errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized
	} else if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrUnsupportedMediaType) {
		return http.StatusUnsupportedMediaType
	} else {
		return http.StatusInternalServerError
	}
}
