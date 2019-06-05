package goorgeous

// a block is a headline if it starts with 1-6 asterisks followed by a space
func isHeadline(items []item) bool {
	// first item has to be an asterisk
	if len(items) > 0 && items[0].typ != itemAsterisk {
		return false
	}
	for idx := 1; idx <= 6; idx++ {
		if len(items) < idx {
			return false
		}
		if items[idx].typ == itemAsterisk {
			continue
		}
		if items[idx].typ == itemSpace {
			return true
		}
		return false
	}
	return false
}

func headlineDepth(items []item) int {
	depth := 0
	for idx := 0; idx <= 6; idx++ {
		if items[idx].typ == itemAsterisk {
			depth++
		}
		if items[idx].typ == itemSpace {
			return depth
		}
	}
	return depth
}
