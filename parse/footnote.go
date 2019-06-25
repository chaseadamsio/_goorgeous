package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
)

func (p *parser) makeFootnoteDefinition(parent ast.Node, start, end int) {
	node := ast.NewFootnoteDefinitionNode(parent, p.items[start:end])
	parent.Append(node)
	p.parseFootnoteDefinition(node, start, end)
}

const footnoteOpenWidth = 3

func (p *parser) parseFootnoteDefinition(node *ast.FootnoteDefinitionNode, start, end int) {
	start = start + footnoteOpenWidth // skip [ fn : because we know they're there
	current := start
	foundInlineDescription := false

	for current < end {
		currItem := p.items[current]

		if currItem.IsColon() || (!foundInlineDescription && currItem.IsBracket() && currItem.Value() == "]") {
			var labelStrs []string
			for idx := start; idx < current; idx++ {
				labelStrs = append(labelStrs, p.items[idx].Value())
			}
			node.Label = strings.Join(labelStrs, "")
		}

		if !foundInlineDescription && currItem.IsColon() {
			foundInlineDescription = true
			start = current + 1
		}

		if currItem.IsBracket() && currItem.Value() == "]" {
			current++
			break
		}

		current++
	}

	if foundInlineDescription {
		child := ast.NewTextNode(node, p.items[start:current])
		node.Append(child)
	}
	p.walkElements(node, current, end)
}

func (p *parser) matchesFootnoteDefinition(current int) (found bool, end int) {
	token := p.items[current]
	itemsLength := len(p.items)

	if 0 < current && !p.items[current-1].IsNewline() {
		return false, -1
	}

	if !(token.IsBracket() && token.Value() == "[") {
		return false, -1
	}

	if !(current < itemsLength && p.items[current+1].Value() == "fn") {
		return false, -1
	}

	if !(current+1 < itemsLength && p.items[current+2].IsColon()) {
		return false, -1
	}

	current = current + 2

	foundFootnoteLabel := false
	foundFootnoteInlineDefinition := false // when there's a definition, we can allow spaces
	foundFootnoteEnd := false

	for current < itemsLength {
		if !foundFootnoteInlineDefinition && p.items[current].IsSpace() {
			return false, -1
		}
		if p.items[current].IsNewline() || p.items[current].IsEOF() {
			if foundFootnoteEnd {
				return true, current
			}
			return false, -1
		}
		if p.items[current].IsColon() {
			foundFootnoteInlineDefinition = true
		}
		if foundFootnoteLabel && p.items[current].IsBracket() && p.items[current].Value() == "]" {
			foundFootnoteEnd = true
		}
		foundFootnoteLabel = true
		current++
	}

	return false, -1
}

func (p *parser) makeFootnoteReference(parent ast.Node, start, end int) {
	node := ast.NewFootnoteReferenceNode(parent, p.items[start:end])
	parent.Append(node)
	p.parseFootnoteReference(node, start, end)
}

func (p *parser) parseFootnoteReference(node *ast.FootnoteReferenceNode, start, end int) {
	start = start + footnoteOpenWidth // skip [ fn : because we know they're there
	current := start
	foundInlineDescription := false

	for current < end {
		currItem := p.items[current]

		if currItem.IsColon() || (!foundInlineDescription && currItem.IsBracket() && currItem.Value() == "]") {
			var labelStrs []string
			for idx := start; idx < current; idx++ {
				labelStrs = append(labelStrs, p.items[idx].Value())
			}
			node.Label = strings.Join(labelStrs, "")
		}

		if !foundInlineDescription && currItem.IsColon() {
			foundInlineDescription = true
			start = current + 1
		}

		if currItem.IsBracket() && currItem.Value() == "]" {
			break
		}

		current++
	}

	if foundInlineDescription {
		node.ReferenceType = "INLINE"
		child := ast.NewFootnoteDefinitionNode(node, p.items[start:current])
		node.Append(child)
		p.walkElements(child, start, current)
	}
	if current != end {
		p.walkElements(node, current, end-1)
	}
}

func (p *parser) matchesFootnoteReference(current int) (found bool, end int) {
	token := p.items[current]
	itemsLength := len(p.items)

	if !(token.IsBracket() && token.Value() == "[") {
		return false, -1
	}

	if !(current < itemsLength && p.items[current+1].Value() == "fn") {
		return false, -1
	}

	if !(current+1 < itemsLength && p.items[current+2].IsColon()) {
		return false, -1
	}

	current = current + 2

	foundFootnoteLabel := false
	foundFootnoteInlineDefinition := false // when there's a definition, we can allow spaces

	for current < itemsLength {
		if !foundFootnoteInlineDefinition && p.items[current].IsSpace() {
			return false, -1
		}
		if p.items[current].IsNewline() || p.items[current].IsEOF() {
			return false, -1
		}
		if p.items[current].IsColon() {
			foundFootnoteInlineDefinition = true
		}
		if foundFootnoteLabel && p.items[current].IsBracket() && p.items[current].Value() == "]" {
			return true, current + 1
		}
		foundFootnoteLabel = true
		current++
	}

	return false, -1
}
