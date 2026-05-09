package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rdfmap2go",
	Short: "RDF mapping tool",
	Long: `A Go based RDF mapping tool, that allows to read data from different sources
	and create RDF triples`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
