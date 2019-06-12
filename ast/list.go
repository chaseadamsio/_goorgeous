package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type ListNode struct {
	NodeType
	listType string
	parent   Node
	value    string
	start    int
	end      int
	children []Node
}

func NewListNode(listType string, parent Node, items []lex.Item) *ListNode {
	node := &ListNode{
		NodeType: "List",
		listType: listType,
		parent:   parent,
		start:    items[0].Offset(),
		end:      items[len(items)-1].Offset(),
	}

	var valStrs []string
	for _, item := range items {
		valStrs = append(valStrs, item.Value())
	}

	node.value = strings.Join(valStrs, "")
	return node
}

func findListItem(items []lex.Item) (start, end int) {
	itemsLength := len(items)
	for idx, item := range items {
		if (start > 0 && item.IsNewline()) || itemsLength == idx {
			end = idx
			return start, end
		}
		if item.Type() == lex.ItemDash && itemsLength > idx && items[idx+1].IsSpace() {
			start = idx + 2
		}
	}
	end = itemsLength
	return start, end
}

// Type returns the type of node this is
func (n *ListNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *ListNode) String() string {
	return n.value
}

func (n ListNode) Children() []Node {
	return n.children
}

func (n *ListNode) Parent() Node {
	return n.parent
}

func (n *ListNode) Append(child Node) {
	n.children = append(n.children, child)
}
