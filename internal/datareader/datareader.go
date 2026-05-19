package datareader

import "github.com/MarkusFank/rdfmap2go/internal/mapping"

type DataReader interface {
	Init(sourceConfig mapping.SourceConfig) error
	Read() (<-chan RowResult, error)
	Close()
}

type RowResult struct {
	Row   DataRow
	Error error
}
type DataRow map[string]any
