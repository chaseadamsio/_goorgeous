package tokens

import "github.com/chaseadamsio/goorgeous/lex"

// IsLink checks if the items it receives matches a link
// this is based on the docs from:
//		https://orgmode.org/manual/Link-format.html#Link-format
func IsLink(items []lex.Item) bool {
	current := 0
	itemsLength := len(items)

	// a link will always start with two brackets [[
	if !((items[current].IsBracket() && items[current].Value() == "[") &&
		(current < itemsLength && items[current+1].IsBracket() && items[current+1].Value() == "[")) {

		return false
	}
	current = current + 2

	foundLinkCloseBracket := false

	for current < itemsLength {
		// we only care about brackets if we've closed the link part
		if items[current].IsBracket() && (foundLinkCloseBracket || items[current].Value() == "]") {
			// we found a ] to get here and if the next bracket is a closing bracket, this is a link!
			if current < itemsLength && items[current+1].IsBracket() && items[current+1].Value() == "]" {
				return true
			}
			// we hadn't found a closing bracket, but the current bracket is a closing bracket
			// because the next bracket opens the description
			if current < itemsLength && items[current+1].IsBracket() && items[current+1].Value() == "[" {
				foundLinkCloseBracket = true
			}
			// this is a self-describing link
			if foundLinkCloseBracket && items[current+1].IsBracket() && items[current+1].Value() == "]" {
				return true
			}
			// if it's a bracket and it's not a closing bracket and we haven't found a closed bracket yet, this isn't a link
		} else if items[current].IsBracket() && items[current].Value() != "]" {
			return false
		}
		current++
	}
	return false
}
