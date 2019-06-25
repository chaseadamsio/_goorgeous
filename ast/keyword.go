package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type KeywordNode struct {
	NodeType
	parent Node
	Key    string
	Value  string
	Start  int
	End    int
}

func NewKeywordNode(parent Node, items []lex.Item) *KeywordNode {
	node := &KeywordNode{
		NodeType: "Keyword",
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

	return node
}

func (n *KeywordNode) Copy() *KeywordNode {
	if n == nil {
		return nil
	}
	return &KeywordNode{
		NodeType: n.NodeType,
		parent:   n.Parent(),
		Start:    n.Start,
		End:      n.End,
		Key:      n.Key,
		Value:    n.Value,
	}
}

// Type returns the type of node this is
func (n *KeywordNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *KeywordNode) String() string {
	return n.Key + ":" + n.Value
}

func (n KeywordNode) Children() []Node {
	return nil
}

func (n *KeywordNode) Parent() Node {
	return n.parent
}

func (n *KeywordNode) Append(child Node) {
}
