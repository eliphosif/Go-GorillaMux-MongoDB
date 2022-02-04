# my-heroku-app-repoUI

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
