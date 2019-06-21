package ast

import "github.com/chaseadamsio/goorgeous/lex"

type HorizontalRuleNode struct {
	NodeType
	parent Node
	Key    string
	Value  string
	Start  int
	End    int
}

func NewHorizontalRuleNode(parent Node, items []lex.Item) *HorizontalRuleNode {
	node := &HorizontalRuleNode{
		NodeType: "HorizontalRule",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *HorizontalRuleNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *HorizontalRuleNode) String() string {
	return n.Key + ":" + n.Value
}

func (n HorizontalRuleNode) Children() []Node {
	return nil
}

func (n *HorizontalRuleNode) Parent() Node {
	return n.parent
}

func (n *HorizontalRuleNode) Append(child Node) {
}
