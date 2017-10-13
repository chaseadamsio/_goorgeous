package parse

import "strconv"

// itemType identifies the type of lex items
type itemType int

const (
	itemError itemType = iota
	itemNewline
	itemEOF
	itemText
)

var itemTypes = [...]string{
	itemError:   "itemError",
	itemNewline: "itemNewline",
	itemEOF:     "itemEOF",
	itemText:    "itemText",
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
