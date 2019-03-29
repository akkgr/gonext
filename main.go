package main

import (
	"log"
	"net/http"
	"os"

	"gonext/server"
	"gonext/handlers"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logger.Println("server starting")

	mux := http.NewServeMux()
	h := handlers.NewHandler(logger)
	h.SetupRoutes(mux)

	srv := server.New(mux, ":8080")
	// openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650
	err := srv.ListenAndServeTLS("./certs/tls.crt", "./certs/tls.key")
	if err != nil {
		logger.Fatal(err)
	}
}
