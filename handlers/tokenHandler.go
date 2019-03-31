package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) getToken(w http.ResponseWriter, r *http.Request) {
	var lgnUser user
	json.NewDecoder(r.Body).Decode(&lgnUser)

	collection := h.client.Database("test").Collection("users")
	filter := bson.D{{"username", lgnUser.Username}}
	var dbUser user
	err := collection.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		h.logger.Printf("%v", err)
	}

	exp := time.Now().Add(time.Hour * 8).Unix()
	claims := &jwt.StandardClaims{
		ExpiresAt: exp,
		Issuer:    "test",
		Subject:   lgnUser.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnText(http.StatusOK, ss, w)
}
