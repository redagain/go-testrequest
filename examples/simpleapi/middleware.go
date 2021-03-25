package simpleapi

import "net/http"

type bearerAuthMiddleware struct {
	introspectToken func(accessToken string) (active bool, err error)
}

func (h *bearerAuthMiddleware) HandleRequest(w http.ResponseWriter, req *http.Request, next RequestHandler) error {
	token, ok := BearerAuth(req)
	if !ok {
		return ErrUnauthorized
	}
	active, err := h.introspectToken(token)
	if err != nil {
		return err
	}
	if !active {
		return ErrUnauthorized
	}
	return next(w, req)
}
