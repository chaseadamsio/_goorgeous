package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) matchesElementMarkup(current int, expectedTypeFunc func(lex.Item) bool) (found bool, end int) {
	foundOpeningChar := false
	itemsLength := len(p.items)
	for current < itemsLength {
		currToken := p.items[current]
		// check the current token against the expected type checker function
		if expectedTypeFunc(currToken) {
			nextTokenIsWhitespace := current+1 < itemsLength && p.items[current+1].IsWhitespace()

			// a match for opening character, it cannot precede a whitespace character
			if !foundOpeningChar && !nextTokenIsWhitespace {
				foundOpeningChar = true
				current++
				continue
			}
			// closing characters cannot follow a whitespace character
			if foundOpeningChar && !(current > 0 && p.items[current-1].IsWhitespace()) {
				// there is no next character in this collection of items
				if current+1 == itemsLength {
					return true, current + 1
				}
				// if it precedes EOF, Newline and Whitespace by this point, it's a match
				if current+1 < itemsLength &&
					(!p.items[current+1].IsWord() || p.items[current+1].IsEOF() || p.items[current+1].IsNewline() || p.items[current+1].IsWhitespace()) {
					return true, current + 1
				}
			}
			// if it's a newline, the first character or the end of the collection, we didn't find the expected type
		} else if !foundOpeningChar || currToken.IsNewline() || current == 0 || current == itemsLength {
			return false, -1
		}
		current++
	}
	return false, -1
}

func (p *parser) newBold(parent ast.Node, start, end int) {
	node := ast.NewBoldNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start+1, end-1)
}

// matchesBold returns true if a collection of items matches bold markup
func (p *parser) matchesBold(current int) (found bool, end int) {
	return p.matchesElementMarkup(current, func(currToken lex.Item) bool {
		return currToken.IsAsterisk()
	})
}

func (p *parser) newItalic(parent ast.Node, start, end int) {
	node := ast.NewItalicNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start+1, end-1)
}

// matchesItalic returns true if a collection of items matches italic markup
func (p *parser) matchesItalic(current int) (found bool, end int) {
	return p.matchesElementMarkup(current, func(currToken lex.Item) bool {
		return currToken.IsForwardSlash()
	})
}

func (p *parser) newVerbatim(parent ast.Node, start, end int) {
	node := ast.NewVerbatimNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start+1, end-1)
}

// matchesVerbatim returns true if a collection of items matches verbatim markup
func (p *parser) matchesVerbatim(current int) (found bool, end int) {
	return p.matchesElementMarkup(current, func(currToken lex.Item) bool {
		return currToken.IsEqual()
	})
}

func (p *parser) newStrikeThrough(parent ast.Node, start, end int) {
	node := ast.NewStrikeThroughNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start+1, end-1)
}

// matchesStrikeThrough returns true if a collection of items matches strike through markup
func (p *parser) matchesStrikeThrough(current int) (found bool, end int) {
	return p.matchesElementMarkup(current, func(currToken lex.Item) bool {
		return currToken.IsPlus()
	})
}

func (p *parser) newUnderline(parent ast.Node, start, end int) {
	node := ast.NewUnderlineNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start+1, end-1)
}

// matchesUnderline returns true if a collection of items matches underline markup
func (p *parser) matchesUnderline(current int) (found bool, end int) {
	return p.matchesElementMarkup(current, func(currToken lex.Item) bool {
		return currToken.IsUnderscore()
	})
}

func (p *parser) newCode(parent ast.Node, start, end int) {
	node := ast.NewCodeNode(parent, p.items[start:end])
	parent.Append(node)
	p.walkElements(node, start+1, end-1)
}

// matchesCode returns true if a collection of items matches code markup
func (p *parser) matchesCode(current int) (found bool, end int) {
	return p.matchesElementMarkup(current, func(currToken lex.Item) bool {
		return currToken.IsTilde()
	})
}
