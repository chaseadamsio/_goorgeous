package ast

// NodeType is the type of a node
type NodeType string

type Node interface {
	Type() NodeType
	Parent() Node
	Children() []Node
	Append(Node)
}

type treeTyp struct {
	NodeType
	Children []treeTyp
	Value    string
}
