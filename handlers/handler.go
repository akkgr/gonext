package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler ...
type Handler struct {
	Logger *log.Logger
	Client *mongo.Client
	Claims map[string]interface{}
}

// NewHandler ...
func NewHandler(logger *log.Logger, client *mongo.Client) *Handler {
	return &Handler{
		Logger: logger,
		Client: client,
	}
}

// SetupRoutes ...
func (h *Handler) SetupRoutes(mux *mux.Router) {
	mux.HandleFunc("/", h.root)
	mux.HandleFunc("/token", h.getToken)
	mux.Use(h.log)

	api := mux.PathPrefix("/api").Subrouter()
	api.HandleFunc("/profile", h.getProfile)
	api.Use(h.auth)
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
