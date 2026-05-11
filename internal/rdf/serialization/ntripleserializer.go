package serialization

import (
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/rdf"
)

type NTripleSerializer struct {
}

func (serializer *NTripleSerializer) Serialize(tripleStore *rdf.TripleStore, outputFile string) {
	stringBuilder := strings.Builder{}

	for _, triple := range tripleStore.Triples {
		stringBuilder.WriteString(formatNode(triple.Subject))
		stringBuilder.WriteString(" ")
		stringBuilder.WriteString(formatNode(triple.Predicate))
		stringBuilder.WriteString(" ")
		stringBuilder.WriteString(formatNode(triple.Object))
		stringBuilder.WriteString("\n")
	}

	os.WriteFile(outputFile, []byte(stringBuilder.String()), os.ModeAppend)
}

func formatNode(n rdf.Node) string {

	switch n.Type {

	case rdf.URI:
		return "<" + n.Value + ">"

	case rdf.Literal:
		return "\"" + n.Value + "\""

	case rdf.BlankNode:
		return "_:" + n.Value
	}

	return ""
}
