package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
)

const (
	projectConfigPath      = ".project.yaml"
	kubeconfigPathTemplate = "%s/.testenv-kubeconifg.yaml"

	// Available commands
	setupCommand    = "setup"
	teardownCommand = "teardown"
	helpCommand     = "help"
)

// ----------------------------------------------------- CONFIG ----------------------------------------------------- //

type projectConfig struct {
	Name string `json:"name"`
}

type config struct {
	ContainerRegistryBaseURL string `env:"CONTAINER_REGISTRY_BASE_URL,required"`
	KindBinary               string `env:"KIND_BINARY,required"`
}

// ----------------------------------------------------- USAGE ------------------------------------------------------ //

const (
	banner = "# KINDENV\n\n"
	usage  = `## Usage

%s [command]

Available commands:
  - %q
  - %q
`
)

// ----------------------------------------------------- MAIN ------------------------------------------------------- //

func main() {
	_, _ = fmt.Fprintf(os.Stdout, banner)

	// 1. Print usage or

	if len(os.Args) < 2 { //nolint:gomnd // if no specified subcommand then print usage and exit.
		_ = help()
		os.Exit(1)
	}

	// 2. Switch command.

	var command func() error

	switch os.Args[1] {
	case setupCommand:
		command = setup
	case teardownCommand:
		command = teardown
	case helpCommand:
		command = help
	}

	// 3. Execute command

	if err := command(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

// ----------------------------------------------------- SETUP ------------------------------------------------------ //

func setup() error {
	// 1. read project config.
	projectCfg, err := readProjectConfig()
	if err != nil {
		return err // TODO: wrap err
	}

	_, _ = fmt.Fprintf(os.Stdout, "Setting up kindenv %q\n", projectCfg.Name)

	// 2. read kindenv config
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", cfg)

	// 3. do stuff

	return nil
}

// ----------------------------------------------------- TEARDOWN --------------------------------------------------- //

func teardown() error {
	return nil
}

// ----------------------------------------------------- HELPERS ---------------------------------------------------- //

func help() error {
	arg0 := fmt.Sprintf("go run \"%s/hack/kindenv\"", os.Getenv("PWD"))
	_, _ = fmt.Fprintf(os.Stdout, usage, arg0, setupCommand, teardownCommand)

	return nil
}

// ----------------------------------------------------- HELPERS ---------------------------------------------------- //

func readProjectConfig() (projectConfig, error) {
	b, err := os.ReadFile(projectConfigPath) //nolint:varnamelen
	if err != nil {
		return projectConfig{}, err // TODO: wrap err
	}

	out := projectConfig{} //nolint:exhaustruct // unmarshal

	if err := yaml.Unmarshal(b, &out); err != nil {
		return projectConfig{}, err // TODO: wrap err
	}

	return out, nil
}

func readConfig() (config, error) {
	out := config{}

	if err := env.Parse(&out); err != nil {
		return config{}, err // TODO: wrap err
	}

	return out, nil
}
