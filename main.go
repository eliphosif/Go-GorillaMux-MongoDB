package main

var Customers []Customer
var Agents []SalesAgent

var loggedinAgent string
func main() {

	c1 := Customer{
		//CustId:    time.Now().Format("20060102150405"),
		CustId:    0,
		FirstName: "kishore",
		LastName:  "ch",
		Email:     "kishore@gmail.com",
		Age:       25,
		Address: Address{
			Hno:    12,
			Street: "lb nagar",
			State:  "Telangna"},
		CreatedbyAgentID: "SM01",
	}

	Customers = append(Customers, c1)

	Agents = append(Agents, SalesAgent01)
	Agents = append(Agents, SalesAgent02)

	AllCustomers = initlizeMongoConnection()
	initlizeRouter()

}
