package upload

import "go.uber.org/fx"

var Module = fx.Module(
	"upload", fx.Options(
		fx.Provide(NewUploadService),
	),
)
