package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

// HeadlineNode is a node that represents a Headline
type HeadlineNode struct {
	NodeType
	Start         int
	End           int
	Depth         int
	parent        Node
	rawvalue      string
	ChildrenNodes []Node
	Keyword       string
}

func NewHeadlineNode(depth int, parent Node, items []lex.Item) *HeadlineNode {
	node := &HeadlineNode{
		NodeType: "Headline",
		Depth:    depth,
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *HeadlineNode) Type() NodeType {
	return n.NodeType
}

func (n *HeadlineNode) String() string {
	return n.rawvalue
}

func (n HeadlineNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *HeadlineNode) Parent() Node {
	return n.parent
}

func (n *HeadlineNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
