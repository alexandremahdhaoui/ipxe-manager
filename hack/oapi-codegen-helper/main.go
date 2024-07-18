package main

import (
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/hack/internal"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	OAPICodegenEnvKey = "OAPI_CODEGEN"

	errEnv = "OAPI_CODEGEN env var must be set to the path to OAPI_CODEGEN or an executable `go run <PACKAGE>` command"

	configName = ".oapi-codegen"
	configPath = "."

	sourceFileTemplate  = "%s.%s.yaml"
	zzGeneratedFilename = "zz_generated.oapi-codegen.go"

	clientTemplate = `---
package: %[1]s
output: %[2]s
generate:
  client: true
  models: true
  embedded-spec: true
output-options:
  # to make sure that all types are generated
  skip-prune: true
`

	serverTemplate = `---
package: %[1]s
output: %[2]s
generate:
  embedded-spec: true
  models: true
  std-http-server: true
  strict-server: true
output-options:
  skip-prune: true
`
)

func main() {
	executable := os.Getenv(OAPICodegenEnvKey)
	if executable == "" {
		_, _ = fmt.Fprintln(os.Stderr, errEnv)
		os.Exit(1)
	}

	cfg, err := readConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := do(executable, cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, _ = fmt.Fprintln(os.Stdout, "successfully generated code")
	os.Exit(0)
}

type (
	GenOpts struct {
		Enabled     bool   `json:"enabled"`
		PackageName string `json:"packageName"`
	}

	Spec struct {
		Name     string   `json:"name"`
		Versions []string `json:"versions"`

		Client GenOpts `json:"client"`
		Server GenOpts `json:"server"`

		Source         string `json:"source,omitempty"`
		DestinationDir string `json:"destinationDir,omitempty"`
	}

	Config struct {
		Specs []Spec `json:"specs"`

		Defaults struct {
			SourceDir      string `json:"sourceDir"`
			DestinationDir string `json:"destinationDir"`
		} `json:"defaults"`
	}
)

// ---

func do(executable string, config Config) error {
	cmdName, args := parseExecutable(executable)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}

	for i := range config.Specs { // for each spec
		i := i
		for _, version := range config.Specs[i].Versions { // for each version
			version := version
			wg.Add(1)

			// for each spec and each version in that spec:

			sourcePath := templateSourcePath(config, i, version)

			for _, pkg := range []struct { // for each client OR server pkg
				opts     GenOpts
				template string
			}{
				{ // Client
					opts:     config.Specs[i].Client,
					template: clientTemplate,
				},
				{ // Server
					opts:     config.Specs[i].Server,
					template: serverTemplate,
				},
			} {
				go func() {
					defer wg.Done()
					if !pkg.opts.Enabled {
						return
					}

					outputPath := templateOutputPath(config, i, pkg.opts.PackageName)
					templatedConfig := fmt.Sprintf(pkg.template, pkg.opts.PackageName, outputPath)

					path, cleanup, err := writeTempCodegenConfig(templatedConfig)
					if err != nil {
						errChan <- err // TODO: wrap err
					}

					defer cleanup()

					if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
						errChan <- err // TODO: wrap err
					}

					args := append(args, "--config", path, sourcePath)
					if err := internal.RunCmdWithStdPipes(exec.Command(cmdName, args...)); err != nil {
						errChan <- err // TODO: wrap err
					}
				}()
			}

		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// if any error occur or the channel is closed, we return the first error early
	if err := <-errChan; err != nil {
		return err // TODO: wrap error
	}

	return nil
}

func parseExecutable(executable string) (string, []string) {
	split := strings.Split(executable, " ")

	return split[0], split[1:]
}

// ---

// writeTempCodegenConfig return the path to the generated config file, a cleanup function and an error.
func writeTempCodegenConfig(templatedConfig string) (string, func(), error) {
	// 1. create tempfile
	tempFile, err := os.CreateTemp("", "oapi-codegen-*.yaml")
	if err != nil {
		return "", nil, err // TODO: wrap err
	}

	// 2. create a cleanup func
	cleanup := func() {
		os.RemoveAll(tempFile.Name())
	}

	// 3. write to file.
	if _, err := tempFile.WriteString(templatedConfig); err != nil {
		cleanup()

		return "", nil, err // TODO: wrap err
	}

	// 4. close file
	if err := tempFile.Close(); err != nil {
		cleanup()

		return "", nil, err // TODO: wrap err
	}

	return tempFile.Name(), cleanup, nil
}

func templateOutputPath(config Config, index int, packageName string) string {
	destDir := config.Defaults.DestinationDir
	if config.Specs[index].DestinationDir != "" { // it takes precedence over defaults.
		destDir = config.Specs[index].DestinationDir
	}

	return filepath.Join(destDir, packageName, zzGeneratedFilename)
}

func templateSourcePath(config Config, index int, version string) string {
	if source := config.Specs[index].Source; source != "" {
		return source
	}

	sourceFile := fmt.Sprintf(sourceFileTemplate, config.Specs[index].Name, version)

	return filepath.Join(config.Defaults.SourceDir, sourceFile)
}

// ---

func readConfig() (Config, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err // TODO: wrap err
	}

	out := Config{} //nolint:exhaustruct // unmarshal

	if err := viper.Unmarshal(&out); err != nil {
		return Config{}, err // TODO: wrap err
	}

	return out, nil
}
