package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/engine"
	"github.com/spf13/cobra"
)

var mappingFiles []string
var outputFile string
var outputType string
var allowedOutputTypes []string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate RDF using one or multiple mapping files",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !slices.Contains(allowedOutputTypes, strings.ToLower(outputType)) {
			return fmt.Errorf("Output type %q is not supported!", outputType)
		}

		return engine.Run(mappingFiles, outputFile, outputType)
	},
}

func init() {
	allowedOutputTypes = []string{"ttl", "nt"}
	generateCmd.Flags().StringArrayVarP(&mappingFiles, "mapping", "m", []string{}, "Mapping file (YAML)")
	generateCmd.Flags().StringVarP(&outputFile, "out", "o", "output.nt", "Output file")
	generateCmd.Flags().StringVarP(&outputType, "type", "t", allowedOutputTypes[0], fmt.Sprintf("Type of the output file (%s)", strings.Join(allowedOutputTypes, ",")))

	generateCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(generateCmd)
}
