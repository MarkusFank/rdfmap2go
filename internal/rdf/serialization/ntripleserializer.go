package serialization

import (
	"os"

	"github.com/MarkusFank/rdfmap2go/internal/rdf"
	"github.com/tggo/goRDFlib/nt"
)

type NTripleSerializer struct {
}

func (serializer *NTripleSerializer) Serialize(tripleStore *rdf.TripleStore, outputFile string) error {

	file, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer file.Close()

	return nt.Serialize(tripleStore.Graph, file)
}
