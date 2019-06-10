package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type LinkNode struct {
	NodeType
	parent   Node
	rawvalue string
	link     string
	start    int
	end      int
	children []Node
}

func NewLinkNode(parent Node, items []lex.Item) *LinkNode {
	node := &LinkNode{
		NodeType: "Link",
		parent:   parent,
		start:    items[0].Offset(),
		end:      items[len(items)-1].Offset(),
	}

	node.parse(items)
	return node
}

func (n *LinkNode) parse(items []lex.Item) {
	current := 0
	itemsLength := len(items)
	var rawvalueStrs []string

	for current < itemsLength {
		rawvalueStrs = append(rawvalueStrs, items[current].Value())
		current++
	}

	n.rawvalue = strings.Join(rawvalueStrs, "")
}

// Type returns the type of node this is
func (n *LinkNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *LinkNode) String() string {
	return n.rawvalue
}

func (n LinkNode) Children() []Node {
	return n.children
}

func (n *LinkNode) Parent() Node {
	return n.parent
}

func (n *LinkNode) Append(child Node) {
	n.children = append(n.children, child)
}
