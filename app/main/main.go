package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

var client *resty.Client

// http://localhost:8080/register - POST
type UserRegInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isUserAdmin"`
}

const (
	MAIN_PATH          = "http://localhost:8080/"
	HELLO_PATH         = "http://localhost:8080/"
	REGISTER_USER_PATH = "http://localhost:8080/register"
)

func main() {
	initClient()
	userRegInput := UserRegInput{"tester7", "tester7@gmail.com", "testedkajfkl", false}
	data, err := json.MarshalIndent(&userRegInput, "", "    ")
	if err != nil {
		panic(err)
	}

	registerResp, err := client.R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		Post(REGISTER_USER_PATH)

	if err != nil {
		panic(err)
	}

	fmt.Println(registerResp.Header())
}

func initClient() {
	client = resty.New()
}
