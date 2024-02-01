package main

import (
	"github.com/go-resty/resty/v2"
	tele "gopkg.in/telebot.v3"
)

func PostCountdown() {

}

func NewUser(user *tele.User) {
}

func GetUser(user *tele.User) {
	client := resty.New()

	client.R()
}
