package parse

import "github.com/chaseadamsio/goorgeous/lex"

func isUnorderedList(token lex.Item, items []lex.Item, current int) bool {
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

func findUnorderedList(items []lex.Item) int {
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		token := items[current]
		if token.IsNewline() {
			if !isUnorderedList(items[current+1], items, current+1) {
				return current
			}
		}
		current++
	}
	return itemsLength
}

func isOrderedList(token lex.Item, items []lex.Item, current int) bool {
	return false
}

func findOrderedList(items []lex.Item) int {
	return 0
}
