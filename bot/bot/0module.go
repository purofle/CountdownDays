package bot

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"bot",
		fx.Provide(NewBot),
		fx.Provide(NewHandlers),
		fx.Invoke(BindHandlers),
	)
}
