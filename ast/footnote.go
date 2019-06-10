package ast

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/lex"
)

type FootnoteDefinitionNode struct {
	NodeType
	parent   Node
	rawvalue string
	label    string
	start    int
	end      int
	children []Node
}

func NewFootnoteDefinitionNode(parent Node, items []lex.Item) *FootnoteDefinitionNode {
	node := &FootnoteDefinitionNode{
		NodeType: "FootnoteDefinition",
		parent:   parent,
		start:    items[0].Offset(),
		end:      items[len(items)-1].Offset(),
	}

	node.parse(items)
	return node
}

// [, fn, :
const footnoteOpenWidth = 3

func (n *FootnoteDefinitionNode) parse(items []lex.Item) {
	// 		{"[fn:1] The link is: https://orgmode.org", 0, true},

	current := footnoteOpenWidth
	start := current
	itemsLength := len(items)
	labelFound := false
	inlineDescription := false
	var rawvalueStrs []string

	for current < itemsLength {
		currItem := items[current]

		// this footnote contains a description
		if (currItem.IsColon() || inlineDescription) && labelFound {
			inlineDescription = true
			rawvalueStrs = append(rawvalueStrs, currItem.Value())
			// TODO footnote needs to have a text node child
		} else if currItem.IsBracket() && currItem.Value() == "]" {
			var labelStrs []string

			if start == current {
				n.label = ""
			} else {
				for idx := start; idx < current; idx++ {
					labelStrs = append(labelStrs, items[idx].Value())
				}
				n.label = strings.Join(labelStrs, "")
			}
			labelFound = true
			start = current
		} else if labelFound {
			rawvalueStrs = append(rawvalueStrs, currItem.Value())
		}

		current++
	}

	n.rawvalue = strings.Join(rawvalueStrs, "")
}

// Type returns the type of node this is
func (n *FootnoteDefinitionNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *FootnoteDefinitionNode) String() string {
	return n.rawvalue
}

func (n FootnoteDefinitionNode) Children() []Node {
	return n.children
}

func (n *FootnoteDefinitionNode) Parent() Node {
	return n.parent
}

func (n *FootnoteDefinitionNode) Append(child Node) {
	n.children = append(n.children, child)
}
