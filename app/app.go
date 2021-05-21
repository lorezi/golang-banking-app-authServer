package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

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

	// dbClient := db.Connect()

	// authRepo := repositories.NewAuthRepositoryDb(dbClient)

	// wiring

	// auth := handlers.AuthHandler{
	// 	Service: service.NewAuthService(authRepo),
	// }

	// defining routes

	// router.HandleFunc("/customers", auth.).Methods("GET")

	// starting serve
	addr := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), router))
}
