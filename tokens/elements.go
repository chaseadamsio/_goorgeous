package tokens

import "github.com/chaseadamsio/goorgeous/lex"

func isElementMarkup(items []lex.Item, expectedTypeFunc func(lex.Item) bool) bool {
	foundOpeningChar := false
	itemsLength := len(items)
	for current, currToken := range items {
		// check the current token against the expected type checker function
		if expectedTypeFunc(currToken) {
			nextTokenIsWhitespace := current < itemsLength && items[current+1].IsWhitespace()

			// a match for opening character, it cannot precede a whitespace character
			if !foundOpeningChar && !nextTokenIsWhitespace {
				foundOpeningChar = true
				continue
			}
			// closing characters cannot follow a whitespace character
			if foundOpeningChar && !items[current-1].IsWhitespace() {
				// if it precedes EOF, Newline and Whitespace by this point, it's a match
				if itemsLength > current &&
					(items[current+1].IsEOF() || items[current+1].IsNewline() || items[current+1].IsWhitespace()) {
					return true
				}
			}
			// if it's a newline, the first character or the end of the collection, we didn't find the expected type
		} else if currToken.IsNewline() || current == 0 || current == itemsLength {
			return false
		}
	}
	return false
}

// IsBold returns true if a collection of items matches bold markup
func IsBold(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsAsterisk()
	})
}

// IsItalic returns true if a collection of items matches italic markup
func IsItalic(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsForwardSlash()
	})
}

// IsVerbatim returns true if a collection of items matches verbatim markup
func IsVerbatim(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsEqual()
	})
}

// IsStrikeThrough returns true if a collection of items matches strike through markup
func IsStrikeThrough(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsPlus()
	})
}

// IsUnderline returns true if a collection of items matches underline markup
func IsUnderline(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsUnderscore()
	})
}

// IsCode returns true if a collection of items matches code markup
func IsCode(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsTilde()
	})
}
