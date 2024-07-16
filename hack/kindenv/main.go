package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	projectConfigPath = ".project.yaml"

	// Available commands
	setupCommand    = "setup"
	teardownCommand = "teardown"
	helpCommand     = "usage"
)

// ----------------------------------------------------- USAGE ------------------------------------------------------ //

const (
	banner        = "# KINDENV\n\n"
	usageTemplate = `## Usage

%s [command]

Available commands:
  - %q
  - %q
`
)

func usage() error {
	arg0 := fmt.Sprintf("go run \"%s/hack/kindenv\"", os.Getenv("PWD"))
	_, _ = fmt.Fprintf(os.Stdout, usageTemplate, arg0, setupCommand, teardownCommand)

	return nil
}

// ----------------------------------------------------- CONFIG ----------------------------------------------------- //

type projectConfig struct {
	Name string `json:"name"`

	Kindenv struct {
		KubeconfigPath string `json:"kubeconfigPath"`
	} `json:"kindenv"`
}

// ----------------------------------------------------- MAIN ------------------------------------------------------- //

func main() {
	_, _ = fmt.Fprint(os.Stdout, banner)

	// 1. Print usageTemplate or

	if len(os.Args) < 2 { //nolint:gomnd // if no specified subcommand then print usageTemplate and exit.
		_ = usage()

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
		command = usage
	}

	// 3. Execute command

	if err := command(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
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

	err = nil
	for key, rule := range map[string]bool{
		"name":                   out.Name == "",
		"kindenv.kubeconfigPath": out.Kindenv.KubeconfigPath == "",
	} {
		if rule {
			err = errors.Join(err, fmt.Errorf("%s must be specified in %s", key, projectConfigPath))
		}
	}
	if err != nil {
		return projectConfig{}, err // TODO: wrap error.
	}

	return out, nil
}

func formatExpectedEnvList[T any]() string {
	optionalEnvs := make([]string, 0)
	requiredEnvs := make([]string, 0)

	observedMaxStrLen := 0

	rt := reflect.TypeFor[T]()
	for i := range rt.NumField() {
		field := rt.Field(i)
		val, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}

		substr := strings.Split(val, ",")
		switch len(substr) {
		case 0:
			continue
		case 1:
			optionalEnvs = append(optionalEnvs, substr[0])
		default:
			requiredEnvs = append(requiredEnvs, substr[0])
		}

		if envStrLen := len(substr[0]); envStrLen > observedMaxStrLen {
			observedMaxStrLen = envStrLen
		}
	}

	envs := ""
	for _, s := range requiredEnvs {
		envs = fmt.Sprintf("%s- %s %s[Required]\n", envs, s, fmtSpaces(s, observedMaxStrLen))
	}

	for _, s := range optionalEnvs {
		envs = fmt.Sprintf("%s- %s %s[Optional]\n", envs, s, fmtSpaces(s, observedMaxStrLen))
	}

	return envs
}

func fmtSpaces(s string, maxLen int) string {
	return strings.Repeat(" ", maxLen-len(s))
}

func runCmdWithStdPipes(cmd *exec.Cmd) error {
	errChan := make(chan error)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			errChan <- err
		}
	}()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		if written, err := io.Copy(os.Stderr, stderr); err != nil {
			errChan <- err

			if written > 0 {
				errChan <- fmt.Errorf("%d bytes written to stderr", written) // TODO: wrap err
			}
		}
	}()

	if err := cmd.Run(); err != nil {
		return err
	}

	close(errChan)
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}
