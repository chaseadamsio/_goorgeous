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

func newTextNode(typ NodeType, parent Node, items []lex.Item) *TextNode {
	var values []string
	for _, item := range items {
		values = append(values, item.Value())
	}
	node := &TextNode{
		NodeType: typ,
		val:      strings.Join(values, ""),
		parent:   parent,
		start:    items[0].Offset(),
		end:      items[len(items)-1].Offset(),
	}

	return node
}

type TextNode struct {
	NodeType
	parent   Node
	children []Node
	val      string
	start    int
	end      int
}

func (n TextNode) Type() NodeType {
	return n.NodeType
}

func (n TextNode) String() string {
	return n.val
}

func (n TextNode) Children() []Node {
	return n.children
}
func (n *TextNode) Parent() Node {
	return n.parent
}
func (n *TextNode) Append(child Node) {
	n.children = append(n.children, child)
}
