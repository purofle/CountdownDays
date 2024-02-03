package api

import (
	"github.com/go-resty/resty/v2"
)

type Client interface {
	NewCountdown(countdown *Countdown) (string, error)
	GetAllCountdown(userId int64) ([]CountdownResponse, error)
	NewUser(user *User) (string, error)
	GetUser(userId int64) (string, error)
}

type client struct {
	*resty.Client
}

func NewClient() Client {
	c := resty.New().
		SetBaseURL("http://localhost:8080").
		SetBasicAuth("bot", "auth_token")

	return &client{c}
}
