package ast

import "github.com/chaseadamsio/goorgeous/lex"

type ListItemNode struct {
	NodeType
	parent   Node
	value    string
	start    int
	end      int
	children []Node
}

func NewListItemNode(start, end int, parent Node, items []lex.Item) *ListItemNode {
	node := &ListItemNode{
		NodeType: "ListItem",
		parent:   parent,
		start:    start,
		end:      end,
	}

	return node
}

// Type returns the type of node this is
func (n *ListItemNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *ListItemNode) String() string {
	return n.value
}

func (n ListItemNode) Children() []Node {
	return n.children
}

func (n *ListItemNode) Parent() Node {
	return n.parent
}

func (n *ListItemNode) Append(child Node) {
	n.children = append(n.children, child)
}
