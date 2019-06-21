package ast

import (
	"strings"

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

	node.parse(items)

	return node
}

func (n *HeadlineNode) parse(items []lex.Item) {
	var headlineVal []string
	for idx, item := range items {
		if hasKeyword(idx, items) {
			n.Keyword = item.Value()
		}
		headlineVal = append(headlineVal, item.Value())
	}
	n.rawvalue = strings.Join(headlineVal, "")
}

var keywords = map[string]struct{}{
	"TODO": struct{}{},
	"DONE": struct{}{},
}

func hasKeyword(idx int, items []lex.Item) bool {
	// keywords will only _ever_ occur in the first space
	if idx != 0 {
		return false
	}
	if _, ok := keywords[items[idx].Value()]; ok {
		return true
	}
	return false
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
