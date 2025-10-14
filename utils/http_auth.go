package utils

import (
	"crypto/subtle"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	handler  http.Handler
	username string
	password string
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(h.username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(h.password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.handler.ServeHTTP(w, r)
}

func BasicAuth(handler http.Handler, username, password string) *AuthHandler {
	return &AuthHandler{
		handler:  handler,
		username: username,
		password: password,
	}
}

func (h *AuthHandler) String() string {
	return fmt.Sprintf("BasicAuth(username=%q, password=%q)", h.username, h.password)
}
