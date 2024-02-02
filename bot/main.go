package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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

		countdown := &Countdown{
			TelegramId: c.Sender().ID,
			Name:       rawCommand[2],
			Date:       rawCommand[1],
		}

		newCountdown, err := NewCountdown(countdown)
		if err != nil {
			return err
		}

		return c.Send(newCountdown)
	})

	b.Handle("/all", func(c tele.Context) error {
		countdowns, err := GetAllCountdown(c.Sender())
		if err != nil {
			return err
		}

		text := "你的倒计时：\n"
		for _, c := range countdowns {
			text += fmt.Sprintf("%s: %s\n", c.Name, c.Date)
		}

		return c.Send(text)
	})

	b.Handle(tele.OnQuery, func(c tele.Context) error {

		countdowns, err := GetAllCountdown(c.Sender())
		if err != nil {
			return err
		}

		results := make(tele.Results, len(countdowns))
		for i, c := range countdowns {

			startDate, _ := time.Parse(time.DateOnly, c.Date)
			days := int(time.Now().Sub(startDate).Hours() / 24)
			middleStr := "已经"
			if days < 0 {
				days = -days
				middleStr = "还有"
			}

			text := fmt.Sprintf("%s: %d天", c.Name, days)

			results[i] = &tele.ArticleResult{
				Title: text,
				Text:  fmt.Sprintf("距离 %s %s%d 天\n\n日期：%s", c.Name, middleStr, days, c.Date),
			}
			results[i].SetResultID(strconv.Itoa(i))
		}

		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 60,
		})
	})

	b.Start()
}
