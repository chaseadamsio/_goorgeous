package parse

import (
	"bytes"
	"testing"
)

var (
	tEOF      = mkItem(itemEOF, "")
	tNewline  = mkItem(itemNewline, "\n")
	tAsterisk = mkItem(itemAsterisk, "*")
	tComment  = mkItem(itemComment, "# ")
)

// testCase is a test input string and
// the expected output items
type testCase struct {
	input string
	items []item
}

var testCases = map[string]testCase{
	"empty string": {
		"", // should handle empty strings gracefully
		[]item{
			tEOF,
		}},

	"simple string no newline": {
		"this is some text",
		[]item{
			mkItem(itemText, "this is some text"),
			tEOF,
		}},

	"simple string with newline": {
		"this is some text\n",
		[]item{
			mkItem(itemText, "this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - header level 1": {
		"* this is some text\n",
		[]item{
			tAsterisk,
			mkItem(itemText, " this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - header level 2": {
		"** this is some text\n",
		[]item{
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - header level 3": {
		"*** this is some text\n",
		[]item{
			tAsterisk,
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - header level 4": {
		"**** this is some text\n",
		[]item{
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - header level 5": {
		"***** this is some text\n",
		[]item{
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - header level 6": {
		"****** this is some text\n",
		[]item{
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " this is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - not header": {
		"this ***** is some text\n",
		[]item{
			mkItem(itemText, "this "),
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk - not header - alt": {
		"this***** is some text\n",
		[]item{
			mkItem(itemText, "this"),
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			tAsterisk,
			mkItem(itemText, " is some text"),
			tNewline,
			tEOF,
		}},

	"asterisk bold": {"this is *some text*\n",
		[]item{
			mkItem(itemText, "this is "),
			tAsterisk,
			mkItem(itemText, "some text"),
			tAsterisk,
			tNewline,
			tEOF,
		}},

	"asterisk - not bold": {"this is *some text\n",
		[]item{
			mkItem(itemText, "this is "),
			tAsterisk,
			mkItem(itemText, "some text"),
			tNewline,
			tEOF,
		}},

	"comment": {"# this is a comment\n",
		[]item{
			tComment,
			mkItem(itemText, "this is a comment"),
			tNewline,
			tEOF,
		}},
}

func TestLexer(t *testing.T) {
	for caseName, tc := range testCases {
		l := NewLexer(tc.input)
		items := collect(l)
		if !equal(tc.items, items, false) {
			t.Errorf("'%s' case failed. items are not equal.\n got %v+\n expected %v\n", caseName, items, tc.items)
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
func collect(l *Lexer) (items []item) {
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
