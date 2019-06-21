package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type TableNode struct {
	NodeType
	parent Node
	Key    string
	Value  string
	Start  int
	End    int
}

func NewTableNode(parent Node, items []lex.Item) *TableNode {
	node := &TableNode{
		NodeType: "Table",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	node.parse(items)
	return node
}

func (n *TableNode) parse(items []lex.Item) {
	var key string
	var val []string
	for idx, item := range items {
		if item.Type() == lex.ItemColon {
			key = items[idx-1].Value()
			continue
		} else if key != "" {
			val = append(val, item.Value())
		}
	}
	n.Key = key
	n.Value = strings.Join(val, "")
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
	return nil
}

func (n *TableNode) Parent() Node {
	return n.parent
}

func (n *TableNode) Append(child Node) {
}
