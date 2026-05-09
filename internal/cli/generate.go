package cli

import (
	"github.com/MarkusFank/rdfmap2go/internal/engine"
	"github.com/spf13/cobra"
)

var mappingFiles []string
var outputFile string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate RDF using one or multiple mapping files",
	RunE: func(cmd *cobra.Command, args []string) error {

		return engine.Run(mappingFiles, outputFile)
	},
}

func init() {
	generateCmd.Flags().StringArrayVarP(&mappingFiles, "mapping", "m", []string{}, "Mapping file (YAML)")
	generateCmd.Flags().StringVarP(&outputFile, "out", "o", "output.ttl", "Output file")

	generateCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(generateCmd)
}
