package parse

import "github.com/chaseadamsio/goorgeous/lex"

func isKeyword(items []lex.Item) bool {
	current := 0
	token := items[0]
	itemsLength := len(items)

	if !token.IsHash() {
		return false
	}
	if current < itemsLength && items[current+1].Type() != lex.ItemPlus {
		return false
	}
	if current == 0 || items[current-1].IsNewline() {
		return true
	}
	return false
}

func findGreaterBlock(items []lex.Item) (found bool, end int) {
	current := 0
	itemsLength := len(items)
	foundEnd := false

	if !isKeyword(items) {
		return false, -1
	}

	if itemsLength > current+1 && items[current+2].Value() != "BEGIN" {
		return false, -1
	}
	current = current + 2

	for current < itemsLength {
		if foundEnd && (current+1 == itemsLength || items[current].IsNewline() || items[current].IsEOF()) {
			return true, itemsLength
		}
		if items[current].Type() == lex.ItemHash {
			if itemsLength > current && items[current+1].Type() == lex.ItemPlus {
				if itemsLength > current+1 && items[current+2].Value() == "END" {
					foundEnd = true
				}
			}
		}
		current++
	}
	return false, -1
}
