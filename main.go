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

	r := mux.NewRouter()
	r.HandleFunc("/", Welcome)

	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
