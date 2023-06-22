package user

import "go.uber.org/fx"

var Module = fx.Module(
	"user", fx.Options(
		fx.Provide(
			NewHandler,
			NewRepository,
			NewService,
			
		),
		fx.Invoke(NewRoutes),
	),
)
