package main

import (
	"log/slog"

	"github.com/MarkusFank/rdfmap2go/internal/cli"
)

func main() {

	slog.SetDefault(slog.Default())

	cli.Execute()
}
