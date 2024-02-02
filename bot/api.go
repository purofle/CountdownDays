package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	tele "gopkg.in/telebot.v3"
)

func client() *resty.Client {
	return resty.New().
		SetBaseURL("http://localhost:8080").
		SetBasicAuth("bot", "auth_token")

}

func PostCountdown() {

}

func NewUser(user *tele.User) (string, error) {

	type User struct {
		TelegramId int64  `json:"telegram_id"`
		Username   string `json:"username"`
		Name       string `json:"name"`
	}

	response, err := client().R().
		SetBody(User{
			TelegramId: user.ID,
			Username:   user.Username,
			Name:       user.FirstName + user.LastName,
		}).
		Post("/user")

	if err != nil {
		return "", err
	}

	return string(response.Body()), nil
}

func GetUser(user *tele.User) (string, error) {
	response, err := client().R().
		Get(fmt.Sprintf("/user/%d", user.ID))
	if err != nil {
		return "", err
	}

	if response.StatusCode() != 200 {
		return "", fmt.Errorf("user not found")
	}

	return string(response.Body()), nil
}
