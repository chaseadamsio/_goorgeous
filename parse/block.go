package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
)

var names = map[string]struct{}{
	"CENTER":  struct{}{},
	"QUOTE":   struct{}{},
	"COMMENT": struct{}{},
	"EXAMPLE": struct{}{},
	"EXPORT":  struct{}{},
	"SRC":     struct{}{},
	"VERSE":   struct{}{},
}

func (p *parser) makeGreaterBlock(parent ast.Node, start, end int) {
	node := ast.NewGreaterBlockNode(parent, p.items[start:end])

	parent.Append(node)
	var val []string

	foundBlockOpening := false

	current := start
	for current < end {
		if !foundBlockOpening && p.items[current].IsNewline() {
			foundBlockOpening = true
			p.parseBlockParameters(node, start, current)
		} else if p.items[current].IsHash() && p.items[current+1].IsPlus() &&
			p.items[current+2].IsText() && strings.ToUpper(p.items[current+2].Value()) == "END" {
			break
		} else if foundBlockOpening {
			item := p.items[current]
			val = append(val, item.Value())
		}
		current++
	}

	node.Value = strings.Join(val, "")
}

func (p *parser) parseBlockParameters(node *ast.GreaterBlockNode, start, end int) {
	current := start
	if p.items[current].IsHash() {
		current++
	}
	if p.items[current].IsPlus() {
		current++
	}
	if p.items[current].IsText() && strings.ToUpper(p.items[current].Value()) == "BEGIN" {
		current++
	}

	if p.items[current].IsUnderscore() {
		current++
	}

	node.Name = strings.ToUpper(p.items[current].Value())
	current++

	if p.items[current].IsSpace() {
		current++
	}

	if node.Name == "SRC" && p.items[current].IsText() {
		node.Language = p.items[current].Value()
	}

	for current < end {
		current++
	}
}

func (p *parser) matchesGreaterBlock(current int) (found bool, end int) {
	itemsLength := len(p.items)
	foundEnd := false
	name := ""

	if found, _ := p.matchesKeyword(current); !found {
		return false, -1
	}

	if !(current+2 < itemsLength && strings.ToUpper(p.items[current+2].Value()) == "BEGIN") &&
		!(current+3 < itemsLength && p.items[current+3].IsUnderscore()) {
		return false, -1
	}
	current = current + 4
	name = p.items[current].Value()

	for current < itemsLength {
		if foundEnd && (current+1 == itemsLength || p.items[current].IsNewline() || p.items[current].IsEOF()) {
			return true, itemsLength
		}
		if p.items[current].IsHash() && current < itemsLength && p.items[current+1].IsPlus() {
			if current+1 < itemsLength && strings.ToUpper(p.items[current+2].Value()) == "END" &&
				current+2 < itemsLength && p.items[current+3].IsUnderscore() &&
				current+3 < itemsLength && p.items[current+4].Value() == name {
				if current+4 < itemsLength && p.items[current+5].IsNewline() {
					current++
				}
				return true, current + 4
			}
		}
		current++
	}
	return false, -1
}
