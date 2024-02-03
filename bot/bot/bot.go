package bot

import (
	"context"
	"log/slog"
	"os"

	"go.uber.org/fx"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func NewBot() (*telego.Bot, error) {
	b, err := telego.NewBot(os.Getenv("TOKEN"), telego.WithHealthCheck())
	if err != nil {
		slog.Error("failed to create new bot", "error", err)
		return nil, err
	}

	return b, nil
}

func NewUpdateChannel(bot *telego.Bot) (<-chan telego.Update, error) {
	ch, err := bot.UpdatesViaLongPolling(&telego.GetUpdatesParams{
		Offset:  0,
		Timeout: 60,
		Limit:   100,
	})
	if err != nil {
		slog.Error("failed to start long polling", "error", err)
		return nil, err
	}

	return ch, nil
}

func NewHandler(bot *telego.Bot, updates <-chan telego.Update, lc fx.Lifecycle) (*telegohandler.BotHandler, error) {
	h, err := telegohandler.NewBotHandler(bot, updates)
	if err != nil {
		slog.Error("failed to create new bot handler", "error", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go h.Start()
			slog.Info("bot started")

			return nil
		},
		OnStop: func(context.Context) error {
			h.Stop()
			slog.Info("bot stopped")

			return nil
		},
	})

	return h, nil
}
