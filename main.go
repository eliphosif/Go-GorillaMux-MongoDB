package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	initlizeRouter()
}

func initlizeRouter() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	r := mux.NewRouter()

	r.HandleFunc("/welcome", Welcome).Methods("GET")

	fmt.Println("server is listening:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to my app deployed on Heroku")
}
