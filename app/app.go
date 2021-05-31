package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lorezi/golang-bank-app-auth/db"
	"github.com/lorezi/golang-bank-app-auth/handlers"
	"github.com/lorezi/golang-bank-app-auth/repositories"
	"github.com/lorezi/golang-bank-app-auth/service"

	"github.com/subosito/gotenv"
)

func sanitizeConfigs() {

	if os.Getenv("SERVER_ADDRESS") == "" || os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Environment variable not defined...")
	}
}

func Start() {
	gotenv.Load()

	sanitizeConfigs()
	// created multiplexer
	router := mux.NewRouter()

	dbClient := db.Connect()

	authRepo := repositories.NewAuthRepositoryDb(dbClient)

	// // wiring

	auth := handlers.AuthHandler{
		Service: service.NewAuthService(authRepo),
	}

	// // defining routes

	router.HandleFunc("/auth/login", auth.Login).Methods("POST")
	// /verify
	router.HandleFunc("/auth/verify", auth.Verify).Methods("GET")
	// /refresh
	router.HandleFunc("/auth/verify", auth.Refresh).Methods(http.MethodPost)

	// starting serve
	addr := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Println(fmt.Sprintf("Starting OAuth server on %s:%s 🤝", addr, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), router))
}
