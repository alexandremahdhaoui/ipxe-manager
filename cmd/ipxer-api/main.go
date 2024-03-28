package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexandremahdhaoui/ipxer/internal/cmd"
	"github.com/alexandremahdhaoui/ipxer/internal/driver/server"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

const (
	Name = "ipxer-api"

	AppServerPort     = 8080 //TODO: specify through config file
	MetricsServerPort = 8081 //TODO: specify through config file
	ProbesServerPort  = 8082 //TODO: specify through config file
)

var (
	Version   = "dev" //nolint:gochecknoglobals // set by ldflags
	CommitSHA = "dev" //nolint:gochecknoglobals // set by ldflags
)

// ------------------------------------------------- Main ----------------------------------------------------------- //

func main() {
	fmt.Printf("Starting %s version %s (%s)\n", Name, Version, CommitSHA)

	// --------------------------------------------- App ------------------------------------------------------------ //

	var handler server.ServerInterface //TODO

	app := echo.New()
	app.Use(echoprometheus.NewMiddleware(Name))
	server.RegisterHandlers(app, handler)

	// --------------------------------------------- Metrics -------------------------------------------------------- //

	metrics := echo.New()
	metrics.GET("/metrics", echoprometheus.NewHandler())

	// --------------------------------------------- Probes --------------------------------------------------------- //

	// TODO: create func initializing a probe server which returns non-200 response when server is considered Unhealthy.

	probes := echo.New()
	probes.GET("/healthz", func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "health probe error")
	})
	probes.GET("/readyz", func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "readiness probe error")
	})

	// --------------------------------------------- Run Server ----------------------------------------------------- //

	if err := cmd.Serve(map[int]*echo.Echo{
		AppServerPort:     app,
		MetricsServerPort: metrics,
		ProbesServerPort:  probes,
	}); err != nil {
		app.Logger.Fatal(err)
	}

	app.Logger.Infof("Successfully stopped %s server", Name)
}

// --------------------------------------------- UTILS -------------------------------------------------------------- //

func wrapErr(err error, s string) error {
	if err != nil {
		return errors.Join(err, errors.New(s))
	}

	return nil
}
