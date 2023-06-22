package internal

import (
	// "ticket-api/internal/chat"
	// "ticket-api/internal/event"
	// "ticket-api/internal/event_participant"
	// "ticket-api/internal/inquiry"
	// "ticket-api/internal/message_attachment"
	// "ticket-api/internal/recaptcha"
	// "ticket-api/internal/subscription"
	"ticket-api/internal/upload"
	"ticket-api/internal/user"
	// "ticket-api/internal/user_payment"

	"go.uber.org/fx"
)

var WebModule = fx.Module(
	"internals",
	fx.Options(
		upload.Module,
		user.Module,
	),
)
