package simpleapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

type BookRequest struct {
	Title   string   `json:"title"`
	ISBN    string   `json:"isbn"`
	Authors []string `json:"authors"`
}

func (req *BookRequest) IsValid() (ok bool) {
	if req.Title == "" || req.ISBN == "" || len(req.Authors) == 0 {
		return
	}
	for _, v := range req.Authors {
		if v == "" {
			return
		}
	}
	ok = true
	return
}

func NewBookRequest(req *http.Request) (*BookRequest, error) {
	ct := req.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		return &BookRequest{
			Title:   req.PostFormValue("title"),
			ISBN:    req.PostFormValue("isbn"),
			Authors: req.PostForm["authors"],
		}, nil
	}
	if strings.HasPrefix(ct, "application/json") {
		b := &BookRequest{}
		err := json.NewDecoder(req.Body).Decode(b)
		if err != nil {
			return nil, ErrBadRequest
		}
		return b, nil
	}
	return nil, ErrUnsupportedMediaType
}
