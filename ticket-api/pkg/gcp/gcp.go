package gcp

import (
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Module(
	"gcp",
	fx.Provide(
		NewBucketService,
		NewBucketStorage,
	),
)
