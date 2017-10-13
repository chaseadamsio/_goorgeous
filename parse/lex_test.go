package parse

import (
	"bytes"
	"testing"
)

var (
	tEOF     = mkItem(itemEOF, "")
	tNewline = mkItem(itemNewline, "\n")
)

// testCase is a test input string and
// the expected output items
type testCase struct {
	input string
	items []item
}

func TestLexer(t *testing.T) {
	testCases := []testCase{
		{"this is some text", []item{
			mkItem(itemText, "this is some text"),
			tEOF,
		}},
		{"this is some text\n", []item{
			mkItem(itemText, "this is some text"),
			tNewline,
			tEOF,
		}},
	}

	for _, tc := range testCases {
		items := collect(&tc)
		if !equal(tc.items, items, false) {
			t.Errorf("items are not equal.\n got %v+\n expected %v\n", items, tc.items)
		}
	}
}

// mkItem is a helper to make it easier to generate items for
// test cases
func mkItem(typ itemType, val string) item {
	return item{
		typ: typ,
		val: []byte(val),
	}
}

// collect runs the lexer and collects all of the items that are
// emitted by nextItem, and returns a slice of item
func collect(tc *testCase) (items []item) {
	l := NewLexer(tc.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return items
}

// equal checks that two slices of item are equal in both type
// and in value
func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if !bytes.Equal(i1[k].val, i2[k].val) {
			return false
		}
		if checkPos && i1[k].end != i2[k].end {
			return false
		}
	}
	return true
}
