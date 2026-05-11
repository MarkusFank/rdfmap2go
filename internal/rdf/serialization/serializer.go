package serialization

import "github.com/MarkusFank/rdfmap2go/internal/rdf"

type TripleStoreSerializer interface {
	Serialize(tripleStore *rdf.TripleStore, outputFile string)
}
