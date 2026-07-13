package rdf

import (
	"fmt"

	rdf "github.com/tggo/goRDFlib"
	"github.com/tggo/goRDFlib/graph"
	"github.com/tggo/goRDFlib/term"
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
	Graph   *graph.Graph // TODO do not expose Graph; just temporary to use it in serializer
	Triples []Triple
}

func NewTripleStore() *TripleStore {
	ts := TripleStore{}
	ts.Graph = graph.NewGraph()
	return &ts
}

func (store *TripleStore) BindPrefixes(prefixes map[string]string) error {
	for prefix, ns := range prefixes {
		uriRef, err := rdf.NewURIRef(ns)

		if err != nil {
			return fmt.Errorf("Unable to bind prefix %q with namespace %q: %w", prefix, ns, err)
		}

		store.Graph.Bind(prefix, uriRef)
	}

	return nil
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

	var o term.Term
	if object.Type == Literal {
		var opt rdf.LiteralOption

		if len(object.DataType) == 0 {
			opt = rdf.WithDatatype(rdf.XSDString)
		} else {
			opt = rdf.WithDatatype(rdf.NewURIRefUnsafe(rdf.XSDNamespace + object.DataType))
		}

		o = rdf.NewLiteral(object.Value, opt)
	} else {
		o, err = rdf.NewURIRef(object.Value)

		if err != nil {
			return fmt.Errorf("Unable to create object %w", err)
		}
	}
	store.Graph.Add(s, p, o)

	return nil
}

func (store *TripleStore) NumTriples() int {
	return store.Graph.Len()
}
