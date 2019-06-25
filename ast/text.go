package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

func NewTextNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Text", parent, items)
	return node
}

func NewBoldNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Bold", parent, items)
	return node
}

func NewItalicNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Italic", parent, items)
	return node
}

func NewVerbatimNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Verbatim", parent, items)
	return node
}

func NewStrikeThroughNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("StrikeThrough", parent, items)
	return node
}

func NewUnderlineNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Underline", parent, items)
	return node
}

func NewCodeNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Code", parent, items)
	return node
}

func NewEnDashNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("EnDash", parent, items)
	return node
}

func NewMDashNode(parent Node, items []lex.Item) *TextNode {
	node := newTextNode("MDash", parent, items)
	return node
}

func newTextNode(typ NodeType, parent Node, items []lex.Item) *TextNode {
	var values []string
	for _, item := range items {
		values = append(values, item.Value())
	}
	node := &TextNode{
		NodeType: typ,
		Val:      strings.Join(values, ""),
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

type TextNode struct {
	NodeType
	parent        Node
	ChildrenNodes []Node
	Val           string
	Start         int
	End           int
}

func (n *TextNode) Copy() *TextNode {
	if n == nil {
		return nil
	}
	return &TextNode{
		NodeType:      n.NodeType,
		parent:        n.Parent(),
		Start:         n.Start,
		End:           n.End,
		Val:           n.Val,
		ChildrenNodes: n.ChildrenNodes,
	}
}

func (n TextNode) Type() NodeType {
	return n.NodeType
}

func (n TextNode) String() string {
	return n.Val
}

func (n TextNode) Children() []Node {
	return n.ChildrenNodes
}
func (n *TextNode) Parent() Node {
	return n.parent
}
func (n *TextNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
