package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeKeyword(parent ast.Node, start, end int) {

	node := ast.NewKeywordNode(parent, p.items[start:end])

	p.parseKeyword(node, start, end)

	parent.Append(node)
}

func (p *parser) parseKeyword(node *ast.KeywordNode, start, end int) {
	var value []string
	current := start
	foundKeywordIdentifier := false
	foundColon := false
	for current < end {
		if !foundKeywordIdentifier && p.items[current].IsHash() && p.items[current+1].IsPlus() {
			foundKeywordIdentifier = true
			current = current + 2
			node.Key = strings.ToUpper(p.items[current].Value())
			continue
		}

		if foundColon {
			value = append(value, p.items[current].Value())
		}

		if p.items[current].IsColon() {
			foundColon = true
		}

		current++
	}

	node.Value = strings.Join(value, "")

}

func (p *parser) matchesKeyword(current int) (found bool, end int) {
	token := p.items[current]
	itemsLength := len(p.items)

	if !token.IsHash() {
		return false, -1
	}
	if current < itemsLength && p.items[current+1].Type() != lex.ItemPlus {
		return false, -1
	}
	if current == 0 || p.items[current-1].IsNewline() {
		for current < itemsLength {
			if p.items[current].IsNewline() {
				return true, current
			}
			current++
		}
	}
	return false, -1
}
