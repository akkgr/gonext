package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/akkgr/gonext/handlers"
	"github.com/akkgr/gonext/server"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
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

	// go http.ListenAndServe(":8080", http.HandlerFunc(redirect))

	mux := mux.NewRouter()
	h := handlers.NewHandler(logger, client)
	h.SetupRoutes(mux)

	srv := server.New(mux, ":8080", logger)
	// openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650
	go func() {
		err = srv.ListenAndServeTLS("./certs/server.rsa.crt", "./certs/server.rsa.key")
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
		logger.Println("Server stopped")
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	var wait time.Duration
	wait = time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	client.Disconnect(ctx)
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func dbInit(logger *log.Logger, client *mongo.Client) {
	collection := client.Database("test").Collection("users")
	filter := bson.D{{Key: "username", Value: "admin"}}

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

// func redirect(w http.ResponseWriter, req *http.Request) {
// 	http.Redirect(w, req,
// 		"https://"+req.Host+req.URL.String(),
// 		http.StatusMovedPermanently)
// }
