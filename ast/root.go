package ast

import (
	"encoding/json"
)

func NewRootNode(start, end int) *RootNode {
	node := &RootNode{
		NodeType: "Root",
		Start:    start,
		End:      end,
	}

	return node
}

type RootNode struct {
	NodeType
	ChildrenNodes []Node
	Start, End    int
}

func (n *RootNode) Copy() *RootNode {
	if n == nil {
		return nil
	}
	return &RootNode{
		NodeType:      n.NodeType,
		Start:         n.Start,
		End:           n.End,
		ChildrenNodes: n.ChildrenNodes,
	}
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
