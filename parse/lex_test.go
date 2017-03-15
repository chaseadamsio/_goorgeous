package parse

import "testing"

var (
	tEOF     = mkItem(itemEOF, "")
	tNewLine = mkItem(itemNewLine, "\n")
)

func TestLex(t *testing.T) {
	for _, tc := range testCases {
		l := lex(tc.input)
		items := []item{}
		for {
			item := l.nextItem()
			items = append(items, item)

			if item.typ == itemEOF {
				break
			}
		}

		if !equal(items, tc.expectedLex, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", tc.name, items, tc.expectedLex)
		}

	}
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}
