package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Counter struct {
	count int64
}

func (i *Counter) increment() int64 {
	i.count = i.count + 1
	return i.count
}

var counter Counter = Counter{0}

type Address struct {
	Hno    int    `json:"hno"`
	Street string `json:"street"`
	State  string `json:"state"`
}

type Customer struct {
	//	gorm.Model
	CustId            int64   `json:"custid"`
	FirstName         string  `json:"firstname"`
	LastName          string  `json:"lastname"`
	Email             string  `json:"email"`
	Age               int     `json:"age"`
	Address           Address `json:"address"`
	CreatedbyAgentID  string  `json:"createdbyagentid"`
	ModifiedbyAgentID string  `json:"modifiedbyagentid"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	json.NewEncoder(w).Encode(Customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	params := mux.Vars(r) // params
	id, _ := strconv.ParseInt(params["id"], 10, 64)
	for _, cust := range Customers {
		if cust.CustId == id {
			json.NewEncoder(w).Encode(cust)
			return
		}
	}
	json.NewEncoder(w).Encode(&Customer{})
}

func updateCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	params := mux.Vars(r) // params

	id, _ := strconv.ParseInt(params["id"], 10, 64)
	for index, item := range Customers {
		if item.CustId == id {
			Customers = append(Customers[:index], Customers[index+1:]...)
			var cust Customer
			_ = json.NewDecoder(r.Body).Decode(&cust)
			cust.CustId = id
			cust.ModifiedbyAgentID = loggedinAgent
			Customers = append(Customers, cust)

			json.NewEncoder(w).Encode(cust)
			return
		}
	}
	json.NewEncoder(w).Encode(Customers)
}

func deleteCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	params := mux.Vars(r) // params
	id, _ := strconv.ParseInt(params["id"], 10, 64)
	for index, cust := range Customers {
		if cust.CustId == id {
			Customers = append(Customers[:index], Customers[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(Customers)
}


//createCustomer
func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	var cust Customer
	_ = json.NewDecoder(r.Body).Decode(&cust)
	cust.CreatedbyAgentID = loggedinAgent
	cust.CustId = counter.increment()
	//Customers = append(Customers, cust)

	res, err:=insertCustomerDoc(AllCustomers, cust, ctx)
	if err != nil {
		return
	}
	fmt.Println(res)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(cust)
}

