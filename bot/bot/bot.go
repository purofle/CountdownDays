package bot

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.uber.org/fx"
	"gopkg.in/telebot.v3"
)

func NewBot(lc fx.Lifecycle) (*telebot.Bot, error) {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		slog.Error("failed to create new bot", "error", err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go b.Start()
			slog.Info("Bot started")

			return nil
		},
		OnStop: func(context.Context) error {
			b.Stop()
			slog.Info("Bot stopped")

			return nil
		},
	})

	return b, nil
}
