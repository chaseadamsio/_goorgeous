package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type FootnoteReferenceNode struct {
	NodeType
	parent        Node
	ReferenceType string
	ID            int
	Label         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func (n *FootnoteReferenceNode) Copy() *FootnoteReferenceNode {
	if n == nil {
		return nil
	}
	return &FootnoteReferenceNode{
		NodeType:      n.NodeType,
		parent:        n.Parent(),
		Label:         n.Label,
		ID:            n.ID,
		ReferenceType: n.ReferenceType,
		Start:         n.Start,
		End:           n.End,
		ChildrenNodes: n.ChildrenNodes,
	}
}

func NewFootnoteReferenceNode(parent Node, items []lex.Item) *FootnoteReferenceNode {
	node := &FootnoteReferenceNode{
		NodeType:      "FootnoteReference",
		ReferenceType: "STANDARD",
		parent:        parent,
		Start:         items[0].Offset(),
		End:           items[len(items)-1].End(),
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
	ID            int
	Start         int
	End           int
	ChildrenNodes []Node
}

func (n *FootnoteDefinitionNode) Copy() *FootnoteDefinitionNode {
	if n == nil {
		return nil
	}
	return &FootnoteDefinitionNode{
		NodeType:      n.NodeType,
		parent:        n.Parent(),
		Label:         n.Label,
		ID:            n.ID,
		rawvalue:      n.rawvalue,
		Start:         n.Start,
		End:           n.End,
		ChildrenNodes: n.ChildrenNodes,
	}
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
