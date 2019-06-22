package app

import (
	"github.com/gorilla/mux"
)

var secretKey = []byte("TooSlowTooLate4u.")

// SetupRoutes ...
func (h *Handler) SetupRoutes(mux *mux.Router) {
	mux.HandleFunc("/", h.root)
	mux.HandleFunc("/token", h.getToken)
	mux.Use(h.log)

	api := mux.PathPrefix("/api").Subrouter()
	api.HandleFunc("/profile", h.getProfile)
	api.Use(h.auth)
}
