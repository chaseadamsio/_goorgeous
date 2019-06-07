package ast

import "encoding/json"

func NewRootNode() *RootNode {
	node := &RootNode{
		NodeType: "Root",
	}

	return node
}

type RootNode struct {
	NodeType
	parent   Node
	children []Node
}

func (n *RootNode) String() string {
	tree := []treeTyp{
		{
			"Root",
			getChildren(n),
			"",
		},
	}
	out, _ := json.MarshalIndent(tree, "", "  ")
	return string(out)
}

func (n RootNode) Tree() RootNode {
	return n
}

func (n RootNode) Type() NodeType {
	return n.NodeType
}

func (n RootNode) Children() []Node {
	return n.children
}

func (n *RootNode) Parent() Node {
	return n.parent
}

func (n *RootNode) Append(child Node) {
	n.children = append(n.children, child)
}
