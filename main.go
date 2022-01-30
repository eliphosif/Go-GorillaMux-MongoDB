package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to my app deployed on Heroku")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}

	r := mux.NewRouter()
	r.HandleFunc("/welcome", Welcome)

	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
