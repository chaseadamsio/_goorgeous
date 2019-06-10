package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

func NewTextNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Text", start, end, parent, items)
	return node
}

func NewBoldNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Bold", start, end, parent, items)
	return node
}

func NewItalicNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Italic", start, end, parent, items)
	return node
}

func NewVerbatimNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Verbatim", start, end, parent, items)
	return node
}

func NewStrikeThroughNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("StrikeThrough", start, end, parent, items)
	return node
}

func NewUnderlineNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Underline", start, end, parent, items)
	return node
}

func NewCodeNode(start, end int, parent Node, items []lex.Item) *TextNode {
	node := newTextNode("Code", start, end, parent, items)
	return node
}

func newTextNode(typ NodeType, start, end int, parent Node, items []lex.Item) *TextNode {
	var values []string
	for _, item := range items {
		values = append(values, item.Value())
	}
	node := &TextNode{
		NodeType: typ,
		val:      strings.Join(values, ""),
		parent:   parent,
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
