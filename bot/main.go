package main

import (
	"go.uber.org/fx"

	"github.com/purofle/countdowndays/api"
	"github.com/purofle/countdowndays/bot"
)

func main() {
	app := fx.New(
		api.Module(),
		bot.Module(),
	)

	app.Run()
}
