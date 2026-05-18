package mapping

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type rawMapping struct {
	Prefixes map[string]string          `yaml:"prefixes"`
	Sources  map[string]yaml.RawMessage `yaml:"sources"`
	Mappings map[string]MappingConfig   `yaml:"mappings"`
}

type Mapping struct {
	Prefixes map[string]string
	Sources  map[string]SourceConfig
	Mappings map[string]MappingConfig
}

type sourceConfigRawBase struct {
	Type string `yaml:"type"`
}

type SourceConfig interface {
	GetSourceType() string
}

type CsvSourceConfig struct {
	Type string `yaml:"type"`
	File string `yaml:"file"`
}

func (s CsvSourceConfig) GetSourceType() string {
	return "csv"
}

type SqliteSourceConfig struct {
	Type  string `yaml:"type"`
	File  string `yaml:"file"`
	Query string `yaml:"query"`
}

func (s SqliteSourceConfig) GetSourceType() string {
	return "sqlite"
}

type JsonSourceConfig struct {
	Type     string `yaml:"type"`
	File     string `yaml:"file"`
	JsonPath string `yaml:"jsonPath"`
}

func (s JsonSourceConfig) GetSourceType() string {
	return "json"
}

type MappingConfig struct {
	Source  string     `yaml:"source"`
	Subject string     `yaml:"subject"`
	Triples [][]string `yaml:"triples"`
}

func ReadMapping(mappingfile string) (Mapping, error) {
	m := Mapping{Sources: map[string]SourceConfig{}}
	var rawMapping rawMapping

	yamlBytes, fileErr := os.ReadFile(mappingfile)
	if fileErr != nil {
		return Mapping{}, fileErr
	}

	if err := yaml.Unmarshal(yamlBytes, &rawMapping); err != nil {
		return Mapping{}, err
	}

	m.Prefixes = rawMapping.Prefixes
	m.Mappings = rawMapping.Mappings

	for sourceName, rawSource := range rawMapping.Sources {
		var baseSource sourceConfigRawBase
		if err := yaml.Unmarshal(rawSource, &baseSource); err != nil {
			return Mapping{}, err
		}

		switch strings.ToLower(baseSource.Type) {
		case "csv":
			var csvSource CsvSourceConfig
			if err := yaml.Unmarshal(rawSource, &csvSource); err != nil {
				return Mapping{}, err
			}

			m.Sources[sourceName] = csvSource
		case "sqlite":
			var sqliteSource SqliteSourceConfig
			if err := yaml.Unmarshal(rawSource, &sqliteSource); err != nil {
				return Mapping{}, err
			}

			m.Sources[sourceName] = sqliteSource
		case "json":
			var jsonSource JsonSourceConfig
			if err := yaml.Unmarshal(rawSource, &jsonSource); err != nil {
				return Mapping{}, err
			}

			m.Sources[sourceName] = jsonSource
		default:
			return Mapping{}, fmt.Errorf("Unable to parse source '%s'", sourceName)
		}
	}

	return m, nil
}

func MergeMappings(mappings []Mapping) Mapping {
	mergedMapping := Mapping{Prefixes: map[string]string{}, Sources: map[string]SourceConfig{}, Mappings: map[string]MappingConfig{}}

	for _, m := range mappings {
		maps.Copy(mergedMapping.Prefixes, m.Prefixes)
		maps.Copy(mergedMapping.Sources, m.Sources)
		maps.Copy(mergedMapping.Mappings, m.Mappings)
	}

	return mergedMapping
}
