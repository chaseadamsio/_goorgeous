package tokens

import "github.com/chaseadamsio/goorgeous/lex"

// a block is a headline if it starts with 1-6 asterisks followed by a space
func IsHeadline(items []lex.Item) bool {
	token := items[0]
	// first item has to be an asterisk
	if len(items) > 0 && token.Type() != lex.ItemAsterisk { // true and false = false
		return false
	}
	for idx := 0; idx <= 6; idx++ { // idx = 1
		if len(items) <= idx { // 22 < 1 => false
			return false
		}
		if items[idx].IsAsterisk() { // " "
			continue
		}
		if items[idx].Type() == lex.ItemSpace { // true
			return true
		}
		return false
	}
	return false
}

func HeadlineDepth(items []lex.Item) int {
	depth := 0
	for depth <= 6 { // 1
		if len(items) > depth && items[depth].Type() == lex.ItemAsterisk { //
			depth++
			continue
		}
		if len(items) > depth && items[depth].Type() == lex.ItemSpace {
			return depth
		}
	}
	return depth
}
