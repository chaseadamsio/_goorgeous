package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) newLink(parent ast.Node, items []lex.Item) (newEnd int) {
	end := findLink(items)
	node := ast.NewLinkNode(parent, items)
	parent.Append(node)
	textCurrent, textEnd := findLinkText(items)
	if textCurrent < textEnd {
		p.walkElements(node, items[textCurrent:textEnd])
	}
	return end
}

func findLink(items []lex.Item) int {
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

func isLink(items []lex.Item) bool {
	current := 0
	itemsLength := len(items)

	// a link will always start with two brackets [[
	if !((items[current].IsBracket() && items[current].Value() == "[") &&
		(current < itemsLength && items[current+1].IsBracket() && items[current+1].Value() == "[")) {

		return false
	}
	current = current + 2

	foundLinkCloseBracket := false

	for current < itemsLength {
		// we only care about brackets if we've closed the link part
		if items[current].IsBracket() && (foundLinkCloseBracket || items[current].Value() == "]") {
			// we found a ] to get here and if the next bracket is a closing bracket, this is a link!
			if current < itemsLength && items[current+1].IsBracket() && items[current+1].Value() == "]" {
				return true
			}
			// we hadn't found a closing bracket, but the current bracket is a closing bracket
			// because the next bracket opens the description
			if current < itemsLength && items[current+1].IsBracket() && items[current+1].Value() == "[" {
				foundLinkCloseBracket = true
			}
			// this is a self-describing link
			if foundLinkCloseBracket && items[current+1].IsBracket() && items[current+1].Value() == "]" {
				return true
			}
			// if it's a bracket and it's not a closing bracket and we haven't found a closed bracket yet, this isn't a link
		} else if items[current].IsBracket() && items[current].Value() != "]" {
			return false
		}
		current++
	}
	return false
}
