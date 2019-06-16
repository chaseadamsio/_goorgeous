package ast

import "fmt"

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *HeadlineNode:

	case *RootNode:
	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))

	}
}

type Inspector func(Node) bool

func (f Inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

func Inspect(node Node, f func(node Node) bool) {
	Walk(Inspector(f), node)
}
