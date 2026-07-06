package rdf

import (
	"fmt"

	rdf "github.com/tggo/goRDFlib"
	"github.com/tggo/goRDFlib/graph"
)

type NodeType int

const (
	URI NodeType = iota
	Literal
	BlankNode
)

type Node struct {
	Type     NodeType
	Value    string
	DataType string
	Language string
}

type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

type TripleStore struct {
	graph   *graph.Graph
	Triples []Triple
}

func NewTripleStore() *TripleStore {
	ts := TripleStore{}
	ts.graph = graph.NewGraph()
	return &ts
}

func (store *TripleStore) AddTriple(subject Node, predicate Node, object Node) error {
	// store.Triples = append(store.Triples, Triple{Subject: subject, Predicate: predicate, Object: object})

	s, err := rdf.NewURIRef(subject.Value)

	if err != nil {
		return fmt.Errorf("Unable to create subject %w", err)
	}

	p, err := rdf.NewURIRef(predicate.Value)

	if err != nil {
		return fmt.Errorf("Unable to create predicate %w", err)
	}

	o := rdf.NewLiteral(object.Value)
	store.graph.Add(s, p, o)

	return nil
}
