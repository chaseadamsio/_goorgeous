package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeTable(parent ast.Node, items []lex.Item) (end int) {
	return 0
}

func isTable(token lex.Item, items []lex.Item, current int) bool {
	itemsLength := len(items)

	if !token.IsPipe() {
		return false
	}
	if current < itemsLength && current == 0 || items[current-1].IsNewline() {
		return true
	}
	return false
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
