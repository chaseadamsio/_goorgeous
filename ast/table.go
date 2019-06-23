package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type TableNode struct {
	NodeType
	parent        Node
	Key           string
	Value         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewTableNode(parent Node, items []lex.Item) *TableNode {
	node := &TableNode{
		NodeType: "Table",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *TableNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *TableNode) String() string {
	return n.Key + ":" + n.Value
}

func (n TableNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *TableNode) Parent() Node {
	return n.parent
}

func (n *TableNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}

type TableRowNode struct {
	NodeType
	parent        Node
	Key           string
	Value         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewTableRowNode(parent Node, items []lex.Item) *TableRowNode {
	node := &TableRowNode{
		NodeType: "TableRow",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *TableRowNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *TableRowNode) String() string {
	return n.Key + ":" + n.Value
}

func (n TableRowNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *TableRowNode) Parent() Node {
	return n.parent
}

func (n *TableRowNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}

type TableCellNode struct {
	NodeType
	parent        Node
	Key           string
	Value         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewTableCellNode(parent Node, items []lex.Item) *TableCellNode {
	node := &TableCellNode{
		NodeType: "TableCell",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

// Type returns the type of node this is
func (n *TableCellNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *TableCellNode) String() string {
	return n.Key + ":" + n.Value
}

func (n TableCellNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *TableCellNode) Parent() Node {
	return n.parent
}

func (n *TableCellNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
