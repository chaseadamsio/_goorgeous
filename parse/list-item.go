package parse

import (
	"regexp"

	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) foundOrderedListItem(current int) bool {
	for current < len(p.items) {
		if !(p.items[current].IsTab() || p.items[current].IsSpace()) {
			break
		}
		current++
	}
	matched, err := regexp.Match(`^\d+\.$`, []byte(p.items[current].Value()))
	if err != nil {
		panic(err)
	}
	return p.items[current].Type() == lex.ItemText && matched
}

func (p *parser) foundUnorderedListItem(current int) bool {
	for current < len(p.items) {
		if !(p.items[current].IsTab() || p.items[current].IsSpace()) {
			break
		}
		current++
	}
	return p.items[current].Type() == lex.ItemDash && (current+1 < len(p.items) && p.items[current+1].IsSpace())
}

func (p *parser) foundListItemTerminatingNewline(current int,
	foundMatchFunc func(current int) bool) (found bool, offset int) {

	itemsLength := len(p.items)

	if p.items[current].IsNewline() { // found the first new line in example
		current++
		if current < itemsLength && foundMatchFunc(current) {
			return true, current
		}

		if current < itemsLength && (p.items[current].IsTab() || p.items[current].IsSpace()) {
			for current < itemsLength {
				if !(p.items[current].IsTab() || p.items[current].IsSpace()) {
					break
				}
				current++
			}
			if current < itemsLength && (p.foundOrderedListItem(current) || p.foundUnorderedListItem(current)) {
				return false, 0
			}
			return false, 0
		}
		return true, current
	}
	return false, 0

}

func (p *parser) findOrderedListItem(current int) (start, end, foundNestedListStart, foundNestedListEnd int) {
	start, end, foundNestedListStart, foundNestedListEnd = p.findListItem(current, p.foundOrderedListItem)
	return start, end, foundNestedListStart, foundNestedListEnd
}

func (p *parser) findUnorderedListItem(current int) (start, end, foundNestedListStart, foundNestedListEnd int) {
	start, end, foundNestedListStart, foundNestedListEnd = p.findListItem(current, p.foundUnorderedListItem)
	return start, end, foundNestedListStart, foundNestedListEnd
}

// findListItem returns
// - start (start being the first character after the (bullet + N space character)
// - end (the character before the newline)
// - next (the next character after the newline. either after end or foundNestedListEnd)
// - foundNestedListStart (if there's a nested list, the beginning [greater than 0])
// - foundNestedListEnd (if there's a nested list, the end of the nested list [greater than 0])
func (p *parser) findListItem(current int,
	foundMatchFunc func(current int) bool) (start, trueEnd, foundNestedListStart, foundNestedListEnd int) {

	var matchesListFuncs = map[string]func(current int) (found bool){
		"UNORDERED": p.foundUnorderedListItem,
		"ORDERED":   p.foundOrderedListItem,
	}

	itemsLength := len(p.items)
	end := current

	offset := current
	baseIndentLevel := 0
	indentLevel := 0
	foundNestedListStart = 0

	if p.items[current].IsTab() || p.items[current].IsSpace() {
		baseIndentLevel, offset = p.getIndentLevel(current)
		indentLevel = baseIndentLevel
		current = offset
	}

	for current < itemsLength {
		if p.items[current].IsNewline() {
			if current+1 < itemsLength && (p.items[current+1].IsTab() || p.items[current+1].IsSpace()) {
				foundIndentLevel, _ := p.getIndentLevel(current + 1)
				if indentLevel < foundIndentLevel {
					if found, listTyp := p.foundListItem(current + 1); found {
						end = current // don't send the newline for closing
						current++     // don't send the newline for opening a new list
						nestedStart, nestedEnd := p.findListBoundaries(current, matchesListFuncs[listTyp])
						return start, end, nestedStart, nestedEnd

					}
				} else if indentLevel == foundIndentLevel {
					end = current // don't send the newline for closing
					return start, end, 0, 0
				}
			}
		}

		if indentLevel == baseIndentLevel && current > 0 {
			if found, _ := p.foundListItemTerminatingNewline(current, foundMatchFunc); found {
				return start, current, foundNestedListStart, foundNestedListEnd
			}
		}

		if foundMatchFunc(current) {
			start = current
		}
		end++
		current++
	}
	return start, end, foundNestedListStart, foundNestedListEnd
}

func (p *parser) foundListItem(current int) (found bool, listTyp string) {
	if found = p.foundOrderedListItem(current); found {
		listTyp = "ORDERED"
	} else if found = p.foundUnorderedListItem(current); found {
		listTyp = "UNORDERED"
	}
	return found, listTyp
}
