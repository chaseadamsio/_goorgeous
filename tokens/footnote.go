package tokens

import "github.com/chaseadamsio/goorgeous/lex"

func IsFootnoteDefinition(token lex.Item, items []lex.Item, current int) bool {
	itemsLength := len(items)
	if current > 0 && !items[current-1].IsNewline() {
		return false
	}

	if !(token.IsBracket() && token.Value() == "[") {
		return false
	}

	if current < itemsLength && items[current+1].Value() != "fn" {
		return false
	}

	if current+1 < itemsLength && !items[current+2].IsColon() {
		return false
	}

	current = current + 2

	foundFootnoteLabel := false
	foundFootnoteInlineDefinition := false // when there's a definition, we can allow spaces

	for current < itemsLength {
		if !foundFootnoteInlineDefinition && items[current].IsSpace() {
			return false
		}
		if items[current].IsNewline() || items[current].IsEOF() {
			return false
		}
		if items[current].IsColon() {
			foundFootnoteInlineDefinition = true
		}
		if foundFootnoteLabel && items[current].IsBracket() && items[current].Value() == "]" {
			return true
		}
		foundFootnoteLabel = true
		current++
	}

	return false
}

func FindFootnoteDefinition(items []lex.Item) int {
	current := 0
	itemsLength := len(items)

	for current < itemsLength {
		currItem := items[current]
		if currItem.IsBracket() && currItem.Value() == "]" {
			return current
		}
		current++
	}
	return itemsLength
}
