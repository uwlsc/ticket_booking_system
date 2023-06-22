package bootstrap

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"ticket-api/config"
	"ticket-api/internal/seeds"
	"ticket-api/internal/user"
	"ticket-api/pkg"
	"ticket-api/pkg/firebase"
)

func SeedApplication() {
	app := fx.New(
		fx.Provide(
			config.New,
			pkg.GetLogger,
			pkg.NewDatabase,
		),
		user.Module,
		firebase.Module,
		seeds.Module,
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}
