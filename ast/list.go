package ast

import (
	"github.com/chaseadamsio/goorgeous/lex"
)

type ListNode struct {
	NodeType
	ListType      string
	parent        Node
	Value         string
	Start         int
	End           int
	ChildrenNodes []Node
}

func NewListNode(listType string, parent Node, items []lex.Item) *ListNode {
	node := &ListNode{
		NodeType: "List",
		ListType: listType,
		parent:   parent,
		Start:    items[0].Offset(),
		End:      items[len(items)-1].End(),
	}

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
	return n.Value
}

func (n ListNode) Children() []Node {
	return n.ChildrenNodes
}

func (n *ListNode) Parent() Node {
	return n.parent
}

func (n *ListNode) Append(child Node) {
	n.ChildrenNodes = append(n.ChildrenNodes, child)
}
