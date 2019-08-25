package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/akkgr/gonext/app"
	"github.com/akkgr/gonext/auth"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) getToken(w http.ResponseWriter, r *http.Request) {
	var lgnUser auth.User
	json.NewDecoder(r.Body).Decode(&lgnUser)

	collection := h.Client.Database("test").Collection("users")
	filter := bson.D{primitive.E{Key: "username", Value: lgnUser.Username}}
	var dbUser auth.User
	err := collection.FindOne(r.Context(), filter).Decode(&dbUser)
	if err != nil {
		h.Logger.Printf("%v", err)
		http.Error(w, "Inavlid username or password.", http.StatusUnauthorized)
		return
	}
	check := checkPasswordHash(lgnUser.Password, dbUser.Password)
	if check == false {
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
	ss, err := token.SignedString(app.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnJSON(http.StatusOK, ss, w)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
