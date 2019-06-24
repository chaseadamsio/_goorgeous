package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type FootnoteReferenceNode struct {
	NodeType
	parent        Node
	ReferenceType string
	Label         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewFootnoteReferenceNode(parent Node, items []lex.Item) *FootnoteReferenceNode {
	node := &FootnoteReferenceNode{
		NodeType: "FootnoteReference",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *FootnoteReferenceNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *FootnoteReferenceNode) String() string {
	return ""
}

func (n FootnoteReferenceNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *FootnoteReferenceNode) Parent() Node {
	return n.parent
}

func (n *FootnoteReferenceNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}

type FootnoteDefinitionNode struct {
	NodeType
	parent        Node
	rawvalue      string
	Label         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewFootnoteDefinitionNode(parent Node, items []lex.Item) *FootnoteDefinitionNode {
	node := &FootnoteDefinitionNode{
		NodeType: "FootnoteDefinition",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *FootnoteDefinitionNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *FootnoteDefinitionNode) String() string {
	return n.rawvalue
}

func (n FootnoteDefinitionNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *FootnoteDefinitionNode) Parent() Node {
	return n.parent
}

func (n *FootnoteDefinitionNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
