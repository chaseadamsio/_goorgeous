package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type ListItemNode struct {
	NodeType
	parent        Node
	Bullet        string
	Value         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func (n *ListItemNode) Copy() *ListItemNode {
	if n == nil {
		return nil
	}
	return &ListItemNode{
		NodeType:      n.NodeType,
		parent:        n.Parent(),
		Bullet:        n.Bullet,
		Value:         n.Value,
		Start:         n.Start,
		End:           n.End,
		ChildrenNodes: n.ChildrenNodes,
	}
}

func NewListItemNode(parent Node, items []lex.Item) *ListItemNode {
	node := &ListItemNode{
		NodeType: "ListItem",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	var valStrs []string
	for _, item := range items {
		valStrs = append(valStrs, item.Value())
	}

	node.Value = strings.Join(valStrs, "")

	return node
}

// Type returns the type of node this is
func (n *ListItemNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *ListItemNode) String() string {
	return n.Value
}

func (n ListItemNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *ListItemNode) Parent() Node {
	return n.parent
}

func (n *ListItemNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
