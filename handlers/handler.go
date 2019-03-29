package handlers

import (
	"log"
	"net/http"
	"time"
)

// Handler ...
type Handler struct {
	logger *log.Logger
}

// SetupRoutes ...
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Logger(h.handle))
}

// NewHandler ...
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Logger ...
func (h *Handler) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello kosta"))
}
