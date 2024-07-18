package main

import (
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/hack/internal"
	"os"
	"os/exec"
)

// ----------------------------------------------------- TEARDOWN --------------------------------------------------- //

func teardown() error {
	// 1. read project setupConfig.
	projectCfg, err := internal.ReadProjectConfig()
	if err != nil {
		return err // TODO: wrap err
	}

	_, _ = fmt.Fprintf(os.Stdout, "Tearing down kindenv %q\n", projectCfg.Name)

	// 2. read kindenv setupConfig
	cfg, err := readSetupConfig()
	if err != nil {
		return fmt.Errorf("%s\nERROR: %w", formatSetupUsage(), err) // TODO: wrap err
	}

	_, _ = fmt.Fprintf(os.Stdout, "%#v\n", cfg)

	// 3. Do
	if err := doTeardown(projectCfg, cfg); err != nil {
		return err // TODO: wrap error
	}

	return nil
}

func doTeardown(pCfg internal.ProjectConfig, cfg setupConfig) error {
	// 1. kind create cluster and wait.
	cmd := exec.Command(
		cfg.KindBinary,
		"delete",
		"cluster",
		"--name", pCfg.Name,
	)

	if err := internal.RunCmdWithStdPipes(cmd); err != nil {
		return err // TODO: wrap error
	}

	if err := os.Remove(pCfg.Kindenv.KubeconfigPath); err != nil {
		return err // TODO: wrap error
	}

	return nil
}
