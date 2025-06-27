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

const (
	configPath = "./config/config.yaml"
)

// App keeps the main application logic in one place.
type App struct {
	conf   *config.Config
	ic     *invopop.Client
	web    *web.Service
	gw     *gateway.Service
	domain *domain.Setup
}

func newApp(configPath string) *App {
	app := new(App)
	app.conf = config.NewConfig(configPath)

	app.ic = invopop.New(
		invopop.WithConfig(app.conf.Invopop),
	)

	app.domain = domain.New()

	app.gw = gateway.New(app.conf.Config, app.domain)
	app.web = web.New(app.domain)

	return app
}

func main() {
	a := newApp(configPath)

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := a.serve(cmd.Context()); err != nil {
				log.Fatal().Err(err).Msg("starting the service")
			}
		},
	}

	root := &cobra.Command{Use: a.conf.Name}
	root.AddCommand(cmdServe)
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msg("exiting")
	}

	log.Info().Msg("process terminated")
}

const (
	defaultWebPort = "8080"
)

// Serve starts the main service.
func (app *App) serve(ctx context.Context) error {
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

	return run.Wait()
}
