package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	tele "gopkg.in/telebot.v3"
)

func client() *resty.Client {
	return resty.New().
		SetBaseURL("http://localhost:8080").
		SetBasicAuth("bot", "auth_token")

}

type Countdown struct {
	TelegramId int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Date       string `json:"date"`
}

func NewCountdown(countdown *Countdown) (string, error) {
	response, err := client().R().
		SetBody(countdown).
		Post("/countdown")

	if err != nil {
		return "", err
	}

	return string(response.Body()), nil
}

type CountdownResponse struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Date            string `json:"date"`
	Description     string `json:"description"`
	ShowAnniversary bool   `json:"show_anniversary"`
	Owner           struct {
		Id         int    `json:"id"`
		TelegramId int64  `json:"telegram_id"`
		Username   string `json:"username"`
		Name       string `json:"name"`
	} `json:"owner"`
}

func GetAllCountdown(user *tele.User) ([]CountdownResponse, error) {
	response, err := client().R().
		Get(fmt.Sprintf("/countdown/%d/all", user.ID))

	if err != nil {
		return nil, err
	}

	var countdowns []CountdownResponse
	err = json.Unmarshal(response.Body(), &countdowns)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("no countdown found")
	}

	return countdowns, nil
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
