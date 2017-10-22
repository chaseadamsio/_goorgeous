package parse

import (
	"bytes"
	"testing"
)

var (
	tEOF          = mkItem(tokenEOF, "")
	tSpace        = mkItem(tokenSpace, " ")
	tNewline      = mkItem(tokenNewline, "\n")
	tAsterisk     = mkItem(tokenAsterisk, "*")
	tHash         = mkItem(tokenHash, "#")
	tPlus         = mkItem(tokenPlus, "+")
	tSlash        = mkItem(tokenSlash, "/")
	tEqual        = mkItem(tokenEqual, "=")
	tTilde        = mkItem(tokenTilde, "~")
	tDash         = mkItem(tokenDash, "-")
	tUnderscore   = mkItem(tokenUnderscore, "_")
	tColon        = mkItem(tokenColon, ":")
	tBracketLeft  = mkItem(tokenBracketLeft, "[")
	tBracketRight = mkItem(tokenBracketRight, "]")
	tPipe         = mkItem(tokenPipe, "|")
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
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tEOF,
		}},

	"simple string with newline": {
		"this is some text\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 1": {
		"* this is some text\n",
		[]item{
			tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"previous text - header level 1": {
		"this is some text.\n* this is some text\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text."),
			tNewline,
			tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 2": {
		"** this is some text\n",
		[]item{
			tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 3": {
		"*** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 4": {
		"**** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 5": {
		"***** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"header level 6": {
		"****** this is some text\n",
		[]item{
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"not header": {
		"this ***** is some text\n",
		[]item{
			mkItem(tokenWord, "this"),
			tSpace,
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"not header alt": {
		"this***** is some text\n",
		[]item{
			mkItem(tokenWord, "this"),
			tAsterisk, tAsterisk, tAsterisk, tAsterisk, tAsterisk,
			tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"bold": {"this is *some text*\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace,
			tAsterisk,
			mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tAsterisk,
			tNewline,
			tEOF,
		}},

	"not bold": {"this is *some text\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace,
			tAsterisk,
			mkItem(tokenWord, "some"), tSpace, mkItem(tokenWord, "text"),
			tNewline,
			tEOF,
		}},

	"comment": {"# this is a comment\n",
		[]item{
			tHash,
			tSpace, mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "comment"),
			tNewline,
			tEOF,
		}},

	"not comment": {"#this is not a comment\n",
		[]item{
			tHash,
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "not"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "comment"),
			tNewline,
			tEOF,
		}},

	"underline": {"_this is a sentence_ with underline.\n",
		[]item{
			tUnderscore,
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "sentence"),
			tUnderscore,
			tSpace, mkItem(tokenWord, "with"), tSpace, mkItem(tokenWord, "underline."),
			tNewline,
			tEOF,
		}},

	"italic": {"/this is a sentence/ with italic.\n",
		[]item{
			tSlash,
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "sentence"),
			tSlash,
			tSpace, mkItem(tokenWord, "with"), tSpace, mkItem(tokenWord, "italic."),
			tNewline,
			tEOF,
		}},

	"strikethrough": {"+this is a sentence+ with strikethrough.\n",
		[]item{
			tPlus,
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "sentence"),
			tPlus,
			tSpace, mkItem(tokenWord, "with"), tSpace, mkItem(tokenWord, "strikethrough."),
			tNewline,
			tEOF,
		}},

	"inline verbatim": {"=this is a sentence= with verbatim.\n",
		[]item{
			tEqual,
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "sentence"),
			tEqual,
			tSpace, mkItem(tokenWord, "with"), tSpace, mkItem(tokenWord, "verbatim."),
			tNewline,
			tEOF,
		}},

	"inline code": {"~this is a sentence~ with code.\n",
		[]item{
			tTilde,
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "is"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "sentence"),
			tTilde,
			tSpace, mkItem(tokenWord, "with"), tSpace, mkItem(tokenWord, "code."),
			tNewline,
			tEOF,
		}},

	"anchor - link as URL": {"this has [[https://github.com/chaseadamsio/goorgeous]] as a link.\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(tokenWord, "https"), tColon, tSlash, tSlash, mkItem(tokenWord, "github.com"), tSlash, mkItem(tokenWord, "chaseadamsio"), tSlash, mkItem(tokenWord, "goorgeous"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(tokenWord, "as"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "link."),
			tNewline,
			tEOF,
		}},

	"anchor - text": {"this has [[https://github.com/chaseadamsio/goorgeous][goorgeous by chaseadamsio]] as a link.\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(tokenWord, "https"), tColon, tSlash, tSlash, mkItem(tokenWord, "github.com"), tSlash, mkItem(tokenWord, "chaseadamsio"), tSlash, mkItem(tokenWord, "goorgeous"),
			tBracketRight, tBracketLeft,
			mkItem(tokenWord, "goorgeous"), tSpace, mkItem(tokenWord, "by"), tSpace, mkItem(tokenWord, "chaseadamsio"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(tokenWord, "as"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "link."),
			tNewline,
			tEOF,
		}},

	"image - basic": {"this has [[file:https://github.com/chaseadamsio/goorgeous/img.png]] as an image.\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(tokenWord, "file"), tColon,
			mkItem(tokenWord, "https"), tColon, tSlash, tSlash, mkItem(tokenWord, "github.com"), tSlash, mkItem(tokenWord, "chaseadamsio"), tSlash, mkItem(tokenWord, "goorgeous"), tSlash, mkItem(tokenWord, "img.png"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(tokenWord, "as"), tSpace, mkItem(tokenWord, "an"), tSpace, mkItem(tokenWord, "image."),
			tNewline,
			tEOF,
		}},

	"image - alt": {"this has [[file:../gopher.gif][a uni-gopher]] as an image.\n",
		[]item{
			mkItem(tokenWord, "this"), tSpace, mkItem(tokenWord, "has"), tSpace,
			tBracketLeft, tBracketLeft,
			mkItem(tokenWord, "file"), tColon,
			mkItem(tokenWord, ".."), tSlash, mkItem(tokenWord, "gopher.gif"),
			tBracketRight, tBracketLeft,
			mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "uni"), tDash, mkItem(tokenWord, "gopher"),
			tBracketRight, tBracketRight,
			tSpace, mkItem(tokenWord, "as"), tSpace, mkItem(tokenWord, "an"), tSpace, mkItem(tokenWord, "image."),
			tNewline,
			tEOF,
		}},

	"definition": {"- definition lists :: these are useful sometimes\n- item 2 :: M-RET again gives another item, and long lines wrap in a tidy way underneath the definition\n",
		[]item{
			tDash,
			tSpace, mkItem(tokenWord, "definition"), tSpace, mkItem(tokenWord, "lists"), tSpace,
			tColon, tColon,
			tSpace, mkItem(tokenWord, "these"), tSpace, mkItem(tokenWord, "are"), tSpace, mkItem(tokenWord, "useful"), tSpace, mkItem(tokenWord, "sometimes"),
			tNewline,
			tDash,
			tSpace, mkItem(tokenWord, "item"), tSpace, mkItem(tokenWord, "2"), tSpace,
			tColon, tColon,
			tSpace, mkItem(tokenWord, "M"), tDash, mkItem(tokenWord, "RET"), tSpace, mkItem(tokenWord, "again"), tSpace, mkItem(tokenWord, "gives"), tSpace, mkItem(tokenWord, "another"),
			tSpace, mkItem(tokenWord, "item,"), tSpace, mkItem(tokenWord, "and"), tSpace, mkItem(tokenWord, "long"), tSpace, mkItem(tokenWord, "lines"), tSpace, mkItem(tokenWord, "wrap"),
			tSpace, mkItem(tokenWord, "in"), tSpace, mkItem(tokenWord, "a"), tSpace, mkItem(tokenWord, "tidy"), tSpace, mkItem(tokenWord, "way"),
			tSpace, mkItem(tokenWord, "underneath"), tSpace, mkItem(tokenWord, "the"), tSpace, mkItem(tokenWord, "definition"),
			tNewline,
			tEOF,
		}},

	"ul - plus": {"+ this\n+ is\n+ an\n+ unordered\n+ list\n",
		[]item{
			tPlus,
			tSpace, mkItem(tokenWord, "this"),
			tNewline,
			tPlus,
			tSpace, mkItem(tokenWord, "is"),
			tNewline,
			tPlus,
			tSpace, mkItem(tokenWord, "an"),
			tNewline,
			tPlus,
			tSpace, mkItem(tokenWord, "unordered"),
			tNewline,
			tPlus,
			tSpace, mkItem(tokenWord, "list"),
			tNewline,
			tEOF,
		}},

	"ul - dash": {"- this\n- is\n- an\n- unordered\n- list\n",
		[]item{
			tDash,
			tSpace, mkItem(tokenWord, "this"),
			tNewline,
			tDash,
			tSpace, mkItem(tokenWord, "is"),
			tNewline,
			tDash,
			tSpace, mkItem(tokenWord, "an"),
			tNewline,
			tDash,
			tSpace, mkItem(tokenWord, "unordered"),
			tNewline,
			tDash,
			tSpace, mkItem(tokenWord, "list"),
			tNewline,
			tEOF,
		}},

	"SRC block": {"#+BEGIN_SRC sh\necho \"foo\"\n#+END_SRC\n",
		[]item{
			tHash, tPlus,
			mkItem(tokenWord, "BEGIN"), tUnderscore, mkItem(tokenWord, "SRC"), tSpace, mkItem(tokenWord, "sh"),
			tNewline,
			mkItem(tokenWord, "echo"), tSpace, mkItem(tokenWord, "\"foo\""),
			tNewline,
			tHash, tPlus,
			mkItem(tokenWord, "END"), tUnderscore, mkItem(tokenWord, "SRC"),
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

	for _, tc := range testCases {
		l := NewLexer(tc.input)
		eval(l)
	}
}

// mkItem is a helper to make it easier to generate items for
// test cases
func mkItem(typ tokenType, val string) item {
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
		if item.typ == tokenEOF || item.typ == tokenError {
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
