package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/add", func(c tele.Context) error {

		c.Sender()

		msgText := c.Message().Text
		rawCommand := strings.Split(msgText, " ")
		rawCommand = slices.DeleteFunc(rawCommand, func(s string) bool {
			return s == ""
		})

		if len(rawCommand) != 3 {
			return c.Send("请按照格式发送。格式：`/add 1989-06-04 名字`", tele.ModeMarkdownV2)
		}

		_, err := time.Parse(time.DateOnly, rawCommand[1])
		if err != nil {
			return c.Send(fmt.Sprintf("时间格式错误：%s", err))
		}

		_, err = GetUser(c.Sender())
		if err != nil {
			// 添加用户到数据库
			user, err := NewUser(c.Sender())
			if err != nil {
				return err
			}
			log.Println("New user: ", user)
		}

		return c.Send("")
	})

	b.Start()
}
