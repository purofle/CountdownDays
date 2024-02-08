package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
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

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("/add 1989-06-04 标题\n/all 查看所有id\n/delete + id 删除")
	})

	b.Handle("/add", func(c tele.Context) error {

		rawCommand := GetRawCommand(c.Message().Text)

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

	b.Handle("/delete", func(c tele.Context) error {

		rawCommand := GetRawCommand(c.Message().Text)

		if len(rawCommand) != 2 {
			return c.Send("请按照格式发送。格式：`/delete 1`", tele.ModeMarkdownV2)
		}

		id, err := strconv.Atoi(rawCommand[1])
		if err != nil {
			return c.Send(fmt.Sprintf("ID格式错误：%s", err))
		}

		_, err = DeleteCountdown(c.Sender(), id)
		if err != nil {
			return c.Send(fmt.Sprintf("删除ID为%d的倒计时失败：%s", id, err))
		}

		return c.Send(fmt.Sprintf("删除ID为%d的倒计时成功", id))
	})

	b.Handle("/all", func(c tele.Context) error {
		countdowns, err := GetAllCountdown(c.Sender())
		if err != nil {
			return err
		}

		text := "你的倒计时：\n"
		for _, c := range countdowns {
			text += fmt.Sprintf("%s(%d): %s\n", c.Name, c.Id, c.Date)
		}

		return c.Send(text)
	})

	b.Handle(tele.OnQuery, func(ctx tele.Context) error {

		log.Println(ctx.Sender().Username, "searching")

		countdowns, err := GetAllCountdown(ctx.Sender())
		if err != nil {
			countdowns = []CountdownResponse{}
		}

		results := make(tele.Results, len(countdowns)+1)

		if len(countdowns) == 0 {
			results[0] = &tele.ArticleResult{
				Title: "tips: 在 inline 里输入的内容在选择倒计时后都会被加到最前面哦～",
				Text:  "你好，我是倒计时 bot，我的 id 是 6872551455，快来私聊 @countdown_days_bot 添加倒计时",
			}
		} else {
			results[0] = &tele.ArticleResult{
				Title: "tips: 在 inline 里输入的内容在选择倒计时后都会被加到最前面哦～",
				Text:  "喵呜",
			}
		}

		for i, c := range countdowns {

			startDate, _ := time.Parse(time.DateOnly, c.Date)
			startDate = startDate.In(time.Local) // SB UTC
			days := int(math.Round(time.Now().Sub(startDate).Hours() / 24))
			middleStr := "已经"
			if days < 0 {
				days = -days
				middleStr = "还有"
			}

			text := fmt.Sprintf("%s: %d天", c.Name, days)

			if len(countdowns) != 0 {
				results[i+1] = &tele.ArticleResult{
					Title: text,
					Text:  fmt.Sprintf("%s\n\n距离 %s %s%d天\n\n日期：%s", ctx.Query().Text, c.Name, middleStr, days, c.Date),
				}
				results[i+1].SetResultID(strconv.Itoa(i))
			}
		}

		return ctx.Answer(&tele.QueryResponse{
			Results:    results,
			CacheTime:  0,
			IsPersonal: true,
		})
	})

	b.Start()
}
