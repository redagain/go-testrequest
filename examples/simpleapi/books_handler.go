package simpleapi

import (
	"net/http"
)

type BooksHandler struct {
	addBook func(req *BookRequest) (id string, err error)
}

func (h *BooksHandler) Post(w http.ResponseWriter, r *http.Request) error {
	req, err := NewBookRequest(r)
	if err != nil {
		return err
	}
	valid := req.IsValid()
	if !valid {
		return ErrBadRequest
	}
	id, err := h.addBook(req)
	if err != nil {
		return err
	}
	return Created(w, map[string]interface{}{"id": id})
}
