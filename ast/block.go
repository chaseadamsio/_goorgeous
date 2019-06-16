package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type GreaterBlockNode struct {
	NodeType
	parent Node
	Key    string
	Value  string
	Start  int
	End    int
}

func NewGreaterBlockNode(parent Node, items []lex.Item) *GreaterBlockNode {
	node := &GreaterBlockNode{
		NodeType: "GreaterBlock",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}
	return node
}

// Type returns the type of node this is
func (n *GreaterBlockNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *GreaterBlockNode) String() string {
	return n.Key + ":" + n.Value
}

func (n GreaterBlockNode) Children() []Node {
	return nil
}

func (n *GreaterBlockNode) Parent() Node {
	return n.parent
}

func (n *GreaterBlockNode) Append(child Node) {
}
