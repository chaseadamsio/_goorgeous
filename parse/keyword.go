package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeKeyword(parent ast.Node, start, end int) {
	if parent.Type() == "Root" {
		node := ast.NewSectionNode(parent, p.items[start:end])
		parent.Append(node)
		parent = node
	}

	node := ast.NewKeywordNode(parent, p.items[start:end])
	parent.Append(node)
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

func (p *parser) matchesGreaterBlock(current int) (found bool, end int) {
	itemsLength := len(p.items)
	foundEnd := false

	if found, _ := p.matchesKeyword(current); !found {
		return false, -1
	}

	if current+2 < itemsLength && p.items[current+2].Value() != "BEGIN" {
		return false, -1
	}
	current = current + 2

	for current < itemsLength {
		if foundEnd && (current+1 == itemsLength || p.items[current].IsNewline() || p.items[current].IsEOF()) {
			return true, itemsLength
		}
		if p.items[current].Type() == lex.ItemHash {
			if itemsLength > current && p.items[current+1].Type() == lex.ItemPlus {
				if itemsLength > current+1 && p.items[current+2].Value() == "END" {
					foundEnd = true
				}
			}
		}
		current++
	}
	return false, -1
}
