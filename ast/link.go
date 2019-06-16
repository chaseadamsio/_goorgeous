package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type LinkNode struct {
	NodeType
	parent        Node
	rawvalue      string
	Link          string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewLinkNode(parent Node, items []lex.Item) *LinkNode {
	node := &LinkNode{
		NodeType: "Link",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
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
	return n.ChildrenNodes
}

func (n *LinkNode) Parent() Node {
	return n.parent
}

func (n *LinkNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
