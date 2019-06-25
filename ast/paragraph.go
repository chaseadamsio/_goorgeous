package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type ParagraphNode struct {
	NodeType
	parent        Node
	rawvalue      string
	Start         int
	End           int
	ChildrenNodes []Node
}

func (n *ParagraphNode) Copy() *ParagraphNode {
	if n == nil {
		return nil
	}
	return &ParagraphNode{
		NodeType:      n.NodeType,
		parent:        n.Parent(),
		rawvalue:      n.rawvalue,
		Start:         n.Start,
		End:           n.End,
		ChildrenNodes: n.ChildrenNodes,
	}
}

func NewParagraphNode(start, end int, parent Node, items []lex.Item) *ParagraphNode {
	node := &ParagraphNode{
		NodeType: "Paragraph",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	var rawvalueStrs []string
	for _, item := range items {
		rawvalueStrs = append(rawvalueStrs, item.Value())
	}
	node.rawvalue = strings.Join(rawvalueStrs, "")

	return node
}

// Type returns the type of node this is
func (n *ParagraphNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *ParagraphNode) String() string {
	return n.rawvalue
}

func (n ParagraphNode) Children() []Node {
	return n.ChildrenNodes
}
func (n *ParagraphNode) Parent() Node {
	return n.parent
}
func (n *ParagraphNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
