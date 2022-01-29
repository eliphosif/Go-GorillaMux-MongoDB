package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
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
}

var userLogin UserLogin
var userLoginError UserLoginError = UserLoginError{
	Message:     "unauthorized",
	Description: "it seems you're not logged in, Please login ang try again ! :)",
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
