package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/cmd"
	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/driver/server"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
)

const (
	Name = "ipxer-api"
)

var (
	Version   = "dev" //nolint:gochecknoglobals // set by ldflags
	CommitSHA = "dev" //nolint:gochecknoglobals // set by ldflags
)

// ------------------------------------------------- Main ----------------------------------------------------------- //

func main() {
	fmt.Printf("Starting %s version %s (%s)\n", Name, Version, CommitSHA)

	// --------------------------------------------- Config --------------------------------------------------------- //

	// TODO: specify all those values through config file
	//      which should derive from a k8s manifests/helmchart values.
	var (
		assignmentNamespace string
		profileNamespace    string

		appServerPort     = 8080
		metricsServerPort = 8081
		probesServerPort  = 8082
		metricsPath       = "/metrics"
		probesHealthPath  = "/healthz"
		probesReadyPath   = "/readyz"
	)

	// --------------------------------------------- Client --------------------------------------------------------- //

	var cl client.Client        // TODO
	var dynCl dynamic.Interface // TODO

	// --------------------------------------------- Adapter -------------------------------------------------------- //

	assignment := adapter.NewAssignment(cl, assignmentNamespace)
	profile := adapter.NewProfile(cl, profileNamespace)

	inlineResolver := adapter.NewInlineResolver()
	objectRefResolver := adapter.NewObjectRefResolver(dynCl)
	webhookResolver := adapter.NewWebhookResolver(objectRefResolver)

	butaneTransformer := adapter.NewButaneTransformer()
	webhookTransformer := adapter.NewWebhookTransformer(objectRefResolver)

	// --------------------------------------------- Controller ----------------------------------------------------- //
	var baseURL string

	mux := controller.NewResolveTransformerMux(
		baseURL,
		map[types.ResolverKind]adapter.Resolver{
			types.InlineResolverKind:    inlineResolver,
			types.ObjectRefResolverKind: objectRefResolver,
			types.WebhookResolverKind:   webhookResolver,
		},
		map[types.TransformerKind]adapter.Transformer{
			types.ButaneTransformerKind:  butaneTransformer,
			types.WebhookTransformerKind: webhookTransformer,
		},
	)

	ipxe := controller.NewIPXE(assignment, profile, mux)
	config := controller.NewConfig(profile, mux)

	// --------------------------------------------- App ------------------------------------------------------------ //

	handler := server.New(ipxe, config)

	app := echo.New()
	app.Use(echoprometheus.NewMiddleware(Name))
	server.RegisterHandlers(app, handler)

	// --------------------------------------------- Metrics -------------------------------------------------------- //

	metrics := echo.New()
	metrics.GET(metricsPath, echoprometheus.NewHandler())

	// --------------------------------------------- Probes --------------------------------------------------------- //

	// TODO: create func initializing a probe server which returns non-200 response when server is considered Unhealthy.

	probes := echo.New()
	probes.GET(probesHealthPath, func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "health probe error")
	})
	probes.GET(probesReadyPath, func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "readiness probe error")
	})

	// --------------------------------------------- Run Server ----------------------------------------------------- //

	if err := cmd.Serve(map[int]*echo.Echo{
		appServerPort:     app,
		metricsServerPort: metrics,
		probesServerPort:  probes,
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
