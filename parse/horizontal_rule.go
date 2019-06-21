package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeHorizontalRule(parent ast.Node, items []lex.Item) (end int) {
	node := ast.NewHorizontalRuleNode(parent, items)
	parent.Append(node)
	return 5
}

func onlyFollowsWhitespace(current int, items []lex.Item) bool {
	current--
	for 0 < current {
		if items[current].IsNewline() {
			return true
		}
		if !(items[current].IsSpace() || items[current].IsTab()) {
			return false
		}
		current--
	}
	return true
}

func (p parser) matchesHorizontalRule(current int) (found bool, start int) {
	itemsLength := len(p.items)
	testWidth := 5
	if !(p.items[current].IsDash() && current+testWidth < itemsLength) {
		return false, -1
	}

	if !onlyFollowsWhitespace(current, p.items) {
		return false, -1
	}

	foundCount := 1
	current++
	for current < itemsLength {
		if !p.items[current].IsDash() {
			break
		}
		current++
		foundCount++
	}

	if foundCount != testWidth {
		return false, -1
	}

	return true, current
}
