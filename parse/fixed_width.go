package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
)

func (p *parser) makeFixedWidth(parent ast.Node, start, end int) {
	node := ast.NewFixedWidthNode(parent, p.items[start:end])
	parent.Append(node)

	var val []string

	current := start
	for current < end {
		// don't include the colons that start new lines
		if p.items[current].IsColon() && p.onlyFollowsWhitespace(current) {
			current++
			continue
		}
		val = append(val, p.items[current].Value())
		current++
	}
	node.Value = strings.Join(val, "")
}

func (p *parser) matchesFixedWidth(current int) (found bool, end int) {

	if !p.onlyFollowsWhitespace(current) {
		return false, -1
	}

	if !p.items[current].IsColon() {
		return false, -1
	}

	current++
	for current < len(p.items) {
		if p.items[current].IsNewline() {
			if nextFound, nextEnd := p.matchesFixedWidth(current + 1); nextFound {
				current = nextEnd
				continue
			}
			break
		}
		current++
	}
	return true, current
}
