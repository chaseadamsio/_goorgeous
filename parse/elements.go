package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func isElementMarkup(items []lex.Item, expectedTypeFunc func(lex.Item) bool) bool {
	foundOpeningChar := false
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		currToken := items[current]
		// check the current token against the expected type checker function
		if expectedTypeFunc(currToken) {
			nextTokenIsWhitespace := current+1 < itemsLength && items[current+1].IsWhitespace()

			// a match for opening character, it cannot precede a whitespace character
			if !foundOpeningChar && !nextTokenIsWhitespace {
				foundOpeningChar = true
				current++
				continue
			}
			// closing characters cannot follow a whitespace character
			if foundOpeningChar && !(current > 0 && items[current-1].IsWhitespace()) {
				// there is no next character in this collection of items
				if current+1 == itemsLength {
					return true
				}
				// if it precedes EOF, Newline and Whitespace by this point, it's a match
				if current+1 < itemsLength &&
					(!items[current+1].IsWord() || items[current+1].IsEOF() || items[current+1].IsNewline() || items[current+1].IsWhitespace()) {
					return true
				}
			}
			// if it's a newline, the first character or the end of the collection, we didn't find the expected type
		} else if currToken.IsNewline() || current == 0 || current == itemsLength {
			return false
		}
		current++
	}
	return false
}

func findElementMarkup(items []lex.Item, expectedTypeFunc func(lex.Item) bool) int {
	current := 0
	itemsLength := len(items)
	foundLeftMarked := false

	for current < itemsLength {
		currItem := items[current]
		if expectedTypeFunc(currItem) {
			if foundLeftMarked {
				return current
			}
			foundLeftMarked = true
		}
		current++
	}
	return -1
}

func (p *parser) newBold(parent ast.Node, items []lex.Item) (end int) {
	current := 0
	end = findBold(items)
	node := ast.NewBoldNode(parent, items)
	parent.Append(node)
	p.walkElements(node, items[current+1:end])
	return end
}

// isBold returns true if a collection of items matches bold markup
func isBold(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsAsterisk()
	})
}

// findBold finds the end item of a bold collection
func findBold(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsAsterisk()
	})
}

func (p *parser) newItalic(parent ast.Node, items []lex.Item) (end int) {
	current := 0
	end = findItalic(items)
	node := ast.NewItalicNode(parent, items)
	parent.Append(node)
	p.walkElements(node, items[current+1:end])
	return end
}

// isItalic returns true if a collection of items matches italic markup
func isItalic(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsForwardSlash()
	})
}

// findItalic finds the end item of a italic collection
func findItalic(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsForwardSlash()
	})
}

func (p *parser) newVerbatim(parent ast.Node, items []lex.Item) (end int) {
	current := 0
	end = findVerbatim(items)
	node := ast.NewVerbatimNode(parent, items)
	parent.Append(node)
	p.walkElements(node, items[current+1:end])
	return end
}

// isVerbatim returns true if a collection of items matches verbatim markup
func isVerbatim(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsEqual()
	})
}

// findVerbatim end item of a verbatim collection
func findVerbatim(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsEqual()
	})
}

func (p *parser) newStrikeThrough(parent ast.Node, items []lex.Item) (end int) {
	current := 0
	end = findStrikeThrough(items)
	node := ast.NewStrikeThroughNode(parent, items)
	parent.Append(node)
	p.walkElements(node, items[current+1:end])
	return end
}

// isStrikeThrough returns true if a collection of items matches strike through markup
func isStrikeThrough(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsPlus()
	})
}

// findStrikeThrough end item of a strike through collection
func findStrikeThrough(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsPlus()
	})
}

func (p *parser) newUnderline(parent ast.Node, items []lex.Item) (end int) {
	current := 0
	end = findUnderline(items)
	node := ast.NewUnderlineNode(parent, items)
	parent.Append(node)
	p.walkElements(node, items[current+1:end])
	return end
}

// isUnderline returns true if a collection of items matches underline markup
func isUnderline(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsUnderscore()
	})
}

// findUnderline end item of a underline collection
func findUnderline(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsUnderscore()
	})
}

func (p *parser) newCode(parent ast.Node, items []lex.Item) (end int) {
	current := 0
	end = findCode(items)
	node := ast.NewCodeNode(parent, items)
	parent.Append(node)
	p.walkElements(node, items[current+1:end])
	return end
}

// isCode returns true if a collection of items matches code markup
func isCode(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsTilde()
	})
}

// findCode end item of a code collection
func findCode(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsTilde()
	})
}
