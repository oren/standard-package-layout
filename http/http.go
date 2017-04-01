package http

import (
	"net/http"

	"github.com/oren/standard-package-layout"
)

type Handler struct {
	UserService myapp.UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle request
}
