package bootstrap

import (
	"context"
	"time"

	"ticket-api/config"
	"ticket-api/internal"
	"ticket-api/internal/middlewares"
	"ticket-api/internal/seeds"
	"ticket-api/pkg"
	"ticket-api/pkg/firebase"
	"ticket-api/pkg/gcp"
	"ticket-api/pkg/stripe"

	"go.uber.org/fx"
)

func WebApplication() {
	app := fx.New(
		fx.Provide(
			config.New,
			pkg.GetLogger,
			pkg.NewHTTPServer,
			pkg.NewDatabase,
			pkg.NewMigrations,
		),
		internal.WebModule,
		middlewares.Module,
		gcp.Module,
		firebase.Module,
		stripe.Module,
		seeds.Module,

		fx.Invoke(loadSpecifiedLocation),
		fx.Invoke(runMigration),
		fx.Invoke(startWebServer),
	)
	app.Run()
}

func loadSpecifiedLocation(lc fx.Lifecycle, config config.Config) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				loc, _ := time.LoadLocation(config.TimeZone)
				time.Local = loc
				return nil
			},
		},
	)
}

func startWebServer(lifecycle fx.Lifecycle, server pkg.HTTPServer, config config.Config, logger pkg.Logger) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					server.Start(config, logger)
				}()
				return nil
			},
			OnStop: func(context context.Context) error {
				logger.Info("Stopping Application")
				return nil
			},
		},
	)
}

func runMigration(lc fx.Lifecycle, migration pkg.Migrations, logger pkg.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Setting Up migration")
				migration.Migrate()
				return nil
			},
		},
	)
}
