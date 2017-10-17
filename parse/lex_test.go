package parse

import (
	"bytes"
	"testing"
)

var (
	tEOF          = mkItem(elEOF, "")
	tSpace        = mkItem(elSpace, " ")
	tNewline      = mkItem(elNewline, "\n")
	tAsterisk     = mkItem(elAsterisk, "*")
	tHash         = mkItem(elHash, "#")
	tSlash        = mkItem(elSlash, "/")
	tEqual        = mkItem(elEqual, "=")
	tTilde        = mkItem(elTilde, "~")
	tUnderscore   = mkItem(elUnderscore, "_")
	tColon        = mkItem(elColon, ":")
	tBracketLeft  = mkItem(elBracketLeft, "[")
	tBracketRight = mkItem(elBracketRight, "]")
	tPipe         = mkItem(elPipe, "|")
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
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tEOF,
		}},

	"simple string with newline": {
		"this is some text\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 1": {
		"* this is some text\n",
		[]item{
			tAsterisk,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 2": {
		"** this is some text\n",
		[]item{
			tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 3": {
		"*** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 4": {
		"**** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 5": {
		"***** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 6": {
		"****** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"not header": {
		"this ***** is some text\n",
		[]item{
			mkItem(elWord, "this"),
			tSpace,
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"not header alt": {
		"this***** is some text\n",
		[]item{
			mkItem(elWord, "this"),
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"bold": {"this is *some text*\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace,
			tAsterisk,
			mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tAsterisk,
			tNewline,
			tEOF,
		}},

	"not bold": {"this is *some text\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace,
			tAsterisk,
			mkItem(elWord, "some"), tSpace, mkItem(elWord, "text"),
			tNewline,
			tEOF,
		}},

	"comment": {"# this is a comment\n",
		[]item{
			tHash,
			tSpace, mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "comment"),
			tNewline,
			tEOF,
		}},

	"not comment": {"#this is not a comment\n",
		[]item{
			tHash,
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "not"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "comment"),
			tNewline,
			tEOF,
		}},

	"underline": {"_this is a sentence_ with underline.\n",
		[]item{
			tUnderscore,
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "sentence"),
			tUnderscore,
			tSpace, mkItem(elWord, "with"), tSpace, mkItem(elWord, "underline."),
			tNewline,
			tEOF,
		}},

	"italic": {"/this is a sentence/ with italic.\n",
		[]item{
			tSlash,
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "sentence"),
			tSlash,
			tSpace, mkItem(elWord, "with"), tSpace, mkItem(elWord, "italic."),
			tNewline,
			tEOF,
		}},

	"inline verbatim": {"=this is a sentence= with verbatim.\n",
		[]item{
			tEqual,
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "sentence"),
			tEqual,
			tSpace, mkItem(elWord, "with"), tSpace, mkItem(elWord, "verbatim."),
			tNewline,
			tEOF,
		}},

	"inline code": {"~this is a sentence~ with code.\n",
		[]item{
			tTilde,
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "is"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "sentence"),
			tTilde,
			tSpace, mkItem(elWord, "with"), tSpace, mkItem(elWord, "code."),
			tNewline,
			tEOF,
		}},

	"anchor - link as URL": {"this has [[https://github.com/chaseadamsio/goorgeous]] as a link.\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(elWord, "https"), tColon, tSlash, tSlash, mkItem(elWord, "github.com"), tSlash, mkItem(elWord, "chaseadamsio"), tSlash, mkItem(elWord, "goorgeous"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(elWord, "as"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "link."),
			tNewline,
			tEOF,
		}},

	"anchor - text": {"this has [[https://github.com/chaseadamsio/goorgeous][goorgeous by chaseadamsio]] as a link.\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(elWord, "https"), tColon, tSlash, tSlash, mkItem(elWord, "github.com"), tSlash, mkItem(elWord, "chaseadamsio"), tSlash, mkItem(elWord, "goorgeous"),
			tBracketRight, tBracketLeft,
			mkItem(elWord, "goorgeous"), tSpace, mkItem(elWord, "by"), tSpace, mkItem(elWord, "chaseadamsio"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(elWord, "as"), tSpace, mkItem(elWord, "a"), tSpace, mkItem(elWord, "link."),
			tNewline,
			tEOF,
		}},

	"image - basic": {"this has [[file:https://github.com/chaseadamsio/goorgeous/img.png]] as an image.\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(elWord, "file"), tColon,
			mkItem(elWord, "https"), tColon, tSlash, tSlash, mkItem(elWord, "github.com"), tSlash, mkItem(elWord, "chaseadamsio"), tSlash, mkItem(elWord, "goorgeous"), tSlash, mkItem(elWord, "img.png"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(elWord, "as"), tSpace, mkItem(elWord, "an"), tSpace, mkItem(elWord, "image."),
			tNewline,
			tEOF,
		}},

	"image - alt": {"this has [[file:../gopher.gif][a uni-gopher]] as an image.\n",
		[]item{
			mkItem(elWord, "this"), tSpace, mkItem(elWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(elWord, "file"), tColon,
			mkItem(elWord, ".."), tSlash, mkItem(elWord, "gopher.gif"),
			tBracketRight, tBracketLeft,
			mkItem(elWord, "a"), tSpace, mkItem(elWord, "uni-gopher"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(elWord, "as"), tSpace, mkItem(elWord, "an"), tSpace, mkItem(elWord, "image."),
			tNewline,
			tEOF,
		}},
}

func TestLexer(t *testing.T) {
	for caseName, tc := range testCases {
		l := NewLexer(tc.input)
		items := collect(l)
		if !equal(tc.items, items, false) {
			t.Errorf("'%s' case failed. items are not equal.\n got  %v+\n want %v\n", caseName, items, tc.items)
		}
	}
}

// mkItem is a helper to make it easier to generate items for
// test cases
func mkItem(typ elType, val string) item {
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
		if item.typ == elEOF || item.typ == elError {
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
