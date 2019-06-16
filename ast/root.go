package ast

import (
	"encoding/json"

	"github.com/chaseadamsio/goorgeous/lex"
)

func NewRootNode(items []lex.Item) *RootNode {
	node := &RootNode{
		NodeType: "Root",
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

type RootNode struct {
	NodeType
	ChildrenNodes []Node
	Start, End    int
}

func (n *RootNode) String() string {
	out, _ := json.MarshalIndent(n, "", "  ")
	return string(out)
}

func (n RootNode) Tree() RootNode {
	return n
}

func (n RootNode) Type() NodeType {
	return n.NodeType
}

func (n RootNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *RootNode) Parent() Node {
	return nil
}

func (n *RootNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
