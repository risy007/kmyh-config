package config

import "go.uber.org/fx"

func NewNatsModule() fx.Option {
	return fx.Module("natsmq",
		fx.Provide(NewNatsService),
		fx.Invoke(RegisterMicroServices),
	)
}
