package internal

import (
	"errors"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	projectConfigPath = ".project.yaml"
)

// ----------------------------------------------------- PROJECT CONFIG --------------------------------------------- //

type ProjectConfig struct {
	Name string `json:"name"`

	Kindenv struct {
		KubeconfigPath string `json:"kubeconfigPath"`
	} `json:"kindenv"`

	LocalContainerRegistry struct {
		Enabled        bool   `json:"enabled"`
		CredentialPath string `json:"credentialPath"`
		CaCrtPath      string `json:"caCrtPath"`
		Namespace      string `json:"namespace"`
	} `json:"localContainerRegistry"`
}

func ReadProjectConfig() (ProjectConfig, error) {
	b, err := os.ReadFile(projectConfigPath) //nolint:varnamelen
	if err != nil {
		return ProjectConfig{}, err // TODO: wrap err
	}

	out := ProjectConfig{} //nolint:exhaustruct // unmarshal

	if err := yaml.Unmarshal(b, &out); err != nil {
		return ProjectConfig{}, err // TODO: wrap err
	}

	err = nil // ensures err is nil

	for key, rule := range map[string]bool{
		"name":                   out.Name == "",
		"kindenv.kubeconfigPath": out.Kindenv.KubeconfigPath == "",
	} {
		if rule {
			err = errors.Join(err, fmt.Errorf("%s must be specified in %s", key, projectConfigPath))
		}
	}

	if err != nil {
		return ProjectConfig{}, err // TODO: wrap error.
	}

	return out, nil
}
