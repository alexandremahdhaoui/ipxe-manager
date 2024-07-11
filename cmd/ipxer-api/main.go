package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/cmd"
	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/driver/server"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/gracefulshutdown"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Name             = "ipxer-api"
	ConfigPathEnvKey = "IPXER_CONFIG_PATH"
)

var (
	Version        = "dev" //nolint:gochecknoglobals // set by ldflags
	CommitSHA      = "n/a" //nolint:gochecknoglobals // set by ldflags
	BuildTimestamp = "n/a" //nolint:gochecknoglobals // set by ldflags
)

type Config struct {
	KubeconfigPath string `json:"kubeconfigPath"`

	// adapters

	AssignmentNamespace string `json:"assignmentNamespace"`
	ProfileNamespace    string `json:"profileNamespace"`

	// APIServer
	APIServer struct {
		Port int `json:"port"`
	} `json:"apiServer"`

	// MetricsServer
	MetricsServer struct {
		Port int    `json:"port"`
		Path string `json:"path"`
	} `json:"metricsServer"`

	// ProbesServer
	ProbesServer struct {
		Port          int    `json:"port"`
		LivenessPath  string `json:"livenessPath"`
		ReadinessPath string `json:"readinessPath"`
	}
}

// ------------------------------------------------- Main ----------------------------------------------------------- //

func main() {
	fmt.Printf("Starting %s version %s (%s) %s\n", Name, Version, CommitSHA, BuildTimestamp)

	gs := gracefulshutdown.New(Name)
	ctx := gs.Context()

	// --------------------------------------------- Config --------------------------------------------------------- //

	ipxerConfigPath := os.Getenv(ConfigPathEnvKey)
	if ipxerConfigPath == "" {
		slog.ErrorContext(ctx, fmt.Sprintf("environment variable %q must be set", ConfigPathEnvKey))
		gs.Shutdown(1)
		return
	}

	b, err := os.ReadFile(ipxerConfigPath)
	if err != nil {
		slog.ErrorContext(ctx, "reading ipxer-api configuration file", "error", err.Error())
		gs.Shutdown(1)
		return
	}

	config := new(Config)
	if err := json.Unmarshal(b, config); err != nil {
		slog.ErrorContext(ctx, "parsing ipxer-api configuration", "error", err.Error())
		gs.Shutdown(1)
		return
	}

	// --------------------------------------------- Client --------------------------------------------------------- //

	var cl client.Client        // TODO
	var dynCl dynamic.Interface // TODO

	// --------------------------------------------- Adapter -------------------------------------------------------- //

	assignment := adapter.NewAssignment(cl, config.AssignmentNamespace)
	profile := adapter.NewProfile(cl, config.ProfileNamespace)

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
	content := controller.NewContent(profile, mux)

	// --------------------------------------------- App ------------------------------------------------------------ //

	handler := server.New(ipxe, content)

	app := echo.New()
	app.Use(echoprometheus.NewMiddleware(Name))
	server.RegisterHandlers(app, handler)

	// --------------------------------------------- Metrics -------------------------------------------------------- //

	metrics := echo.New()
	metrics.GET(config.MetricsServer.Path, echoprometheus.NewHandler())

	// --------------------------------------------- Probes --------------------------------------------------------- //

	// TODO: create func initializing a probe server which returns non-200 response when server is considered Unhealthy.

	probes := echo.New()
	probes.GET(config.ProbesServer.LivenessPath, func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "liveness probe error")
	})
	probes.GET(config.ProbesServer.ReadinessPath, func(c echo.Context) error {
		return wrapErr(c.NoContent(http.StatusOK), "readiness probe error")
	})

	// --------------------------------------------- Run Server ----------------------------------------------------- //

	servers := map[int]*echo.Echo{
		config.APIServer.Port:     app,
		config.MetricsServer.Port: metrics,
		config.ProbesServer.Port:  probes,
	}

	cmd.Serve(servers, gs)

	app.Logger.Infof("✅ successfully stopped %s", Name)
}

// --------------------------------------------- UTILS -------------------------------------------------------------- //

func wrapErr(err error, s string) error {
	if err != nil {
		return errors.Join(err, errors.New(s)) //nolint: goerr113
	}

	return nil
}
