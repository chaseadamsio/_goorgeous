package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type ParagraphNode struct {
	NodeType
	parent   Node
	rawvalue string
	start    int
	end      int
	children []Node
}

func NewParagraphNode(start, end int, parent Node, items []lex.Item) *ParagraphNode {
	node := &ParagraphNode{
		NodeType: "Paragraph",
		parent:   parent,
		start:    start,
		end:      end,
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
	return n.children
}
func (n *ParagraphNode) Parent() Node {
	return n.parent
}
func (n *ParagraphNode) Append(child Node) {
	n.children = append(n.children, child)
}
