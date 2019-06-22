package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// Handler ...
type Handler struct {
	logger *log.Logger
	client *mongo.Client
	claims map[string]interface{}
}

// New ...
func New(logger *log.Logger, client *mongo.Client) *Handler {
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
