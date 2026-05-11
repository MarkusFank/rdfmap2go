package rdf

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
	Triples []Triple
}

func (store *TripleStore) AddTriple(subject Node, predicate Node, object Node) {
	store.Triples = append(store.Triples, Triple{Subject: subject, Predicate: predicate, Object: object})
}
