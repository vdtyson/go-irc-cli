package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

type ChannelAccessType string // defines how users can join; channel owner can set this
const (
	CLOSED = ChannelAccessType("DIRECT_MESSAGE") // direct message
	INVITE = ChannelAccessType("INVITE")         // channel is invite only
	SECRET = ChannelAccessType("SECRET")         // can join through secret key
	OPEN   = ChannelAccessType("PUBLIC")         // can join through channel name
)

// http://localhost:8080/register - POST
type UserRegInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isUserAdmin"`
}

// http://localhost:8080/channel/new - POST
type NewGroupChannelInput struct {
	OwnerUsername string            `json:"ownerUsername"`
	ChannelName   string            `json:"channelName"`
	AccessType    ChannelAccessType `json:"accessType"`
}

// http://localhost:8080/channels/message - POST
type NewMessageInput struct {
	ChannelName    string `json:"channelName"`
	Message        string `json:"message"`
	SenderUsername string `json:"senderUsername"`
}

// http://localhost:8080/channels/messages - GET
type AllChannelMessagesInput struct {
	UserName    string `json:"username"`
	ChannelName string `json:"channelName"`
}

// http://localhost:8080/channels/users/kick - POST TODO: Not yet implemented
type KickUserInput struct {
	ChannelName   string `json:"channelName"`
	OwnerUsername string `json:"ownerUsername"`
	UserToKick    string `json:"userToKick"`
}

// http://localhost:8080/channels/users - PUT
type AddUserToChannelInput struct {
	ChannelName   string `json:"channelName"`
	OwnerUsername string `json:"ownerUsername"`
	UserToAdd     string `json:"userToAdd"`
	PrivilegeType string `json:"privilegeType"`
}

// http://localhost:8080/admin/ban - PUT
type BanUserInput struct {
	AdminUsername     string `json:"adminUsername"`
	UserToBanUsername string `json:"userToBanUsername"`
}

/* PATHS without body */

// New DM Channel: http://localhost:8080/channels/direct/{username1}/{username2} - POST
// All Channels a user is member of: http://localhost:8080/users/{username}/channels - GET

const (
	MAIN_PATH          = "http://localhost:8080/"
	HELLO_PATH         = "http://localhost:8080/"
	REGISTER_USER_PATH = MAIN_PATH + "register"
)

func main() {
	initClient()

}

func initClient() {
	client = resty.New()
}

func registerUser(userName, email, password string, isAdmin bool) int { // http://localhost:8080/register - POST
	userRegInput := UserRegInput{userName, email, password, isAdmin}
	data, err := json.Marshal(&userRegInput)
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
	return registerResp.StatusCode()

}
