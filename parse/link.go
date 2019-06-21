package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) newLink(parent ast.Node, start, end int) {
	node := ast.NewLinkNode(parent, p.items[start:end])
	parent.Append(node)
	textCurrent, textEnd := findLinkText(p.items[start:end])
	if textCurrent < textEnd {
		p.walkElements(node, textCurrent, textEnd)
	}
}

func findLinkBoundary(items []lex.Item) int {
	current := 0
	itemsLength := len(items)
	foundInsideClosingBracket := false

	for current < itemsLength {
		currItem := items[current]
		if currItem.IsBracket() && currItem.Value() == "]" {
			if foundInsideClosingBracket ||
				(current+1 < itemsLength && items[current+1].Value() == "]") {
				return current + 1
			}
			if current+1 < itemsLength && items[current+1].Value() == "[" {
				foundInsideClosingBracket = true
			}
		}
		current++
	}
	return itemsLength
}

func findLinkText(items []lex.Item) (start, end int) {
	current := 0
	itemsLength := len(items)
	foundLinkOpenBracket := false

	for current < itemsLength {
		currItem := items[current]
		if currItem.IsBracket() && currItem.Value() == "]" {
			if current+1 < itemsLength && items[current+1].Value() == "[" {
				foundLinkOpenBracket = true
				current = current + 2 // get to the content of the link
				start = current
				continue
			}
			if !foundLinkOpenBracket {
				return 0, 0
			}
			end = current
			return start, end
		}
		current++
	}
	return start, end
}

func (p *parser) matchesLink(current int) (found bool, end int) {
	itemsLength := len(p.items)

	// a link will always start with two brackets [[
	if !((p.items[current].IsBracket() && p.items[current].Value() == "[") &&
		(current < itemsLength && p.items[current+1].IsBracket() && p.items[current+1].Value() == "[")) {

		return false, -1
	}
	current = current + 2

	foundLinkCloseBracket := false

	for current < itemsLength {
		// we only care about brackets if we've closed the link part
		if p.items[current].IsBracket() && (foundLinkCloseBracket || p.items[current].Value() == "]") {
			// we found a ] to get here and if the next bracket is a closing bracket, this is a link!
			if current < itemsLength && p.items[current+1].IsBracket() && p.items[current+1].Value() == "]" {
				return true, current
			}
			// we hadn't found a closing bracket, but the current bracket is a closing bracket
			// because the next bracket opens the description
			if current < itemsLength && p.items[current+1].IsBracket() && p.items[current+1].Value() == "[" {
				foundLinkCloseBracket = true
			}
			// this is a self-describing link
			if foundLinkCloseBracket && p.items[current+1].IsBracket() && p.items[current+1].Value() == "]" {
				return true, -1
			}
			// if it's a bracket and it's not a closing bracket and we haven't found a closed bracket yet, this isn't a link
		} else if p.items[current].IsBracket() && p.items[current].Value() != "]" {
			return false, -1
		}
		current++
	}
	return false, -1
}
