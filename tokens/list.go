package tokens

import "github.com/chaseadamsio/goorgeous/lex"

func IsUnorderedList(token lex.Item, items []lex.Item, current int) bool {
	itemsLength := len(items)
	if token.Type() == lex.ItemDash {

		if current != 0 {

			for current >= 0 {
				if current > 0 && !(items[current-1].IsWhitespace() || items[current-1].IsNewline()) {
					return false
				} else if items[current-1].IsNewline() {
					return true
				}
				current--
			}
		}

		if current < itemsLength && items[current+1].IsWhitespace() {
			return true
		}

		return false

	}
	return false
}

func FindUnorderedList(items []lex.Item) int {
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		token := items[current]
		if token.IsNewline() {
			if !IsUnorderedList(items[current+1], items, current+1) {
				return current
			}
		}
		current++
	}
	return itemsLength
}

func IsOrderedList(token lex.Item, items []lex.Item, current int) bool {
	return false
}

func FindOrderedList(items []lex.Item) int {
	return 0
}