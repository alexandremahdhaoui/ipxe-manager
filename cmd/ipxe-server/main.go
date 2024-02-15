package main

import (
	"errors"
	"github.com/alexandremahdhaoui/ipxe-api/internal/cmd"
	"github.com/alexandremahdhaoui/ipxe-api/internal/interface/server"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	Name = "ipxe-api"

	APIServerPort     = 8080
	MetricsServerPort = 8081
	ProbesServerPort  = 8082
)

func main() {
	var handler server.ServerInterface

	api := echo.New()
	api.Use(echoprometheus.NewMiddleware(Name))
	server.RegisterHandlers(api, handler)

	metrics := echo.New()
	metrics.GET("/metrics", echoprometheus.NewHandler())

	//TODO: create func initializing a probe server which returns non-200 response when server is considered Unhealthy.

	probes := echo.New()
	probes.GET("/healthz", func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "health probe error")
	})
	probes.GET("/readyz", func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "readiness probe error")
	})

	if err := cmd.Serve(map[int]*echo.Echo{
		APIServerPort:     api,
		MetricsServerPort: metrics,
		ProbesServerPort:  probes,
	}); err != nil {
		api.Logger.Fatal(err)
	}

	api.Logger.Infof("Successfully stopped %s server", Name)
}

func wrapErr(err error, s string) error {
	if err != nil {
		return errors.Join(err, errors.New(s))
	}

	return nil
}
