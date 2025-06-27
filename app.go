package main

import (
	"github.com/invopop/client.go/invopop"
	"github.com/invopop/popapp/internal/config"
	"github.com/invopop/popapp/internal/domain"
	"github.com/invopop/popapp/internal/interfaces/gateway"
	"github.com/invopop/popapp/internal/interfaces/web"
)

// ABOUT: The app structure is used to encapsulate all the components of the application.
// It provides a convenient way to access and manage these components.
// Any new domain config that is needed should be passed from the config object to the domain here.

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
