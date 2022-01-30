package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome World")
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", Welcome)
	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
