package tokens

import "github.com/chaseadamsio/goorgeous/lex"

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

// IsBold returns true if a collection of items matches bold markup
func IsBold(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsAsterisk()
	})
}

// FindBold finds the end item of a bold collection
func FindBold(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsAsterisk()
	})
}

// IsItalic returns true if a collection of items matches italic markup
func IsItalic(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsForwardSlash()
	})
}

// FindItalic finds the end item of a italic collection
func FindItalic(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsForwardSlash()
	})
}

// IsVerbatim returns true if a collection of items matches verbatim markup
func IsVerbatim(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsEqual()
	})
}

// FindVerbatim end item of a verbatim collection
func FindVerbatim(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsEqual()
	})
}

// IsStrikeThrough returns true if a collection of items matches strike through markup
func IsStrikeThrough(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsPlus()
	})
}

// FindStrikeThrough end item of a strike through collection
func FindStrikeThrough(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsPlus()
	})
}

// IsUnderline returns true if a collection of items matches underline markup
func IsUnderline(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsUnderscore()
	})
}

// FindUnderline end item of a underline collection
func FindUnderline(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsUnderscore()
	})
}

// IsCode returns true if a collection of items matches code markup
func IsCode(items []lex.Item) bool {
	return isElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsTilde()
	})
}

// FindCode end item of a code collection
func FindCode(items []lex.Item) int {
	return findElementMarkup(items, func(currToken lex.Item) bool {
		return currToken.IsTilde()
	})
}
