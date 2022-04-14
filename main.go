package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserLogin struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Jwt      string `json:"jwt"`
}

type UserLoginError struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Route       string `json:"route"`
}

var userLogin UserLogin
var userLoginError UserLoginError = UserLoginError{
	Message:     "unauthorized",
	Description: "it seems you're not logged in, Please login ang try again ! :)",
	Route: "https://go-gorillamux-mongodb-myapp01.herokuapp.com/api/login",
}

func isLoggedin(r *http.Request) bool {
	if jwt := r.Header.Get("jwt"); jwt == userLogin.Jwt && jwt != "" {
		return true
	}
	return false
}

func CreateToken(userid string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "KishoreJWT")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	userLogin.Jwt = token
	return token, nil
}

func agentLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&userLogin)
	var count int = 0
	for _, agent := range Agents {
		if userLogin.UserID == agent.AgentID {
			count = count + 1
			if userLogin.Password == agent.AgentPassword {
				userLogin.Username = agent.AgentName
				loggedinAgent = agent.AgentID
			} else {

				fmt.Fprintf(w, "Please provide valid login details")
				return
			}
		}
	}
	if count == 0 {
		fmt.Fprintf(w, "the Agent does not exist")
		return
	}

	//	userLogin.UserID = string(time.Now().Format("20060102150405"))
	token, err := CreateToken(userLogin.UserID)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	userLogin.Jwt = token
	json.NewEncoder(w).Encode(userLogin)
}

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
	//	similar to gorm.Model
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	custsfromDB := getAllDocuments()

	json.NewEncoder(w).Encode(custsfromDB)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	params := mux.Vars(r) // params

	idPrimitive, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	}
	cursor, finderr := AllCustomers.Find(ctx, bson.M{"_id": idPrimitive})
	if finderr != nil {
		log.Fatal(finderr)
	}

	var custSlice []bson.M
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var custfromDB bson.M
		if err = cursor.Decode(&custfromDB); err != nil {
			log.Fatal(err)
		}
		custSlice = append(custSlice, custfromDB)

	}
	json.NewEncoder(w).Encode(custSlice)
}

func updateOneCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	params := mux.Vars(r) // params

	idPrimitive, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	}

	var cust Customer
	_ = json.NewDecoder(r.Body).Decode(&cust)

	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", bson.D{
		{"firstname", cust.FirstName},
		{"lastname", cust.LastName},
		{"age", cust.Age},
		{"email", cust.Email},
		{"address.hno", cust.Address.Hno},
		{"address.state", cust.Address.State},
		{"address.street", cust.Address.Street},
	}}}
	result, err := AllCustomers.UpdateOne(ctx, bson.M{"_id": idPrimitive}, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
}

//delete one customer
func deleteOneCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !isLoggedin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	params := mux.Vars(r) // params

	idPrimitive, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	}
	result, err := AllCustomers.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNoContent)
		msg := "pls check the ID passed"
		io.WriteString(w, msg)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

//delete all customer
func deleteAllCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isLoggedin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}

	// BSON filter for all documents with a value of 1
	f := bson.M{}
	fmt.Println("nDeleting documents with filter:", f)

	// Call the DeleteMany() method to delete docs matching filter
	result, err := AllCustomers.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
	json.NewEncoder(w).Encode(result)
}

//createCustomer
func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isLoggedin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(userLoginError)
		return
	}
	var cust Customer
	_ = json.NewDecoder(r.Body).Decode(&cust)
	cust.CreatedbyAgentID = loggedinAgent
	cust.CustId = counter.increment()
	//Customers = append(Customers, cust)

	res, err := insertCustomerDoc(AllCustomers, cust, ctx)
	if err != nil {
		return
	}
	fmt.Println(res)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(cust)
} 
var Agents []SalesAgent

var loggedinAgent string
var AllCustomers *mongo.Collection
var ctx = context.Background()

func initlizeMongoConnection() *mongo.Collection {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	USERNAME := os.Getenv("USERNAME")
	PASSWORD := os.Getenv("PASSWORD")

	MongoDBURI := "mongodb+srv://" + USERNAME + ":" + PASSWORD + "@cluster0.0imgv.mongodb.net/test?authSource=admin&replicaSet=atlas-fydu9m-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true"

	//defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBURI))

	golangMongoDB := client.Database("GolangMongo")
	AllCustomers := golangMongoDB.Collection("AllCustomers")
	return AllCustomers

}

func insertCustomerDoc(AllCustomers *mongo.Collection, cust Customer, ctx context.Context) (*mongo.InsertOneResult, error) {
	//insert document in to MongoDB collection
	res, err := AllCustomers.InsertOne(ctx, cust)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}

func getAllDocuments() []bson.M {

	cursor, err := AllCustomers.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var custSlice []bson.M
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var custfromDB bson.M
		if err = cursor.Decode(&custfromDB); err != nil {
			log.Fatal(err)
		}
		custSlice = append(custSlice, custfromDB)

	}
	return custSlice
}

func main() {
 
	Agents = append(Agents, SalesAgent01)
	Agents = append(Agents, SalesAgent02)

	AllCustomers = initlizeMongoConnection()
	initlizeRouter()
}

func initlizeRouter() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	r := mux.NewRouter()

	r.HandleFunc("/", Welcome).Methods("GET")
	r.HandleFunc("/api/login", agentLogin).Methods("GET")
	r.HandleFunc("/api/customers", getCustomers).Methods("GET")
	r.HandleFunc("/api/customer/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/api/customers", createCustomer).Methods("POST")
	r.HandleFunc("/api/customer/{id}", updateOneCustomer).Methods("PUT")
	r.HandleFunc("/api/customer/{id}", deleteOneCustomer).Methods("DELETE")
	r.HandleFunc("/api/customers", deleteAllCustomers).Methods("DELETE")

	fmt.Println("server is listening:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	msg := ` 

Welcome to my app built using GOLang + GorillaMux + MongoDB and deployed to Heroku	
	below are the routes : 
	 
	("/api/login", agentLogin).Methods("GET")
	("/api/customers", getCustomers).Methods("GET")
	("/api/customers/{id}", getCustomer).Methods("GET")
	("/api/customers", createCustomer).Methods("POST")
	("/api/customers/{id}", updateOneCustomer).Methods("PUT")
	("/api/customers/{id}", deleteOneCustomer).Methods("DELETE")
	("/api/customers", deleteAllCustomers).Methods("DELETE")
	
	for the routes you have to login first
	use payload:   
	{
		"userid": "SA01",   
		"password": "agent01"    
	} 
	or  
    {
		"userid": "SA01",  
		"password": "agent02"  
	}   
	you will get a JWT token as a response which is valid for 15 mins,
	add the "jwt" in the request header as key and use the token as value for all the routes
	note: you will NOT be redirected to login
	Thankyou :)


app server : https://go-gorillamux-mongodb-myapp01.herokuapp.com/
	`

	io.WriteString(w, msg)
}

type SalesAgent struct {
	AgentID       string `json:"agentid"`
	AgentName     string `json:"agentname"`
	AgentEmail    string `json:"agentemail"`
	AgentPassword string `json:"agentpassword"`
}

type SalesManager struct {
	ManagerID    string `json:"managerid"`
	ManagerName  string `json:"managername"`
	ManagerEmail string `json:"manageremail"`
}

var salesManager SalesManager = SalesManager{
	ManagerID:    "SM01",
	ManagerName:  "Shawn Mendy",
	ManagerEmail: "shawnm@gmail.com",
}

var SalesAgent01 SalesAgent = SalesAgent{
	AgentID:       "SA01",
	AgentName:     "John Doe",
	AgentEmail:    "johnDoe@example.com",
	AgentPassword: "agent01",
}

var SalesAgent02 SalesAgent = SalesAgent{
	AgentID:       "SA02",
	AgentName:     "Jessica Jones",
	AgentEmail:    "jessicajones@example.com",
	AgentPassword: "agent02",
}
