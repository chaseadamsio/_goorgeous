package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
)

func (p *parser) makeHorizontalRule(parent ast.Node, current, end int) {
	node := ast.NewHorizontalRuleNode(parent, p.items[current:end])
	parent.Append(node)
}

func (p *parser) onlyFollowsWhitespace(current int) bool {
	current--
	for 0 <= current {
		if p.items[current].IsNewline() {
			return true
		}
		if !(p.items[current].IsSpace() || p.items[current].IsTab()) {
			return false
		}
		current--
	}
	return true
}

func (p *parser) matchesHorizontalRule(current int) (found bool, start int) {
	itemsLength := len(p.items)
	testWidth := 5
	if !(p.items[current].IsDash() && current+testWidth < itemsLength) {
		return false, -1
	}

	if !p.onlyFollowsWhitespace(current) {
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
