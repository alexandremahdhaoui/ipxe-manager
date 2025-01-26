package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexandremahdhaoui/ipxer/internal/util/httputil"
	ipxerv1alpha1 "github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/alexandremahdhaoui/ipxer/pkg/generated/ipxerserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/driver/server"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/gracefulshutdown"
)

const (
	Name             = "ipxer-api"
	ConfigPathEnvKey = "IPXER_CONFIG_PATH"

	KubeconfigFromServiceAccount = ">>> Kubeconfig From Service Account"
)

var (
	Version        = "dev" //nolint:gochecknoglobals // set by ldflags
	CommitSHA      = "n/a" //nolint:gochecknoglobals // set by ldflags
	BuildTimestamp = "n/a" //nolint:gochecknoglobals // set by ldflags
)

type Config struct {
	// Adapters

	AssignmentNamespace string `json:"assignmentNamespace"`
	ProfileNamespace    string `json:"profileNamespace"`

	// Kubeconfig

	KubeconfigPath string `json:"kubeconfigPath"`

	// ProbesServer
	ProbesServer struct {
		LivenessPath  string `json:"livenessPath"`
		ReadinessPath string `json:"readinessPath"`
		Port          int    `json:"port"`
	} `json:"probesServer"`

	// MetricsServer
	MetricsServer struct {
		Path string `json:"path"`
		Port int    `json:"port"`
	} `json:"metricsServer"`

	// APIServer
	APIServer struct {
		Port int `json:"port"`
	} `json:"apiServer"`
}

// ------------------------------------------------- Main ----------------------------------------------------------- //

func main() {
	_, _ = fmt.Fprintf(
		os.Stdout,
		"Starting %s version %s (%s) %s\n",
		Name,
		Version,
		CommitSHA,
		BuildTimestamp,
	)

	gs := gracefulshutdown.New(Name)
	ctx := gs.Context()

	// --------------------------------------------- Config --------------------------------------------------------- //

	ipxerConfigPath := os.Getenv(ConfigPathEnvKey)
	if ipxerConfigPath == "" {
		slog.ErrorContext(ctx, fmt.Sprintf("environment variable %q must be set", ConfigPathEnvKey))
		gs.Shutdown(1)
	}

	b, err := os.ReadFile(ipxerConfigPath)
	if err != nil {
		slog.ErrorContext(ctx, "reading ipxer-api configuration file", "error", err.Error())
		gs.Shutdown(1)
	}

	config := new(Config)
	if err = json.Unmarshal(b, config); err != nil {
		slog.ErrorContext(ctx, "parsing ipxer-api configuration", "error", err.Error())
		gs.Shutdown(1)
	}

	// --------------------------------------------- Client --------------------------------------------------------- //

	restConfig, err := newKubeRestConfig(config.KubeconfigPath)
	if err != nil {
		slog.ErrorContext(ctx, "creating kube rest config", "error", err.Error())
		gs.Shutdown(1)
	}

	cl, err := newKubeClient(restConfig)
	if err != nil {
		slog.ErrorContext(ctx, "creating kube client", "error", err.Error())
		gs.Shutdown(1)
	}

	dynCl, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		slog.ErrorContext(ctx, "creating dynamic client", "error", err.Error())
		gs.Shutdown(1)
	}

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

	ipxerHandler := ipxerserver.Handler(ipxerserver.NewStrictHandler(
		server.New(ipxe, content),
		nil, // TODO: prometheus middleware
	))

	ipxerServer := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf(":%d", config.APIServer.Port),
		Handler:           ipxerHandler,
		ReadHeaderTimeout: time.Second,
		// TODO: set fields etc...
	}

	// --------------------------------------------- Metrics -------------------------------------------------------- //

	metricsHandler := http.NewServeMux()
	metricsHandler.Handle(config.MetricsServer.Path, promhttp.Handler())

	metrics := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf(":%d", config.MetricsServer.Port),
		Handler:           metricsHandler,
		ReadHeaderTimeout: time.Second,
	}

	// --------------------------------------------- Probes --------------------------------------------------------- //

	probesHandler := http.NewServeMux()

	probesHandler.Handle(config.ProbesServer.LivenessPath, http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	probesHandler.Handle(config.ProbesServer.ReadinessPath, http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	probes := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf(":%d", config.ProbesServer.Port),
		Handler:           probesHandler,
		ReadHeaderTimeout: time.Second,
	}

	// --------------------------------------------- Run Server ----------------------------------------------------- //

	httputil.Serve(map[string]*http.Server{
		"ipxer":   ipxerServer,
		"metrics": metrics,
		"probes":  probes,
	}, gs)

	slog.Info("✅ gracefully stopped", "binary", Name)
}

// ------------------------------------------------- Helpers -------------------------------------------------------- //

func newKubeRestConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == KubeconfigFromServiceAccount {
		return rest.InClusterConfig() // TODO: wrap err
	}

	b, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(b)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	return restConfig, nil
}

func newKubeClient(restConfig *rest.Config) (client.Client, error) { //nolint:ireturn
	sch := runtime.NewScheme()

	if err := corev1.AddToScheme(sch); err != nil {
		return nil, err // TODO: wrap err
	}

	if err := ipxerv1alpha1.AddToScheme(sch); err != nil {
		return nil, err // TODO: wrap err
	}

	cl, err := client.New(restConfig, client.Options{Scheme: sch}) //nolint:exhaustruct
	if err != nil {
		return nil, err // TODO: wrap err
	}

	return cl, nil
}
