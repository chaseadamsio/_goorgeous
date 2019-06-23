package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
)

func (p *parser) makeOrderedList(parent ast.Node, current, end int) {
	p.makeList("ORDERED", parent, current, end)
}

func (p *parser) makeUnorderedList(parent ast.Node, current, end int) {
	p.makeList("UNORDERED", parent, current, end)
}

func (p *parser) makeList(listTyp string, parent ast.Node, current, end int) {

	findItemFuncMap := map[string]func(current, listEnd int) (start, end, foundNestedListStart, foundNestedListEnd int){
		"UNORDERED": p.findUnorderedListItem,
		"ORDERED":   p.findOrderedListItem,
	}
	findFunc := findItemFuncMap[listTyp]

	listNode := ast.NewListNode(listTyp, parent, p.items[current:end])
	parent.Append(listNode)

	for current < end {

		foundStart, foundEnd, foundNestedListStart, foundNestedListEnd := findFunc(current, end)
		node := ast.NewListItemNode(parent, p.items[foundStart:foundEnd])
		listNode.Append(node)

		if foundNestedListStart > 0 {
			p.walk(node, foundNestedListStart, foundNestedListEnd)

			current = foundNestedListEnd
		} else {
			current = foundEnd + 1
		}
		// current++ // this is bumping it forward when it's nested
	}
}

func (p *parser) matchesOrderedList(current int) (found bool, end int) {
	return p.matchesList(current, p.foundOrderedListItem)
}

func (p *parser) matchesUnorderedList(current int) (found bool, end int) {
	return p.matchesList(current, p.foundUnorderedListItem)
}

// checks if the current collection of items is a list by checking
// it against a provided list function.
func (p *parser) matchesList(current int, isListFunc func(current int) bool) (found bool, end int) {
	itemsLength := len(p.items)

	if 0 < current && !p.items[current-1].IsNewline() {
		return false, -1
	}

	if isListFunc(current) {
		_, end := p.findListBoundaries(current, isListFunc)

		if current < itemsLength && p.items[current+1].IsWhitespace() {
			return true, end
		}
		return false, -1

	}
	return false, -1
}

// findListBoundaries takes a place in a collection of items and a function to check a collection
// against and returns the end of the list
func (p *parser) findListBoundaries(current int, isListFunc func(current int) (found bool)) (start, end int) {
	startFound := false
	itemsLength := len(p.items)
	baseIndentLevel, _ := p.getIndentLevel(current)
	for current < itemsLength {
		if found := isListFunc(current); !startFound && found { // this *should* always be 0.
			start = current
			startFound = true
		} else if found, _ := p.foundTerminatingNewline(current, baseIndentLevel, isListFunc); found {
			if current+1 < itemsLength {
				current++
				currIndentLevel, _ := p.getIndentLevel(current)
				if found := isListFunc(current); !found ||
					(found && currIndentLevel < baseIndentLevel) { // could match list type but be diff. level of nesting
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

func (p *parser) getIndentLevel(current int) (indentLevel, itemOffset int) {
	itemsLength := len(p.items)
	itemOffset = current
	foundIndentLevel := 0
	tabIndentSize := 4 // tabs have to be equal to spaces and this is the best way to handle it
	spaceIndentSize := 1
	for current < itemsLength {
		if !(p.items[current].IsTab() || p.items[current].IsSpace()) {
			break
		}
		if p.items[current].IsTab() {
			itemOffset++
			foundIndentLevel = foundIndentLevel + tabIndentSize
		}
		if p.items[current].IsSpace() {
			itemOffset++
			foundIndentLevel = foundIndentLevel + spaceIndentSize
		}
		current++
	}
	return foundIndentLevel, itemOffset
}

func (p *parser) foundTerminatingNewline(current int, indentLevel int,
	foundMatchFunc func(current int) (found bool)) (found bool, offset int) {
	itemsLength := len(p.items)
	if p.items[current].IsNewline() && current+1 < itemsLength { // found the first new line in example
		next := current + 1
		currIndentLevel, _ := p.getIndentLevel(next)

		if p.items[next].IsTab() || p.items[next].IsSpace() {
			if currIndentLevel < indentLevel {
				return true, 0
			}
			return false, -1
			// TODO handle if curr indent level < indentLevel
		} else if indentLevel > 0 { // no tab or space, but indentlevel means the list is closed
			return true, 0
		}

		if currIndentLevel < indentLevel {
			return true, 0
		}

		if found := foundMatchFunc(next); found {
			return false, -1
		}

		if p.items[next].IsWord() || p.items[next].IsNonWord() {
			return true, 0
		}

		if p.items[next].IsEOF() {
			return true, 0
		}

		if p.items[next].IsNewline() { // found the second newline
			next++
			if next < itemsLength && p.items[next].IsNewline() {
				return true, 0
			}
		}
	} else if itemsLength == 1 {
		return true, 0
	}

	return false, 0
}
