package tokens

import "github.com/chaseadamsio/goorgeous/lex"

const maxHeadlineDepth = 6

// IsHeadline determines inf a collection of items is a headline
func IsHeadline(items []lex.Item) bool {
	current := 0
	itemsLength := len(items)
	token := items[0]
	// first item has to be an asterisk
	if !token.IsAsterisk() {
		return false
	}
	for current <= maxHeadlineDepth && current < itemsLength {
		currItem := items[current]
		// it's still a potential heading
		if currItem.IsAsterisk() {
			current++
			continue
		}
		// space terminates the headline "stars"
		if currItem.IsSpace() {
			return true
		}
		return false
	}
	return false
}

// HeadlineDepth determines the depth of a headline
func HeadlineDepth(items []lex.Item) int {
	depth := 0
	itemsLength := len(items)
	for depth <= maxHeadlineDepth {
		hasNextItem := itemsLength > depth
		currItem := items[depth]
		if hasNextItem && currItem.IsAsterisk() {
			depth++
			continue
		}
		if hasNextItem && currItem.IsSpace() {
			return depth
		}
	}
	return depth
}
