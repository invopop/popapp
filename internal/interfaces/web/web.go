// Package web contains the web interface for the application.
package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/invopop/client.go/pkg/echopop"
	"github.com/invopop/popapp/internal/domain"
	"github.com/invopop/popapp/internal/interfaces/web/assets"
	popui "github.com/invopop/popui/go"
	"github.com/labstack/echo/v4"
)

// ABOUT: The web service uses an echo server under the hood.
// It lets you define routes and controllers for the application.
// The controller is only responsible for parsing the request and calling the appropriate domain method.
// Once this is done, the controller should return the correct response.
// Controllers should be added to the service struct and initialized during the serve call.
// Each controller struct should contain the necessary domain structs.

// Service holds together the web service.
type Service struct {
	echo *echopop.Service
}

// New instantiates a new web service.
func New(domain *domain.Setup) *Service {
	s := new(Service)

	s.echo = echopop.NewService()
	s.echo.Serve(func(e *echo.Echo) {
		e.StaticFS(popui.AssetPath, popui.Assets)
		e.StaticFS("/", assets.Content)
	})

	return s
}

// Start the web service.
func (s *Service) Start(port string) error {
	return s.echo.Start(port)
}

// Stop the web service.
func (s *Service) Stop(ctx context.Context) error {
	return s.echo.Stop(ctx)
}

// ServeHTTP exposes the echo's ServeHTTP function which is useful for testing.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.Serve(func(e *echo.Echo) {
		e.ServeHTTP(w, r)
	})
}

// render provides a wrapper around the component to make it nice to render.
func render(c echo.Context, status int, t templ.Component) error { //nolint:unparam
	c.Response().Writer.WriteHeader(status)

	if err := t.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func clientError(format string, a ...any) error {
	return httpError(http.StatusBadRequest, format, a...)
}

func serverError(format string, a ...any) error {
	return httpError(http.StatusInternalServerError, format, a...)
}

func httpError(code int, format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	return echo.NewHTTPError(code, echo.Map{"error": err.Error()})
}
