package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
)

func (p *parser) newLink(parent ast.Node, start, end int) {
	node := ast.NewLinkNode(parent, p.items[start:end])
	parent.Append(node)

	linkStart, linkEnd, descriptionStart, descriptionEnd := p.parseLink(start, end)

	var link []string
	for linkStart <= linkEnd {
		link = append(link, p.items[linkStart].Value())
		linkStart++
	}

	node.Link = strings.Join(link, "")
	if descriptionStart < descriptionEnd {
		p.walkElements(node, descriptionStart, descriptionEnd)
	}
}

func (p *parser) parseLink(current, end int) (linkStart, linkEnd, descriptionStart, descriptionEnd int) {
	foundLinkOpenBracket := false        // track the opening of the link
	foundDescriptionOpenBracket := false // track the opening of the description
	if p.items[current].Value() == "[" {
		current++ // advance and get rid of the overall link opening bracket
	}
	if p.items[end].Value() == "]" {
		end-- // remove the end overall link closing bracket
	}

	for current <= end {
		currItem := p.items[current]
		if currItem.IsBracket() {
			if currItem.Value() == "[" && !foundLinkOpenBracket { // found the opening
				foundLinkOpenBracket = true
				current++
				linkStart = current // skip the bracket itself
			} else if currItem.Value() == "[" && foundLinkOpenBracket {
				foundDescriptionOpenBracket = true
				current++
				descriptionStart = current // skip the bracket itself
			} else if currItem.Value() == "]" && foundLinkOpenBracket && !foundDescriptionOpenBracket {
				linkEnd = current - 1 // don't include the bracket itself
			} else if currItem.Value() == "]" && foundDescriptionOpenBracket {
				descriptionEnd = current
			}
		}
		current++
	}
	return linkStart, linkEnd, descriptionStart, descriptionEnd
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
				return true, current + 1
			}
			// we hadn't found a closing bracket, but the current bracket is a closing bracket
			// because the next bracket opens the description
			if current < itemsLength && p.items[current+1].IsBracket() && p.items[current+1].Value() == "[" {
				foundLinkCloseBracket = true
			}
			// this is a self-describing link
			if foundLinkCloseBracket && p.items[current+1].IsBracket() && p.items[current+1].Value() == "]" {
				return true, current + 1
			}
			// if it's a bracket and it's not a closing bracket and we haven't found a closed bracket yet, this isn't a link
		} else if p.items[current].IsBracket() && p.items[current].Value() != "]" {
			return false, -1
		}
		current++
	}
	return false, -1
}
