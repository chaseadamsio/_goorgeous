package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type GreaterBlockNode struct {
	NodeType
	parent Node
	key    string
	value  string
	start  int
	end    int
}

func NewGreaterBlockNode(start, end int, parent Node, items []lex.Item) *GreaterBlockNode {
	node := &GreaterBlockNode{
		NodeType: "GreaterBlock",
		parent:   parent,
		start:    start,
		end:      end,
	}

	node.parse(items)
	return node
}

func (n *GreaterBlockNode) parse(items []lex.Item) {
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
func (n *GreaterBlockNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *GreaterBlockNode) String() string {
	return n.key + ":" + n.value
}

func (n GreaterBlockNode) Children() []Node {
	return nil
}

func (n *GreaterBlockNode) Parent() Node {
	return n.parent
}

func (n *GreaterBlockNode) Append(child Node) {
}
