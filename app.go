package main

import (
	"context"

	"github.com/invopop/client.go/invopop"
	"github.com/invopop/client.go/pkg/runner"
	"github.com/invopop/cron/internal/config"
	"github.com/invopop/cron/internal/domain"
	"github.com/invopop/cron/internal/interfaces/gateway"
	"github.com/invopop/cron/internal/interfaces/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

/*
 * ABOUT: The app structure is used to encapsulate all the components of the application.
 * It provides a convenient way to access and manage these components.
 */

const (
	defaultWebPort = "8080"
)

// App keeps the main application logic in one place.
type App struct {
	conf   *config.Config
	ic     *invopop.Client
	web    *web.Service
	gw     *gateway.Service
	domain *domain.Setup
}

// New instantiates a new App.
func New(conf *config.Config) *App {
	app := new(App)
	app.conf = conf
	app.ic = invopop.New(
		invopop.WithConfig(app.conf.Invopop),
	)
	app.domain = domain.New()
	app.gw = gateway.New(app.conf.Config, app.domain)
	app.web = web.New(app.domain)
	return app
}

// AddCommands will add the application's own cobra commands to the provided
// root.
func (app *App) AddCommands(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Start accepting tasks and serving HTML assets",
		Long:  "",
		Args:  nil,
		Run: func(_ *cobra.Command, _ []string) {
			app.serve()
		},
	})
}

func (app *App) serve() {
	app.web = web.New(app.domain)

	run := new(runner.Group)

	run.Start(func(_ context.Context) error {
		return app.web.Start(defaultWebPort)
	})
	run.Start(func(_ context.Context) error {
		return app.gw.Start()
	})
	run.Stop(func(ctx context.Context) error {
		return app.web.Stop(ctx)
	})
	run.Stop(func(_ context.Context) error {
		app.gw.Stop()
		return nil
	})

	if err := run.Wait(); err != nil {
		log.Error().Err(err).Msg("shutting down")
	}
}
