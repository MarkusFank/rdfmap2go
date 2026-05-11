package engine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/datareader/csv"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
	"github.com/MarkusFank/rdfmap2go/internal/rdf"
	"github.com/MarkusFank/rdfmap2go/internal/rdf/serialization"
)

func Process(mapping *mapping.Mapping, outputFile string) error {

	sourcesToMappings, err := mapSourcesToMappings(mapping)

	if err != nil {
		return err
	}

	tripleStore := rdf.TripleStore{}

	for source, mappings := range sourcesToMappings {
		if len(mappings) == 0 {
			fmt.Printf("Warning: Source '%s' is used in no mapping\n", source)
			continue
		}

		err := processSource(source, &mappings, mapping, &tripleStore)

		if err != nil {
			return err
		}
	}

	fmt.Printf("Created %d triples\n", len(tripleStore.Triples))

	serializer := serialization.NTripleSerializer{} // TODO let serializer type be set from options

	serializer.Serialize(&tripleStore, outputFile)

	return nil
}

func processSource(sourceName string, mappingsForSource *[]string, mainMapping *mapping.Mapping, tripleStore *rdf.TripleStore) error {

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
			processDataRowWithMapping(*row, &mapping, mainMapping.Prefixes, tripleStore)
		}

	}

	return nil
}

func processDataRowWithMapping(dataRow datareader.DataRow, mapping *mapping.MappingConfig, prefixes map[string]string, tripleStore *rdf.TripleStore) {

	subject := expandPrefix(expandDataColumns(mapping.Subject, dataRow), prefixes)
	fmt.Printf("The subject is: %s\n", subject)

	for _, tripleConfig := range mapping.Triples {
		if len(tripleConfig) == 2 {
			predicateConf := tripleConfig[0]
			objectConf := tripleConfig[1]

			predicate := expandPrefix(expandDataColumns(predicateConf, dataRow), prefixes)
			object := expandPrefix(expandDataColumns(objectConf, dataRow), prefixes)

			fmt.Printf("\tpredicate %s; object: %s\n", predicate, object)

			fillToTripleStore(subject, predicate, object, tripleStore)

		} else {
			fmt.Printf("Warning: Unable to process triple %v\n", tripleConfig)
		}
	}
}

func fillToTripleStore(subject, predicate, object string, tripleStore *rdf.TripleStore) {
	subjectNode := createNodeForValue(subject)
	predicateNode := createNodeForValue(predicate)
	objectNode := createNodeForValue(object)

	tripleStore.AddTriple(subjectNode, predicateNode, objectNode)
}

func createNodeForValue(value string) rdf.Node {
	node := rdf.Node{}
	if strings.HasPrefix(value, "http://") {
		node.Type = rdf.URI
	} else {
		node.Type = rdf.Literal
	}

	node.Value = value

	return node
}

func expandDataColumns(templateString string, dataRow datareader.DataRow) string {
	valRegex := regexp.MustCompile(`\$\{([a-zA-Z0-9_]+)\}`)
	result := valRegex.ReplaceAllStringFunc(templateString, func(m string) string {
		sub := valRegex.FindStringSubmatch(m)
		val := dataRow[sub[1]].(string) // TODO handle other types than string
		return val
	})
	return result
}

func expandPrefix(value string, prefixes map[string]string) string {
	for prefix, replacement := range prefixes {
		prefixWithColon := prefix + ":"
		if after, ok := strings.CutPrefix(value, prefixWithColon); ok {
			return replacement + after
		}
	}

	return value
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
