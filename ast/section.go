package ast

func NewSectionNode(parent Node) *SectionNode {
	node := &SectionNode{
		NodeType: "Section",
		parent:   parent,
	}

	return node
}

type SectionNode struct {
	NodeType
	parent   Node
	children []Node
}

func (n SectionNode) Type() NodeType {
	return n.NodeType
}

func (n SectionNode) String() string {
	return ""
}

func (n SectionNode) Children() []Node {
	return n.children
}

func (n *SectionNode) Parent() Node {
	return n.parent
}

func (n *SectionNode) Append(child Node) {
	n.children = append(n.children, child)
}
