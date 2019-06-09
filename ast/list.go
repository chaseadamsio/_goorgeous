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

func NewUnorderedListNode(start, end int, parent Node, items []lex.Item) *ListNode {
	return newListNode("UNORDERED", start, end, parent, items)
}

func NewOrderedListNode(start, end int, parent Node, items []lex.Item) *ListNode {
	return newListNode("ORDERED", start, end, parent, items)
}

func newListNode(listType string, start, end int, parent Node, items []lex.Item) *ListNode {
	node := &ListNode{
		NodeType: "List",
		listType: listType,
		parent:   parent,
		start:    start,
		end:      end,
	}

	var valStrs []string
	for _, item := range items {
		valStrs = append(valStrs, item.Value())
	}

	node.value = strings.Join(valStrs, "")
	node.parse(items)
	return node
}

func (n *ListNode) parse(items []lex.Item) {
	start, end := 0, 0
	itemsLength := len(items)
	for end < itemsLength {
		nextStart, nextEnd := findListItem(items[start:itemsLength])
		start = nextStart + start
		end = nextEnd + end
		node := NewListItemNode(start, end, n, items[start:end])
		n.Append(node)
		start = end
	}
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
