package bot

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
	"go.uber.org/fx"

	"github.com/purofle/countdowndays/api"
)

type Handlers interface {
	HelloHandler(bot *telego.Bot, update telego.Update)
	AddHandler(bot *telego.Bot, update telego.Update)
	AllHandler(bot *telego.Bot, update telego.Update)
	QueryHandler(bot *telego.Bot, update telego.Update)
}

type handlers struct {
	fx.In
	ApiClient api.Client
}

func (h *handlers) HelloHandler(bot *telego.Bot, update telego.Update) {
	if _, err := bot.SendMessage(telegoutil.Message(telegoutil.ID(update.Message.Chat.ID), "Hello, World!")); err != nil {
		slog.Error("failed to send message", "error", err)
	}
}

func (h *handlers) AddHandler(bot *telego.Bot, update telego.Update) {
	id := update.Message.From.ID

	msgText := update.Message.Text
	rawCommand := strings.Split(msgText, " ")
	rawCommand = slices.DeleteFunc(rawCommand, func(s string) bool {
		return s == ""
	})

	if len(rawCommand) != 3 {
		bot.SendMessage(telegoutil.Message(telegoutil.ID(id), "请按照格式发送。格式：`/add 1989-06-04 名字`").WithParseMode(telego.ModeMarkdownV2))
		return
	}

	if _, err := time.Parse(time.DateOnly, rawCommand[1]); err != nil {
		bot.SendMessage(telegoutil.Message(telegoutil.ID(id), fmt.Sprintf("时间格式错误：%s", err)))
		return
	}

	if _, err := h.ApiClient.GetUser(id); err != nil {
		// 添加用户到数据库
		user, err := h.ApiClient.NewUser(&api.User{
			TelegramId: id,
			Username:   update.Message.From.Username,
			Name:       update.Message.From.FirstName,
		})
		if err != nil {
			bot.SendMessage(telegoutil.Message(telegoutil.ID(id), "创建用户失败"))
			slog.Error("failed to create new user", "error", err)
			return
		}

		slog.Info("New user: ", user)
	}

	countdown := &api.Countdown{
		TelegramId: id,
		Name:       rawCommand[2],
		Date:       rawCommand[1],
	}

	newCountdown, err := h.ApiClient.NewCountdown(countdown)
	if err != nil {
		bot.SendMessage(telegoutil.Message(telegoutil.ID(id), "创建倒计时失败"))
		slog.Error("failed to create new countdown", "error", err)
		return
	}

	bot.SendMessage(telegoutil.Message(telegoutil.ID(id), newCountdown))
}

func (h *handlers) AllHandler(bot *telego.Bot, update telego.Update) {
	id := update.Message.From.ID

	countdowns, err := h.ApiClient.GetAllCountdown(id)
	if err != nil {
		slog.Error("failed to get all countdowns", "error", err)
		bot.SendMessage(telegoutil.Message(telegoutil.ID(id), "获取倒计时失败"))
		return
	}

	text := "你的倒计时：\n"
	for _, c := range countdowns {
		text += fmt.Sprintf("%s: %s\n", c.Name, c.Date)
	}

	bot.SendMessage(telegoutil.Message(telegoutil.ID(id), text))
}

func (h *handlers) QueryHandler(bot *telego.Bot, update telego.Update) {
	id := update.InlineQuery.From.ID

	slog.Debug("inline query", "query", update.InlineQuery.Query, "username", update.InlineQuery.From.Username)

	countdowns, err := h.ApiClient.GetAllCountdown(id)
	if err != nil {
		slog.Error("failed to get all countdowns", "error", err)
		return
	}

	results := make([]telego.InlineQueryResult, len(countdowns))
	for i, c := range countdowns {
		startDate, _ := time.Parse(time.DateOnly, c.Date)
		days := int(time.Now().Sub(startDate).Hours() / 24)
		middleStr := "已经"
		if days < 0 {
			days = -days
			middleStr = "还有"
		}

		text := fmt.Sprintf("%s: %d天", c.Name, days)

		results[i] = telegoutil.ResultArticle(
			update.InlineQuery.ID,
			text,
			telegoutil.TextMessage(fmt.Sprintf("距离 %s %s%d天\n\n日期：%s", c.Name, middleStr, days, c.Date)),
		)
	}

	bot.AnswerInlineQuery(telegoutil.InlineQuery(update.InlineQuery.ID, results...).WithCacheTime(0).WithIsPersonal())
}

func NewHandlers(h handlers) Handlers {
	return &h
}

func BindHandlers(h *telegohandler.BotHandler, handlers Handlers) {
	group := h.Group()
	group.Use(telegohandler.PanicRecovery())
	// TODO: limit in pm

	group.Handle(handlers.HelloHandler, telegohandler.CommandEqual("hello"))
	group.Handle(handlers.AddHandler, telegohandler.CommandEqual("add"))
	group.Handle(handlers.AllHandler, telegohandler.CommandEqual("all"))
	group.Handle(handlers.QueryHandler, telegohandler.AnyInlineQuery())
}
