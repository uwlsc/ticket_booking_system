package seeds

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"seeds", fx.Options(
		fx.Invoke(NewAdminSeed),
	),
)
