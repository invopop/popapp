// Package config stores the main application configuration.
package config

import (
	"github.com/invopop/client.go/gateway"
	"github.com/invopop/client.go/invopop"
	"github.com/invopop/configure"
	"github.com/rs/zerolog/log"
)

// ABOUT: Any new configuration that you add to the application should be added here.
// Remember to match the name with the name in config.yaml.

// Config is the configuration for the application.
type Config struct {
	*gateway.Config

	Invopop *invopop.Config `json:"invopop"`
}

// NewConfig instantiates a new configuration.
func NewConfig(file string) *Config {
	c := new(Config)
	if err := configure.Load(file, c); err != nil {
		log.Fatal().Err(err).Msg("loading configuration")
	}
	c.Log.Init(c.Name)
	return c
}
