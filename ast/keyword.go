package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type KeywordNode struct {
	NodeType
	parent Node
	key    string
	value  string
	start  int
	end    int
}

func NewKeywordNode(start, end int, parent Node, items []lex.Item) *KeywordNode {
	node := &KeywordNode{
		NodeType: "Keyword",
		parent:   parent,
		start:    start,
		end:      end,
	}

	node.parse(items)
	return node
}

func (n *KeywordNode) parse(items []lex.Item) {
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
func (n *KeywordNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *KeywordNode) String() string {
	return n.key + ":" + n.value
}

func (n KeywordNode) Children() []Node {
	return nil
}

func (n *KeywordNode) Parent() Node {
	return n.parent
}

func (n *KeywordNode) Append(child Node) {
}
