package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p parser) makeFootnote(parent ast.Node, items []lex.Item) {
	// if start < current {
	// 	if parent.Type() != "Section" {
	// 		parent = findClosestSectionNode(parent, p.items[current:end])
	// 	}
	// 	node := ast.NewParagraphNode(start, current-1, parent, p.items[start:current-1])
	// 	p.walkElements(node, start, end)
	// 	parent.Append(node)
	// }
	// peekStart := current
	// end := peekStart + peekToNewLine(p.items[peekStart:])

	// node := ast.NewFootnoteDefinitionNode(parent, p.items[current:end])
	// parent.Append(node)
}

func (p parser) makeFootnoteDefinition(parent ast.Node, start, end int) {
	node := ast.NewFootnoteDefinitionNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start, end)
}

func (p *parser) matchesFootnoteDefinition(current int) (found bool, end int) {
	itemsLength := len(p.items)
	token := p.items[current]

	if !(token.IsBracket() && token.Value() == "[") {
		return false, -1
	}

	if current < itemsLength && p.items[current+1].Value() != "fn" {
		return false, -1
	}

	if current+1 < itemsLength && !p.items[current+2].IsColon() {
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
			return true, current
		}
		foundFootnoteLabel = true
		current++
	}

	return false, -1
}
