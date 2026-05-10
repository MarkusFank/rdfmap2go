package engine

import (
	"fmt"
	"regexp"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/datareader/csv"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

func Process(mapping *mapping.Mapping) error {

	sourcesToMappings, err := mapSourcesToMappings(mapping)

	if err != nil {
		return err
	}

	for source, mappings := range sourcesToMappings {
		if len(mappings) == 0 {
			fmt.Printf("Warning: Source '%s' is used in no mapping\n", source)
			continue
		}

		err := processSource(source, &mappings, mapping)

		if err != nil {
			return err
		}
	}

	return nil
}

func processSource(sourceName string, mappingsForSource *[]string, mainMapping *mapping.Mapping) error {

	sourceConfig := mainMapping.Sources[sourceName]

	reader, err := createDataReaderForSource(sourceName, sourceConfig)

	if err != nil {
		return err
	}

	defer reader.Close()

	mappings := []mapping.MappingConfig{}

	for _, m := range *mappingsForSource {
		mappings = append(mappings, mainMapping.Mappings[m])
	}

	for {
		row, err := reader.ReadRow()

		if err != nil {
			return err
		}

		if row == nil {
			break
		}

		for _, mapping := range mappings {
			processDataRowWithMapping(*row, &mapping)
		}

	}

	return nil
}

func processDataRowWithMapping(dataRow datareader.DataRow, mapping *mapping.MappingConfig) {
	re := regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)

	result := re.ReplaceAllStringFunc(mapping.Subject, func(m string) string {
		sub := re.FindStringSubmatch(m)
		val := dataRow[sub[1]].(string) // TODO handle other types than string
		return val
	})

	fmt.Printf("The subject is: %s\n", result)
}

func createDataReaderForSource(sourceName string, sourceConfig mapping.SourceConfig) (datareader.DataReader, error) {

	switch sourceConfig.GetSourceType() {
	case "csv":
		csvSourceConfig := sourceConfig.(mapping.CsvSourceConfig)
		csvReader := csv.CsvDataReader{}
		csvReader.Init(csvSourceConfig.File)

		return &csvReader, nil
	}

	return nil, fmt.Errorf("Unable to create data reader for source '%s'", sourceName)
}

func mapSourcesToMappings(mapping *mapping.Mapping) (map[string][]string, error) {

	sourcesToMappings := map[string][]string{}
	for sourceName := range mapping.Sources {
		sourcesToMappings[sourceName] = []string{}
	}

	for mappingName, m := range mapping.Mappings {
		arr, hasSource := sourcesToMappings[m.Source]

		if !hasSource {
			return nil, fmt.Errorf("Mapping '%s' refers to a source '%s' which is not defined!", mappingName, m.Source)
		}

		sourcesToMappings[m.Source] = append(arr, mappingName)
	}

	return sourcesToMappings, nil
}
