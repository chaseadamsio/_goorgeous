package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type TableNode struct {
	NodeType
	parent Node
	key    string
	value  string
	start  int
	end    int
}

func NewTableNode(start, end int, parent Node, items []lex.Item) *TableNode {
	node := &TableNode{
		NodeType: "Table",
		parent:   parent,
		start:    start,
		end:      end,
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
	n.key = key
	n.value = strings.Join(val, "")
}

// Type returns the type of node this is
func (n *TableNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *TableNode) String() string {
	return n.key + ":" + n.value
}

func (n TableNode) Children() []Node {
	return nil
}

func (n *TableNode) Parent() Node {
	return n.parent
}

func (n *TableNode) Append(child Node) {
}
