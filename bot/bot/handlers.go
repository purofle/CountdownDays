package bot

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.uber.org/fx"
	"gopkg.in/telebot.v3"

	"github.com/purofle/countdowndays/api"
)

type Handlers interface {
	HelloHandler(c telebot.Context) error
	AddHandler(c telebot.Context) error
	AllHandler(c telebot.Context) error
	QueryHandler(c telebot.Context) error
}

type handlers struct {
	fx.In
	ApiClient api.Client
}

func (h *handlers) HelloHandler(c telebot.Context) error {
	return c.Send("Hello!")
}

func (h *handlers) AddHandler(c telebot.Context) error {
	sender := c.Sender()
	if sender == nil || sender.IsBot {
		return c.Send("Sender not legal")
	}

	msgText := c.Message().Text
	rawCommand := strings.Split(msgText, " ")
	rawCommand = slices.DeleteFunc(rawCommand, func(s string) bool {
		return s == ""
	})

	if len(rawCommand) != 3 {
		return c.Send("请按照格式发送。格式：`/add 1989-06-04 名字`", telebot.ModeMarkdownV2)
	}

	if _, err := time.Parse(time.DateOnly, rawCommand[1]); err != nil {
		return c.Send(fmt.Sprintf("时间格式错误：%s", err))
	}

	if _, err := h.ApiClient.GetUser(sender.ID); err != nil {
		// 添加用户到数据库
		user, err := h.ApiClient.NewUser(&api.User{
			TelegramId: sender.ID,
			Username:   sender.Username,
			Name:       sender.FirstName,
		})
		if err != nil {
			return err
		}

		slog.Info("New user: ", user)
	}

	countdown := &api.Countdown{
		TelegramId: sender.ID,
		Name:       rawCommand[2],
		Date:       rawCommand[1],
	}

	newCountdown, err := h.ApiClient.NewCountdown(countdown)
	if err != nil {
		return err
	}

	return c.Send(newCountdown)
}

func (h *handlers) AllHandler(c telebot.Context) error {
	sender := c.Sender()
	if sender == nil || sender.IsBot {
		return c.Send("Sender not legal")
	}

	countdowns, err := h.ApiClient.GetAllCountdown(sender.ID)
	if err != nil {
		return err
	}

	text := "你的倒计时：\n"
	for _, c := range countdowns {
		text += fmt.Sprintf("%s: %s\n", c.Name, c.Date)
	}

	return c.Send(text)
}

func (h *handlers) QueryHandler(c telebot.Context) error {
	sender := c.Sender()
	if sender == nil || sender.IsBot {
		return c.Send("Sender not legal")
	}

	slog.Debug(sender.Username, "searching")

	countdowns, err := h.ApiClient.GetAllCountdown(sender.ID)
	if err != nil {
		return err
	}

	results := make(telebot.Results, len(countdowns))
	for i, c := range countdowns {

		startDate, _ := time.Parse(time.DateOnly, c.Date)
		days := int(time.Now().Sub(startDate).Hours() / 24)
		middleStr := "已经"
		if days < 0 {
			days = -days
			middleStr = "还有"
		}

		text := fmt.Sprintf("%s: %d天", c.Name, days)

		results[i] = &telebot.ArticleResult{
			Title: text,
			Text:  fmt.Sprintf("距离 %s %s%d天\n\n日期：%s", c.Name, middleStr, days, c.Date),
		}
		results[i].SetResultID(strconv.Itoa(i))
	}

	return c.Answer(&telebot.QueryResponse{
		Results:    results,
		CacheTime:  0,
		IsPersonal: true,
	})
}

func NewHandlers(h handlers) Handlers {
	return &h
}

func BindHandlers(bot *telebot.Bot, h Handlers) {
	bot.Handle("/hello", h.HelloHandler)
	bot.Handle("/add", h.AddHandler)
	bot.Handle("/all", h.AllHandler)
	bot.Handle(telebot.OnQuery, h.QueryHandler)
}
