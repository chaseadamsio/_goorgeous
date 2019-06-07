package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

func NewTextNode(start, end int, parent Node, items []lex.Item) *TextNode {
	var values []string
	for _, item := range items {
		values = append(values, item.Value())
	}
	node := &TextNode{
		NodeType: "Text",
		val:      strings.Join(values, ""),
		parent:   parent,
	}

	return node
}

type TextNode struct {
	NodeType
	parent   Node
	children []Node
	val      string
	start    int
	end      int
}

func (n TextNode) Type() NodeType {
	return n.NodeType
}

func (n TextNode) String() string {
	return n.val
}

func (n TextNode) Children() []Node {
	return n.children
}
func (n *TextNode) Parent() Node {
	return n.parent
}
func (n *TextNode) Append(child Node) {
	n.children = append(n.children, child)
}
