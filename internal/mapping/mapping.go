package mapping

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Mapping struct {
	Prefixes map[string]string        `yaml:"prefixes"`
	Sources  map[string]SourceConfig  `yaml:"sources"`
	Mappings map[string]MappingConfig `yaml:"mappings"`
}

type SourceConfig struct {
	Type string `yaml:"type"`
}

type MappingConfig struct {
	Source  string `yaml:"source"`
	Subject string `yaml:"subject"`
}

func ReadMapping(mappingfile string) (Mapping, error) {
	var m Mapping

	yamlBytes, fileErr := os.ReadFile(mappingfile)
	if fileErr != nil {
		return Mapping{}, fileErr
	}

	if err := yaml.Unmarshal(yamlBytes, &m); err != nil {
		return Mapping{}, err
	}

	return m, nil
}
