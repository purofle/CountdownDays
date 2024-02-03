package api

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/pkg/errors"
)

type Countdown struct {
	TelegramId int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Date       string `json:"date"`
}

func (c *client) NewCountdown(countdown *Countdown) (string, error) {
	response, err := c.R().
		SetBody(countdown).
		Post("/countdown")
	if err != nil {
		return "", err
	}

	if response.StatusCode() != 200 {
		slog.Error("failed to create new countdown", "response", response.String())
		return "", errors.New("countdown not created")
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

func (c *client) GetAllCountdown(userId int64) ([]CountdownResponse, error) {
	response, err := c.R().
		SetPathParam("userId", strconv.FormatInt(userId, 10)).
		Get("/countdown/{userId}/all")
	if err != nil {
		slog.Error("failed to get all countdowns", "error", err)
		return nil, err
	}

	var countdowns []CountdownResponse
	if err := json.Unmarshal(response.Body(), &countdowns); err != nil {
		slog.Error("failed to unmarshal countdowns", "error", err)
		return nil, err
	}

	if response.StatusCode() != 200 {
		slog.Error("failed to get all countdowns", "response", response.String())
		return nil, errors.New("countdowns not found")
	}

	return countdowns, nil
}

type User struct {
	TelegramId int64  `json:"telegram_id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
}

func (c *client) NewUser(user *User) (string, error) {
	response, err := c.R().
		SetBody(user).
		Post("/user")
	if err != nil {
		slog.Error("failed to create new user", "error", err)
		return "", err
	}

	if response.StatusCode() != 200 {
		slog.Error("failed to create new user", "response", response.String())
		return "", errors.New("user not created")
	}

	return string(response.Body()), nil
}

func (c *client) GetUser(userId int64) (string, error) {
	response, err := c.R().
		SetPathParam("userId", strconv.FormatInt(userId, 10)).
		Get("/user/{userId}")
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return "", err
	}

	if response.StatusCode() != 200 {
		slog.Error("failed to get user", "response", response.String())
		return "", errors.New("user not found")
	}

	return string(response.Body()), nil
}
