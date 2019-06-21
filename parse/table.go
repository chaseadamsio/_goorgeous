package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeTable(parent ast.Node, start, end int) {
	node := ast.NewTableNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start, end)
}

func (p *parser) matchesTable(current int) (found bool, end int) {
	itemsLength := len(p.items)
	token := p.items[current]

	if !token.IsPipe() {
		return false, -1
	}
	if current < itemsLength && current == 0 || p.items[current-1].IsNewline() {
		for current < itemsLength {
			token := p.items[current]
			if token.IsNewline() {
				if current < itemsLength && (p.items[current+1].Type() != lex.ItemPipe) {
					if p.items[current+1].IsEOF() {
						current++
						continue
					}
					break
				}
			}
			current++
		}
		return true, current
	}
	return false, -1
}

func findTable(items []lex.Item) (end int) {
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		token := items[current]
		if token.IsNewline() {
			if current < itemsLength && (items[current+1].Type() != lex.ItemPipe) {
				if items[current+1].IsEOF() {
					current++
					continue
				}
				return current
			}
		}
		current++
	}
	return itemsLength
}
