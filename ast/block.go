package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type GreaterBlockNode struct {
	NodeType
	parent   Node
	Name     string
	Value    string
	Language string
	Start    int
	End      int
}

func NewGreaterBlockNode(parent Node, items []lex.Item) *GreaterBlockNode {
	node := &GreaterBlockNode{
		NodeType: "GreaterBlock",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}
	return node
}

// Type returns the type of node this is
func (n *GreaterBlockNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *GreaterBlockNode) String() string {
	return n.Value
}

func (n GreaterBlockNode) Children() []Node {
	return nil
}

func (n *GreaterBlockNode) Parent() Node {
	return n.parent
}

func (n *GreaterBlockNode) Append(child Node) {
}

func (n *GreaterBlockNode) Copy() *GreaterBlockNode {
	if n == nil {
		return nil
	}
	return &GreaterBlockNode{
		NodeType: n.NodeType,
		Name:     n.Name,
		parent:   n.Parent(),
		Value:    n.Value,
		Language: n.Language,
		Start:    n.Start,
		End:      n.End,
	}
}

type FixedWidthNode struct {
	NodeType
	parent Node
	Value  string
	Start  int
	End    int
}

func NewFixedWidthNode(parent Node, items []lex.Item) *FixedWidthNode {
	node := &FixedWidthNode{
		NodeType: "FixedWidth",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}
	return node
}

// Type returns the type of node this is
func (n *FixedWidthNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *FixedWidthNode) String() string {
	return n.Value
}

func (n FixedWidthNode) Children() []Node {
	return nil
}

func (n *FixedWidthNode) Parent() Node {
	return n.parent
}

func (n *FixedWidthNode) Append(child Node) {
}

func (n *FixedWidthNode) Copy() *FixedWidthNode {
	if n == nil {
		return nil
	}
	return &FixedWidthNode{
		NodeType: n.NodeType,
		parent:   n.Parent(),
		Value:    n.Value,
		Start:    n.Start,
		End:      n.End,
	}
}
