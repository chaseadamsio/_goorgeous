package goorgeous

// a block is a headline if it starts with 1-6 asterisks followed by a space
func isHeadline(items []item, token item) bool {
	// first item has to be an asterisk
	if len(items) > 0 && token.typ != itemAsterisk { // true and false = false
		return false
	}
	for idx := 0; idx <= 5; idx++ { // idx = 1
		if len(items) <= idx { // 22 < 1 => false
			return false
		}
		if items[idx].typ == itemAsterisk { // " "
			continue
		}
		if items[idx].typ == itemSpace { // true
			return true
		}
		return false
	}
	return false
}

func headlineDepth(items []item) int {
	depth := 0
	for depth <= 6 { // 1
		if len(items) > depth && items[depth].typ == itemAsterisk { //
			depth++
			continue
		}
		if len(items) > depth && items[depth].typ == itemSpace {
			return depth
		}
	}
	return depth
}
