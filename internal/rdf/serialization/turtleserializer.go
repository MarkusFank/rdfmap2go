package serialization

import (
	"os"

	"github.com/MarkusFank/rdfmap2go/internal/rdf"
	"github.com/tggo/goRDFlib/turtle"
)

type TurtleSerializer struct {
}

func (serializer *TurtleSerializer) Serialize(tripleStore *rdf.TripleStore, outputFile string) error {
	file, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer file.Close()

	return turtle.Serialize(tripleStore.Graph, file)
}
