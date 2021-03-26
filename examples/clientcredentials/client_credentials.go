package clientcredentials

import (
	"net/http"
	"strings"
)

func BasicCredentials(req *http.Request) (clientID, secret string, ok bool) {
	clientID, secret, ok = req.BasicAuth()
	if ok {
		return
	}
	if req.Method == http.MethodPost && strings.HasPrefix(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		err := req.ParseForm()
		if err != nil {
			return
		}
		if _, sOk := req.PostForm["client_secret"]; sOk {
			clientID = req.FormValue("client_id")
			secret = req.PostFormValue("client_secret")
			ok = true
		}
	}
	return
}
