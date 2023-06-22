package firebase

import (
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Module(
	"firebase",
	fx.Provide(
		NewFBApp,
		NewFBAuth,
		NewFCMClient,
		NewFirestoreClient,
		NewFirebaseService,
	),
)
