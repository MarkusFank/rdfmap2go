package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

func Run(mappingFiles []string, outputFile string) error {

	if len(mappingFiles) == 0 {
		return errors.New("At least one mapping file has to be specified")
	}

	var allfiles strings.Builder
	for _, f := range mappingFiles {

		if fileExists := checkIfFileExists(f); !fileExists {
			return fmt.Errorf("The specified mapping file '%s' does not exist", f)
		}
		allfiles.WriteString("'" + f + "' ")
	}

	fmt.Println("Generating RDF file using " + allfiles.String())

	mappings, parseError := parseMappings(mappingFiles)

	if parseError != nil {
		return parseError
	}

	fmt.Printf("Successfully parsed %d mapping(s)\n", len(mappings))

	// TODO do a mapping validation after the mappings are parsed and before actual processing begins

	mapping := mapping.MergeMappings(mappings)

	processingErr := Process(&mapping, outputFile)

	if processingErr != nil {
		return processingErr
	}

	return nil
}

func checkIfFileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}

	return true
}

func parseMappings(mappingFiles []string) ([]mapping.Mapping, error) {
	var mappings []mapping.Mapping

	for _, f := range mappingFiles {
		m, err := mapping.ReadMapping(f)

		if err != nil {
			return nil, errors.Join(fmt.Errorf("Error while parsing mapping file '%s'", f), err)
		}

		mappings = append(mappings, m)
	}

	return mappings, nil
}
