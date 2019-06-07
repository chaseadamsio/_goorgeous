package ast

// NodeType is the type of a node
type NodeType string

type Node interface {
	Type() NodeType
	Parent() Node
	String() string
	Children() []Node
	Append(Node)
}

type treeTyp struct {
	NodeType
	Children []treeTyp
	Value    string
}

func getChildren(node Node) []treeTyp {
	var tree []treeTyp

	if len(node.Children()) == 0 {
		return nil
	}

	for _, child := range node.Children() {
		tree = append(tree, treeTyp{
			child.Type(),
			getChildren(child),
			child.String(),
		})
	}

	return tree
}

type Pos struct {
	Start int
	End   int
}
