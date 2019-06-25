package ast

import "github.com/chaseadamsio/goorgeous/lex"

func NewSectionNode(parent Node, items []lex.Item) *SectionNode {
	node := &SectionNode{
		NodeType: "Section",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

func (n *SectionNode) Copy() *SectionNode {
	if n == nil {
		return nil
	}
	return &SectionNode{
		NodeType:      n.NodeType,
		parent:        n.Parent(),
		Start:         n.Start,
		End:           n.End,
		ChildrenNodes: n.ChildrenNodes,
	}
}

type SectionNode struct {
	NodeType
	parent        Node
	ChildrenNodes []Node
	Start, End    int
}

func (n SectionNode) Type() NodeType {
	return n.NodeType
}

func (n SectionNode) String() string {
	return ""
}

func (n SectionNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *SectionNode) Parent() Node {
	return n.parent
}

func (n *SectionNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
