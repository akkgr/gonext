package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	mux.HandleFunc("/", h.Logger(h.Auth(h.handle)))
	mux.HandleFunc("/token", h.Logger(h.getToken))
}

// NewHandler ...
func NewHandler(logger *log.Logger, client *mongo.Client) *Handler {
	return &Handler{
		logger: logger,
		client: client,
	}
}

// Auth ...
func (h *Handler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authStr := r.Header.Get("Authorization")
		if !strings.HasPrefix(authStr, "Bearer ") {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		tokenString := authStr[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			h.claims = claims
			next(w, r)
		} else {
			http.Error(w, "invalid token", http.StatusUnauthorized)
		}
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

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	returnText(http.StatusOK, h.claims["sub"].(string), w)
}

func (h *Handler) getToken(w http.ResponseWriter, r *http.Request) {

	var user User
	json.NewDecoder(r.Body).Decode(&user)
	h.logger.Printf("%v", user)
	exp := time.Now().Add(time.Hour * 8).Unix()
	claims := &jwt.StandardClaims{
		ExpiresAt: exp,
		Issuer:    "test",
		Subject:   user.username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnText(http.StatusOK, ss, w)
}

type User struct {
	username string
}
