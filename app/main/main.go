package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

var httpClient *http.Client
var client *resty.Client
var username = "vdtyson"

type ChannelAccessType string // defines how users can join; channel owner can set this

// http://localhost:8080/register - POST
type UserRegInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isUserAdmin"`
}

// http://localhost:8080/channel/new - POST
type NewGroupChannelInput struct {
	OwnerUsername string `json:"ownerUsername"`
	ChannelName   string `json:"channelName"`
	AccessType    string `json:"accessType"`
}

// http://localhost:8080/channels/message - POST
type NewMessageInput struct {
	ChannelName    string `json:"channelName"`
	Message        string `json:"message"`
	SenderUsername string `json:"senderUsername"`
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

// http://localhost:8080/channels/messages - GET
type AllChannelMessagesInput struct {
	UserName    string `json:"username"`
	ChannelName string `json:"channelName"`
}

/* PATHS without body */

// New DM Channel: http://localhost:8080/channels/direct/{username1}/{username2} - POST
// All Channels a user is member of: http://localhost:8080/users/{username}/channels - GET

const (
	//
	BASE_URL                   = "https://mthree-go-irc.herokuapp.com"
	WELCOME_ENDPOINT           = "/"
	REGISTER_ENDPOINT          = "/register"
	BAN_ENDPOINT               = "/admin/ban"
	INVITE_USER_ENDPOINT       = "/channels/users"
	ALL_CHAN_MESSAGES_ENDPOINT = "/channels/messages"
	USER_CHANNELS_ENDPOINT     = "/users/{username}/channels"
	NEW_DM_ENDPOINT            = "/channels/direct/{username1}/{username2}"
	NEW_CHANNEL_ENDPOINT       = "/channels/new"
	NEW_MESSAGE_ENDPOINT       = "/channels/message"
)

func main() {
	initClient()
	input := NewMessageInput{"#midnight-coders", "Coding at midnight is very fun", "vdtyson"}
	err := addMessageToChannel(input)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func initClient() {
	client = resty.New()
	client.HostURL = BASE_URL
}

func addMessageToChannel(input NewMessageInput) error {
	resp, err := client.R().SetBody(input).Post(NEW_MESSAGE_ENDPOINT)
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		fmt.Println("Message added!")
	} else {
		fmt.Println(string(resp.Body()))
	}
	return nil
}
func newGroupChannel(input NewGroupChannelInput) error {
	resp, err := client.R().SetBody(input).Post(NEW_CHANNEL_ENDPOINT)
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		fmt.Println("New channel created!")
	} else {
		fmt.Println(string(resp.Body()))
	}
	return nil
}

func newDirectMessageChannel(username1, username2 string) error {
	resp, err := client.R().SetPathParams(map[string]string{"username1": username1, "username2": username2}).Post(NEW_DM_ENDPOINT)
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		fmt.Println("new dm channel created")
	} else {
		fmt.Println(resp.Body())
	}

	return nil
}

func getUserChannels(username string) error {
	resp, err := client.R().SetPathParams(map[string]string{"username": username}).Get(USER_CHANNELS_ENDPOINT)
	if err != nil {
		return err
	} else {
		fmt.Println(resp.String())
		return nil
	}
}

func getAllChannelMessages(input *AllChannelMessagesInput) error {
	resp, err := client.R().
		SetBody(input).
		Put(ALL_CHAN_MESSAGES_ENDPOINT)

	if err != nil {
		return err
	}

	fmt.Print(resp.String())
	return nil
}

func inviteUser(input *AddUserToChannelInput) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	inviteUserResp, err := client.R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		Put(INVITE_USER_ENDPOINT)

	if err != nil {
		return err
	}

	if inviteUserResp.StatusCode() == 200 {
		fmt.Println("User invited!")
	} else {
		fmt.Println(inviteUserResp.String())
	}

	return nil
}
func banUser(input BanUserInput) error {
	data, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	banResp, err := client.R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		Put(BAN_ENDPOINT)

	if err != nil {
		return err
	}

	fmt.Println(banResp.String())

	return nil
}

func registerUser(input UserRegInput) error { // http://localhost:8080/register - POST
	data, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	registerResp, err := client.R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		Post(REGISTER_ENDPOINT)

	if err != nil {
		return err
	}

	if registerResp.StatusCode() == 200 {
		fmt.Println("User registered!")
	} else {
		fmt.Println(registerResp.String())
	}

	return nil
}

func changeUser(newUsername string) {
	username = newUsername
}
