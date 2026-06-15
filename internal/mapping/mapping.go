package mapping

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type rawMapping struct {
	Prefixes  map[string]string           `yaml:"prefixes"`
	Functions map[string]string           `yaml:"functions"`
	Sources   map[string]yaml.RawMessage  `yaml:"sources"`
	Mappings  map[string]rawMappingConfig `yaml:"mappings"`
}

type Mapping struct {
	Prefixes  map[string]string
	Functions map[string]string
	Sources   map[string]SourceConfig
	Mappings  map[string]MappingConfig
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
	Source  string
	Subject string
	Triples []TripleMapping
}

type TripleMapping struct {
	Predicate      string         `yaml:"predicate"`
	Object         string         // can be null if mapping defines function call
	FunctionName   string         `yaml:"function"`
	FunctionParams map[string]any `yaml:"params"`
}

type rawMappingConfig struct {
	Source  string            `yaml:"source"`
	Subject string            `yaml:"subject"`
	Triples []yaml.RawMessage `yaml:"triples"`
}

func ReadMapping(mappingfile string) (Mapping, error) {
	m := Mapping{Sources: map[string]SourceConfig{}, Mappings: map[string]MappingConfig{}}
	var rawMapping rawMapping

	yamlBytes, fileErr := os.ReadFile(mappingfile)
	if fileErr != nil {
		return Mapping{}, fileErr
	}

	if err := yaml.Unmarshal(yamlBytes, &rawMapping); err != nil {
		return Mapping{}, err
	}

	m.Prefixes = rawMapping.Prefixes
	m.Functions = rawMapping.Functions

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

	for mappingName, rawMappingConfig := range rawMapping.Mappings {
		mc := MappingConfig{Source: rawMappingConfig.Source, Subject: rawMappingConfig.Subject, Triples: []TripleMapping{}}

		var tripleAsArray []string

		for idx, tripleConf := range rawMappingConfig.Triples {
			if err := yaml.Unmarshal(tripleConf, &tripleAsArray); err == nil {
				// if we land here, parsing yaml to array was successfull, so we assume the triple mapping was defined via the simple array notation
				predicate := tripleAsArray[0]
				object := tripleAsArray[1]

				if len(strings.TrimSpace(predicate)) == 0 {
					return Mapping{}, fmt.Errorf("Unable to parse mapping definition for mapping '%s' on row %d: predicate is empty", mappingName, idx)
				}

				if len(strings.TrimSpace(object)) == 0 {
					return Mapping{}, fmt.Errorf("Unable to parse mapping definition for mapping '%s' on row %d: object is empty", mappingName, idx)
				}

				mc.Triples = append(mc.Triples, TripleMapping{Predicate: predicate, Object: object})
			} else {
				// if we land here, we assume that the tripe mapping was defined in the object notation to use a function call

				var tripleAsObject TripleMapping
				if err := yaml.Unmarshal(tripleConf, &tripleAsObject); err != nil {
					return Mapping{}, fmt.Errorf("Unable to parse mapping definition for mapping '%s' on row %d: %w", mappingName, idx, err)
				}

				mc.Triples = append(mc.Triples, tripleAsObject) // TODO validate triple mapping
			}
		}

		m.Mappings[mappingName] = mc
	}

	return m, nil
}

func MergeMappings(mappings []Mapping) Mapping {
	mergedMapping := Mapping{Prefixes: map[string]string{}, Functions: map[string]string{}, Sources: map[string]SourceConfig{}, Mappings: map[string]MappingConfig{}}

	for _, m := range mappings {
		maps.Copy(mergedMapping.Prefixes, m.Prefixes)
		maps.Copy(mergedMapping.Functions, m.Functions)
		maps.Copy(mergedMapping.Sources, m.Sources)
		maps.Copy(mergedMapping.Mappings, m.Mappings)
	}

	return mergedMapping
}
