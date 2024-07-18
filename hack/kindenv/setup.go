package main

import (
	"errors"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/hack/internal"
	"os"
	"os/exec"

	"github.com/caarlos0/env/v11"
)

// ----------------------------------------------------- USAGE ------------------------------------------------------ //

const (
	//nolint:dupword
	setupUsageTemplate = `
## Setup

The setup command may expect the following env variables:
%s`
)

func formatSetupUsage() string {
	return fmt.Sprintf(setupUsageTemplate, formatExpectedEnvList[setupConfig]())
}

// ----------------------------------------------------- CONFIG ----------------------------------------------------- //

type setupConfig struct {
	KindBinary string `env:"KIND_BINARY,required"`

	// TODO: make use of the below variables.
	ContainerRegistryBaseURL string `env:"CONTAINER_REGISTRY_BASE_URL"`
	ContainerEngineBinary    string `env:"CONTAINER_ENGINE_BINARY"`
	HelmBinary               string `env:"HELM_BINARY"`
}

func readSetupConfig() (setupConfig, error) {
	out := setupConfig{} //nolint:exhaustruct // unmarshal

	if err := env.Parse(&out); err != nil {
		return setupConfig{}, err // TODO: wrap err
	}

	return out, nil
}

// ----------------------------------------------------- SETUP ------------------------------------------------------ //

func setup() error {
	// 1. read project setupConfig.
	projectCfg, err := internal.ReadProjectConfig()
	if err != nil {
		return err // TODO: wrap err
	}

	_, _ = fmt.Fprintf(os.Stdout, "Setting up kindenv %q\n", projectCfg.Name)

	// 2. read kindenv setupConfig
	cfg, err := readSetupConfig()
	if err != nil {
		return fmt.Errorf("%s\nERROR: %w", formatSetupUsage(), err) // TODO: wrap err
	}

	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", cfg)

	// 3. Do
	if err := doSetup(projectCfg, cfg); err != nil {
		return errors.Join(err, doTeardown(projectCfg, cfg))
	}

	return nil
}

func doSetup(pCfg internal.ProjectConfig, cfg setupConfig) error {
	// 1. kind create cluster and wait.
	cmd := exec.Command(
		cfg.KindBinary,
		"create",
		"cluster",
		"--name", pCfg.Name,
		"--kubeconfig", pCfg.Kindenv.KubeconfigPath,
		"--wait", "5m",
	)

	if err := internal.RunCmdWithStdPipes(cmd); err != nil {
		return err // TODO: wrap error
	}

	// 2. TODO: setup communication towards local-registry.

	// 3. TODO: setup communication towards any provided registry (e.g. required if users wants to install some apps into their kind cluster). It can be any OCI registry. (to support helm chart)

	// 4. TODO: setup communication CONTAINER_ENGINE login & HELM login.

	return nil
}
