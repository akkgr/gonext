package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/akkgr/gonext/handlers"
	"github.com/akkgr/gonext/server"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logger.Println("server starting")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	logger.Println("Connected to MongoDB!")
	dbInit(logger, client)

	mux := http.NewServeMux()
	h := handlers.NewHandler(logger, client)
	h.SetupRoutes(mux)

	srv := server.New(mux, ":8080")
	// openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650
	err = srv.ListenAndServeTLS("./certs/tls.crt", "./certs/tls.key")
	if err != nil {
		logger.Fatal(err)
	}
}

func dbInit(logger *log.Logger, client *mongo.Client) {
	collection := client.Database("test").Collection("users")
	filter := bson.D{{"username", "admin"}}

	res, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		logger.Fatal(err)
	}
	if res == 0 {
		pass, _ := hashPassword("Abc.123")
		_, err = collection.InsertOne(context.TODO(), bson.M{"username": "admin", "password": pass})
		if err != nil {
			logger.Fatal(err)
		}
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
