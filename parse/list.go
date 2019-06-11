package parse

import (
	"regexp"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func isList(items []lex.Item, isListFunc func([]lex.Item) bool) bool {
	itemsLength := len(items)
	offset := 1 // we want to check if previous character is a newline
	current := 0
	for items[current].IsTab() {
		current++
		offset++
	}
	if isListFunc(items[current:]) {

		if current != 0 {

			for current >= 0 {
				if current < offset || items[current-offset].IsNewline() {
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

func isUnorderedList(items []lex.Item) bool {
	return isList(items, foundUnorderedListItem)
}

func isOrderedList(items []lex.Item) bool {
	return isList(items, foundOrderedListItem)
}

// findList takes a collection of Item and a function to check a collection
// against and returns the end of the list
func findList(items []lex.Item, isListFunc func(items []lex.Item) bool) int {
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		if foundTerminatingNewline(items[current:], isListFunc) {
			if current+1 < itemsLength && !isListFunc(items[current+1:]) {
				return current
			}
		}
		current++
	}
	return itemsLength
}

func findUnorderedList(items []lex.Item) int {
	return findList(items, isUnorderedList)
}

func findOrderedList(items []lex.Item) int {
	return findList(items, isOrderedList)
}

func (p *parser) makeListItems(parent ast.Node, items []lex.Item,
	findFunc func(items []lex.Item) (foundStart, foundEnd, foundNestedListStart, foundNestedListEnd int)) (end int) {
	start, current := 0, 0
	itemsLength := len(items)
	for current < itemsLength {
		foundStart, foundEnd, foundNestedListStart, foundNestedListEnd := findFunc(items[current:itemsLength])
		start = start + foundStart
		foundEnd = current + foundEnd
		node := ast.NewListItemNode(parent, items[start:foundEnd])
		if foundNestedListStart > 0 {
			p.walk(node, items[foundNestedListStart:foundNestedListEnd])
		}
		parent.Append(node)
		current = foundEnd
		start = current
	}
	return current
}

func (p *parser) makeUnorderedListItems(parent ast.Node, items []lex.Item) (end int) {
	return p.makeListItems(parent, items, findUnorderedListItem)
}

func (p *parser) makeOrderedListItems(parent ast.Node, items []lex.Item) (end int) {
	return p.makeListItems(parent, items, findOrderedListItem)
}

// org mode is super weird in that this is a valid list:
// - foo
//
// - bar
// so we need to check if the next three characters are Newlines
func foundTerminatingNewline(items []lex.Item,
	foundMatchFunc func(items []lex.Item) bool) bool {
	itemsLength := len(items)
	if items[0].IsNewline() { // found the first new line in example

		if 1 < itemsLength && items[1].IsWord() {
			return true
		}

		if 1 < itemsLength && foundMatchFunc(items[1:]) {
			return false
		}

		if 1 < itemsLength && (items[1].IsTab()) {
			return false
		}

		if itemsLength == 1 {
			return true
		}

		if 1 < itemsLength && items[1].IsEOF() {
			return true
		}

		if 1 < itemsLength && items[1].IsNewline() { // found the second newline
			if 2 < itemsLength && items[2].IsNewline() {
				return true
			}
		}
	}
	return false
}

func foundListItemTerminatingNewline(items []lex.Item,
	foundMatchFunc func(items []lex.Item) bool) bool {

	current := 0
	itemsLength := len(items)

	if items[current].IsNewline() { // found the first new line in example
		current++
		if current < itemsLength && foundMatchFunc(items[current:]) {
			return true
		}

		if current < itemsLength && (items[current].IsTab()) {
			for current < itemsLength {
				if !items[current].IsTab() {
					break
				}
				current++
			}
			if current < itemsLength && foundMatchFunc(items[current:]) {
				return true
			}
			return false
		}
		return true
	}
	return false

}

func isDeeperIndent(currIndentLevel int, items []lex.Item) bool {
	itemsLength := len(items)
	foundIndentLevel := 0
	for foundIndentLevel < itemsLength {
		if items[foundIndentLevel].Type() != lex.ItemTab {
			break
		}
		foundIndentLevel++
	}
	return currIndentLevel < foundIndentLevel
}

func findListItem(items []lex.Item,
	foundMatchFunc func(items []lex.Item) bool) (start, trueEnd, foundNestedListStart, foundNestedListEnd int) {
	itemsLength := len(items)
	current := 0
	end := 0
	trueEnd = itemsLength
	indentBaseLevel := -1
	indentLevel := 0
	foundNestedListStart = 0

	if items[current].IsTab() {
		indentBaseLevel = 1
		current = current + 1
		indentLevel = indentBaseLevel
		for current < itemsLength {
			if !items[current+1].IsTab() {
				break
			}
			indentBaseLevel++
			current++
		}
	}

	for current < itemsLength {
		// if indentBaseLevel == indentLevel && items[current].IsNewline() { // this will fix the panic
		// if items[current].IsNewline() {
		// 	if current < itemsLength && items[current+1].IsTab() {
		// 		if !isDeeperIndent(indentLevel, items[current:]) {
		// 			indentLevel++
		// 		}
		// 		if foundNestedListStart == 0 {
		// 			foundNestedListStart = current
		// 		}
		// 	} else {
		// 		foundNestedListEnd = current
		// 		indentLevel = 0
		// 	}
		// }

		if indentLevel == 0 && (current > 0 && foundListItemTerminatingNewline(items[current:], foundMatchFunc)) {
			return start, current, foundNestedListStart, foundNestedListEnd
		}

		if foundMatchFunc(items[current:]) {
			start = end
		}

		end++
		current++
	}
	return start, current, foundNestedListStart, foundNestedListEnd
}

func findOrderedListItem(items []lex.Item) (start, end, foundNestedListStart, foundNestedListEnd int) {
	start, end, foundNestedListStart, foundNestedListEnd = findListItem(items, foundOrderedListItem)
	return start, end, foundNestedListStart, foundNestedListEnd
}

func foundOrderedListItem(items []lex.Item) bool {
	matched, err := regexp.Match(`^\d+\.$`, []byte(items[0].Value()))
	if err != nil {
		panic(err)
	}
	return items[0].Type() == lex.ItemText && matched
}

func findUnorderedListItem(items []lex.Item) (start, end, foundNestedListStart, foundNestedListEnd int) {
	start, end, foundNestedListStart, foundNestedListEnd = findListItem(items, foundUnorderedListItem)
	return start, end, foundNestedListStart, foundNestedListEnd
}

func foundUnorderedListItem(items []lex.Item) bool {
	return items[0].Type() == lex.ItemDash && (1 < len(items) && items[1].IsSpace())
}
