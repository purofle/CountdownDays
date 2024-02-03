package api

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("api",
		fx.Provide(NewClient),
	)
}
