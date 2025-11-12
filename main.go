// Package main contains the entry point of the application.
package main

import (
	"context"

	"github.com/invopop/client.go/invopop"
	"github.com/invopop/client.go/pkg/runner"
	"github.com/invopop/popapp/internal/config"
	"github.com/invopop/popapp/internal/domain"
	"github.com/invopop/popapp/internal/interfaces/gateway"
	"github.com/invopop/popapp/internal/interfaces/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// ABOUT: The main function initializes the application.
// It starts the web server and the gateway service.
// This is done using a command-line interface.

const (
	configPath     = "./config/config.yaml"
	defaultWebPort = "8080"
)

// ABOUT: The app structure is used to encapsulate all the components of the application.
// It provides a convenient way to access and manage these components.
// Any new domain config that is needed should be passed from the config object to the domain here.

// App keeps the main application logic in one place.
type App struct {
	config *config.Config
	ic     *invopop.Client
	web    *web.Service
	gw     *gateway.Service
	domain *domain.Setup
}

func main() {
	app := new(App)
	// Here initialize the app components
	app.config = config.NewConfig(configPath)
	app.ic = invopop.New(
		invopop.WithConfig(app.config.Invopop),
	)

	app.domain = domain.New() // pass as arguments the elements needed for the domain

	app.gw = gateway.New(app.config.Config, app.domain)

	app.run()
}

func (app *App) run() {
	cmd := app.root()
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("exiting")
	}
	log.Info().Msg("process terminated")
}

func (app *App) root() *cobra.Command {
	root := &cobra.Command{Use: app.config.Name}

	root.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Start accepting tasks and serving HTML assets",
		Long:  "",
		Args:  nil,
		Run: func(_ *cobra.Command, _ []string) {
			app.serve()
		},
	})

	return root
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
