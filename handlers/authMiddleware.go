package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func (h *Handler) auth(next http.HandlerFunc) http.HandlerFunc {
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
