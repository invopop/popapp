// Package main contains the entry point of the application.
package main

import (
	"context"

	"github.com/invopop/client.go/pkg/runner"
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

func main() {
	a := newApp(configPath)

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Start the service",
		Run: func(_ *cobra.Command, _ []string) {
			if err := a.serve(); err != nil {
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

// Serve starts the main service.
func (app *App) serve() error {
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
