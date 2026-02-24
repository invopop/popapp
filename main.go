// Package main contains the entry point of the application.
package main

import (
	"github.com/invopop/cron/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

/*
 * ABOUT: The main function initializes and passes control over to
 * the applicaiton.
 */

const (
	configPath = "./config/config.yaml"
)

func main() {
	conf := config.NewConfig(configPath)

	app := New(conf)

	root := &cobra.Command{Use: app.conf.Name}
	app.AddCommands(root)
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msg("exiting")
	}

	log.Info().Msg("process terminated")
}
