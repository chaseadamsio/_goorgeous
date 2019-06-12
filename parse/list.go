package parse

import (
	"regexp"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

// checks if the current collection of items is a list by checking
// it against a provided list function.
func isList(items []lex.Item, isListFunc func([]lex.Item) bool) (start int, found bool) {
	itemsLength := len(items)
	offset := 1 // we want to check if previous character is a newline
	current := 0

	for current < itemsLength {
		if !(items[current].IsTab() || items[current].IsSpace()) {
			break
		}
		current++
		offset++
	}

	if isListFunc(items[current:]) {

		if current != 0 {

			for current >= 0 {
				if current < offset || items[current-offset].IsNewline() {
					return current, true
				}
				current--
			}
		}

		if current < itemsLength && items[current+1].IsWhitespace() {
			return current, true
		}

		return -1, false

	}
	return -1, false
}

func isUnorderedList(items []lex.Item) (start int, found bool) {
	return isList(items, foundUnorderedListItem)
}

func isOrderedList(items []lex.Item) (start int, found bool) {
	return isList(items, foundOrderedListItem)
}

// findList takes a collection of Item and a function to check a collection
// against and returns the end of the list
func findList(items []lex.Item, isListFunc func(items []lex.Item) (start int, found bool)) int {
	current := 0
	itemsLength := len(items)
	baseIndentLevel, _ := getIndentLevel(items[current:])
	for current < itemsLength {
		if _, found := foundTerminatingNewline(items[current:], baseIndentLevel, isListFunc); found {
			if current+1 < itemsLength {
				if _, found := isListFunc(items[current+1:]); !found {
					return current
				}
			}
		}
		current++
	}
	return current
}

func maybeList(items []lex.Item) (listTyp string, start, end int, found bool) {
	if _, found = isOrderedList(items); found {
		listTyp = "ORDERED"
		start, end = findOrderedListBoundaries(items)
	} else if _, found = isUnorderedList(items); found {
		listTyp = "UNORDERED"
		start, end = findUnorderedListBoundaries(items)
	}
	return listTyp, start, end, found
}

func findOrderedListBoundaries(items []lex.Item) (start, end int) {
	start, end = findListBoundaries(items, isOrderedList)
	return start, end
}

func findUnorderedListBoundaries(items []lex.Item) (start, end int) {
	start, end = findListBoundaries(items, isUnorderedList)
	return start, end
}

// findList takes a collection of Item and a function to check a collection
// against and returns the end of the list
func findListBoundaries(items []lex.Item, isListFunc func(items []lex.Item) (start int, found bool)) (start, end int) {
	current := 0
	startFound := false
	itemsLength := len(items)
	baseIndentLevel, _ := getIndentLevel(items[current:])
	for current < itemsLength {
		if _, found := isListFunc(items[current:]); startFound && found { // this *should* always be 0.
			start = current
			startFound = true
		} else if _, found := foundTerminatingNewline(items[current:], baseIndentLevel, isListFunc); found {
			if current+1 < itemsLength {
				current++
				if _, found := isListFunc(items[current:]); !found {
					end = current
					return start, end
				}
			}
		}
		current++
	}
	end = current
	return start, end
}

func findUnorderedList(items []lex.Item) int {
	return findList(items, isUnorderedList)
}

func findOrderedList(items []lex.Item) int {
	return findList(items, isOrderedList)
}

var findItemFuncMap = map[string]func(items []lex.Item) (start, end, foundNestedListStart, foundNestedListEnd int){
	"UNORDERED": findUnorderedListItem,
	"ORDERED":   findOrderedListItem,
}

func (p *parser) makeList(listTyp string, parent ast.Node, items []lex.Item) (end int) {
	start, current := 0, 0
	itemsLength := len(items)
	findFunc := findItemFuncMap[listTyp]

	listNode := ast.NewListNode(listTyp, parent, items)
	parent.Append(listNode)

	for current < itemsLength {

		foundStart, end, foundNestedListStart, foundNestedListEnd := findFunc(items[current:])
		end = current + end
		node := ast.NewListItemNode(parent, items[start+foundStart:end])
		listNode.Append(node)

		if foundNestedListStart > 0 {
			foundNestedListStart = start + foundNestedListStart
			foundNestedListEnd = start + foundNestedListEnd
			p.walk(node, items[foundNestedListStart:foundNestedListEnd])

			current = foundNestedListEnd
		} else {
			current = end + 1 // skip the newline
		}
		start = current
	}
	return current
}

func foundTerminatingNewline(items []lex.Item, indentLevel int,
	foundMatchFunc func(items []lex.Item) (start int, found bool)) (offset int, found bool) {
	itemsLength := len(items)
	if items[0].IsNewline() { // found the first new line in example

		if (1 < itemsLength && items[1].IsTab()) || (1 < itemsLength && items[1].IsSpace()) {
			if currIndentLevel, _ := getIndentLevel(items[1:]); currIndentLevel < indentLevel {
				return 0, true
			}
			return 0, false
		} else if indentLevel > 0 { // no tab or space, but indentlevel means the list is closed
			return 0, true
		}

		if 1 < itemsLength {
			if _, found := foundMatchFunc(items[1:]); found {
				return 0, false
			}
		}

		if 1 < itemsLength && items[1].IsWord() || 1 < itemsLength && items[1].IsNonWord() {
			return 0, true
		}

		if itemsLength == 1 {
			return 0, true
		}

		if 1 < itemsLength && items[1].IsEOF() {
			return 0, true
		}

		if 1 < itemsLength && items[1].IsNewline() { // found the second newline
			if 2 < itemsLength && items[2].IsNewline() {
				return 0, true
			}
		}
	}
	return 0, false
}

func foundListItemTerminatingNewline(items []lex.Item,
	foundMatchFunc func(items []lex.Item) bool) (offset int, found bool) {

	current := 0
	itemsLength := len(items)

	if items[current].IsNewline() { // found the first new line in example
		current++
		if current < itemsLength && foundMatchFunc(items[current:]) {
			return current, true
		}

		if current < itemsLength && (items[current].IsTab() || items[current].IsSpace()) {
			for current < itemsLength {
				if !(items[current].IsTab() || items[current].IsSpace()) {
					break
				}
				current++
			}
			if current < itemsLength && (foundOrderedListItem(items[current:]) || foundUnorderedListItem(items[current:])) {
				return 0, false
			}
			return 0, false
		}
		return current, true
	}
	return 0, false

}

func getIndentLevel(items []lex.Item) (indentLevel, itemOffset int) {
	itemsLength := len(items)
	foundIndentLevel := 0
	tabIndentSize := 4 // tabs have to be equal to spaces and this is the best way to handle it
	spaceIndentSize := 1
	for foundIndentLevel < itemsLength {
		if !(items[foundIndentLevel].IsTab() || items[foundIndentLevel].IsSpace()) {
			break
		}
		if items[foundIndentLevel].IsTab() {
			itemOffset++
			foundIndentLevel = foundIndentLevel + tabIndentSize
		}
		if items[foundIndentLevel].IsSpace() {
			itemOffset++
			foundIndentLevel = foundIndentLevel + spaceIndentSize
		}
	}
	return foundIndentLevel, itemOffset
}

var isListFuncs = map[string]func([]lex.Item) (start int, found bool){
	"UNORDERED": isUnorderedList,
	"ORDERED":   isOrderedList,
}

// findListItem returns
// - start (start being the first character after the (bullet + N space character)
// - end (the character before the newline)
// - next (the next character after the newline. either after end or foundNestedListEnd)
// - foundNestedListStart (if there's a nested list, the beginning [greater than 0])
// - foundNestedListEnd (if there's a nested list, the end of the nested list [greater than 0])
func findListItem(items []lex.Item,
	foundMatchFunc func(items []lex.Item) bool) (start, trueEnd, foundNestedListStart, foundNestedListEnd int) {

	itemsLength := len(items)
	current := 0
	end := 0

	offset := 0
	baseIndentLevel := 0
	indentLevel := 0
	foundNestedListStart = 0

	if items[current].IsTab() || items[current].IsSpace() {
		baseIndentLevel, offset = getIndentLevel(items[current:])
		indentLevel = baseIndentLevel
		current = current + offset
	}

	for current < itemsLength {
		if items[current].IsNewline() {
			if current+1 < itemsLength && (items[current+1].IsTab() || items[current+1].IsSpace()) {
				foundIndentLevel, _ := getIndentLevel(items[current+1:])
				if indentLevel < foundIndentLevel {
					if listTyp, found := foundListItem(items[current+1:]); found {
						end = current // don't send the newline for closing
						current++     // don't send the newline for opening a new list
						nestedStart, nestedEnd := findListBoundaries(items[current:], isListFuncs[listTyp])
						return start, end, current + nestedStart, current + nestedEnd

					}
				} else if indentLevel == foundIndentLevel {
					end = current // don't send the newline for closing
					return start, end, 0, 0
				}
			}
		}

		if indentLevel == baseIndentLevel && current > 0 {
			if _, found := foundListItemTerminatingNewline(items[current:], foundMatchFunc); found {
				return start, current, foundNestedListStart, foundNestedListEnd
			}
		}

		if foundMatchFunc(items[current:]) {
			start = current
		}
		end++
		current++
	}
	return start, end + 1, foundNestedListStart, foundNestedListEnd
}

func findOrderedListItem(items []lex.Item) (start, end, foundNestedListStart, foundNestedListEnd int) {
	start, end, foundNestedListStart, foundNestedListEnd = findListItem(items, foundOrderedListItem)
	return start, end, foundNestedListStart, foundNestedListEnd
}

func foundListItem(items []lex.Item) (listTyp string, found bool) {
	if found = foundOrderedListItem(items); found {
		listTyp = "ORDERED"
	} else if found = foundUnorderedListItem(items); found {
		listTyp = "UNORDERED"
	}
	return listTyp, found
}

func foundOrderedListItem(items []lex.Item) bool {
	current := 0
	for current < len(items) {
		if !(items[current].IsTab() || items[current].IsSpace()) {
			break
		}
		current++
	}
	matched, err := regexp.Match(`^\d+\.$`, []byte(items[current].Value()))
	if err != nil {
		panic(err)
	}
	return items[current].Type() == lex.ItemText && matched
}

func findUnorderedListItem(items []lex.Item) (start, end, foundNestedListStart, foundNestedListEnd int) {
	start, end, foundNestedListStart, foundNestedListEnd = findListItem(items, foundUnorderedListItem)
	return start, end, foundNestedListStart, foundNestedListEnd
}

func foundUnorderedListItem(items []lex.Item) bool {
	current := 0
	for current < len(items) {
		if !(items[current].IsTab() || items[current].IsSpace()) {
			break
		}
		current++
	}
	return items[current].Type() == lex.ItemDash && (current+1 < len(items) && items[current+1].IsSpace())
}
