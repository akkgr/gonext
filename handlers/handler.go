package handlers

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"encoding/json"
)

var secretKey = []byte("TooSlowTooLate4u.")

// Handler ...
type Handler struct {
	logger *log.Logger
	client *mongo.Client
	claims map[string]interface{}
}

// SetupRoutes ...
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.log(h.root))
	mux.HandleFunc("/token", h.log(h.getToken))
	mux.HandleFunc("/profile", h.log(h.auth(h.getProfile)))
}

// NewHandler ...
func NewHandler(logger *log.Logger, client *mongo.Client) *Handler {
	return &Handler{
		logger: logger,
		client: client,
	}
}

func returnJSON(status int, data interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.WriteHeader(status)
	w.Write(js)
}

func returnText(status int, data string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.WriteHeader(status)
	w.Write([]byte(data))
}
