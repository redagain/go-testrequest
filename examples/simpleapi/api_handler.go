package simpleapi

import "net/http"

type RequestHandler func(w http.ResponseWriter, req *http.Request) error
