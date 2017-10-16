package parse

import "strconv"

// itemType identifies the type of lex items
type itemType int

const (
	itemError itemType = iota
	itemNewline
	itemEOF
	itemText
	itemAsterisk // "*" indicates either a headline or a bold token
	itemComment  // "#  " indicates a comment token
)

var itemTypes = [...]string{
	itemError:    "itemError",
	itemNewline:  "itemNewline",
	itemEOF:      "itemEOF",
	itemText:     "itemText",
	itemAsterisk: "*",
}

func (typ itemType) String() string {
	s := ""
	if 0 <= typ && typ < itemType(len(itemTypes)) {
		s = itemTypes[typ]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(typ)) + ")"
	}
	return s
}
