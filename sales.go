package main

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
