package main

import (
	"context"
	"log"

	"github.com/overlorddamygod/go-auth/configs"
	"github.com/overlorddamygod/go-auth/controllers/auth"
	"github.com/overlorddamygod/go-auth/controllers/auth/admin"
	"github.com/overlorddamygod/go-auth/db"
	"github.com/overlorddamygod/go-auth/mailer"
	"github.com/overlorddamygod/go-auth/middlewares"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	app := fx.New(
		fx.Provide(
			configs.NewConfig(".env"),
			db.NewDB,
			mailer.NewMailer,
			models.NewLogger,
			middlewares.NewLimiter,
			auth.NewAuthController,
			admin.NewAdminController,
			server.NewRouter,
		),
		fx.Populate(&configs.MainConfig),
		fx.Invoke(server.RegisterServer),
		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
	)

	startCtx := context.Background()

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	// <-app.Done()
}
