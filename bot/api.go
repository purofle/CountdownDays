package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	tele "gopkg.in/telebot.v3"
)

func PostCountdown() {

}

func NewUser(user *tele.User) {
}

func GetUser(user *tele.User) (string, error) {
	client := resty.New()
	client.SetBaseURL("http://localhost:8080")

	response, err := client.R().
		Get(fmt.Sprintf("/user/%d", user.ID))
	if err != nil {
		return "", err
	}

	fmt.Println("Response: ", response)

	return "", nil
}
