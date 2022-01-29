package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func initlizeRouter() {

	r := mux.NewRouter()
	r.HandleFunc("/", welcome).Methods("GET")
	r.HandleFunc("/api/login", agentLogin).Methods("GET")
	r.HandleFunc("/api/customers", getCustomers).Methods("GET")
	r.HandleFunc("/api/customer/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/api/customers", createCustomer).Methods("POST")
	r.HandleFunc("/api/customer/{id}", updateCustomers).Methods("PUT")
	r.HandleFunc("/api/customer/{id}", deleteCustomers).Methods("DELETE")

	fmt.Println("server is listening")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "welcome to my app")
}
