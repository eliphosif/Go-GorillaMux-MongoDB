package main

import (  
	"fmt"
	"log"
	"net/http" 
	"github.com/gorilla/mux"  
)
func initlizeRouter() {

	r := mux.NewRouter()

	r.HandleFunc("/api/login", agentLogin).Methods("GET")
	r.HandleFunc("/api/customers", getCustomers).Methods("GET")
	r.HandleFunc("/api/customer/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/api/customers", createCustomer).Methods("POST")
	r.HandleFunc("/api/customer/{id}", updateCustomers).Methods("PUT")
	r.HandleFunc("/api/customer/{id}", deleteCustomers).Methods("DELETE")

	fmt.Println("server is listening")
	log.Fatal(http.ListenAndServe(":8090", r))

}