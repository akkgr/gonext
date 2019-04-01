package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/akkgr/gonext/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) getToken(w http.ResponseWriter, r *http.Request) {
	var lgnUser models.User
	json.NewDecoder(r.Body).Decode(&lgnUser)

	collection := h.client.Database("test").Collection("users")
	filter := bson.D{{"username", lgnUser.Username}}
	var dbUser models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil {
		h.logger.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	check := checkPasswordHash(lgnUser.Password, dbUser.Password)
	if check {
		http.Error(w, "Inavlid username or password.", http.StatusUnauthorized)
		return
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

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
