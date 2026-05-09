package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

	return nil
}

func checkIfFileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}

	return true
}
