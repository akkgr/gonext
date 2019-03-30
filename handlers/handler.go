package handlers

import (
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Handler ...
type Handler struct {
	logger *log.Logger
	client *mongo.Client
}

// SetupRoutes ...
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Logger(h.handle))
}

// NewHandler ...
func NewHandler(logger *log.Logger, client *mongo.Client) *Handler {
	return &Handler{
		logger: logger,
		client: client,
	}
}

// Logger ...
func (h *Handler) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s\n", time.Now().Sub(startTime))
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next(w, r)
	}
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello kosta"))
}
