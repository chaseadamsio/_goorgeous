package goorgeous

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/russross/blackfriday"
)

var update = flag.Bool("update", false, "update golden files")

func TestIsTable(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"|some table", true},
		{"| some table", true},
		{" | not a table", false},
		{"not a table", false},
		{"*not a table", false},
		{"-not a table", false},
		{"+not a table", false},
	}

	for _, tc := range testCases {
		isTable := isTable([]byte(tc.in))
		if isTable != tc.expected {
			t.Errorf("isTable(%s) = %T\nwants: %T", tc.in, isTable, tc.expected)
		}
	}
}

func TestSkipChar(t *testing.T) {
	testCases := []struct {
		in       string
		start    int
		char     byte
		expected int
	}{
		{"   check for spaces", 0, ' ', 3},
		{" -  check for spaces", 0, ' ', 1},
		{"check     for spaces", 1, ' ', 1},
		{"check     for spaces", 5, ' ', 10},
		{"check-----for spaces", 5, '-', 10},
	}

	for _, tc := range testCases {
		skipped := skipChar([]byte(tc.in), tc.start, tc.char)
		if skipped != tc.expected {
			t.Errorf("skipChar(%s, %d, %q) = %d\nwants: %d", tc.in, tc.start, tc.char, skipped, tc.expected)
		}
	}
}

func TestIsSpace(t *testing.T) {
	testCases := []struct {
		char     byte
		expected bool
	}{

		{' ', true},
		{'+', false},
	}

	for _, tc := range testCases {
		isSpace := isSpace(tc.char)
		if isSpace != tc.expected {
			t.Errorf("isSpace(%q) = %T/nwants: %T", tc.char, isSpace, tc.expected)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"\n", true},
		{"\t", true},
		{"\t\n\t", true},
		{"\tfoo\n\t", false},
	}

	for _, tc := range testCases {
		isEmpty := isEmpty([]byte(tc.in))
		if isEmpty != tc.expected {
			t.Errorf("isEmpty(%s) = %T\nwants: %T", tc.in, isEmpty, tc.expected)
		}
	}
}

func TestCharMatch(t *testing.T) {
	testCases := []struct {
		a        byte
		b        byte
		expected bool
	}{
		{' ', ' ', true},
		{'+', '+', true},
		{'#', '#', true},
		{'-', '-', true},
		{'+', '-', false},
		{' ', '-', false},
	}

	for _, tc := range testCases {
		matches := charMatches(tc.a, tc.b)
		if matches != tc.expected {
			t.Errorf("charMatches(%q, %q) = %T\nwants: %T", tc.a, tc.b, matches, tc.expected)
		}
	}
}

func TestIsHeader(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"#+TITLE", true},
		{"#+ TITLE", false},
		{"# TITLE", false},
		{"TITLE", false},
	}

	for _, tc := range testCases {
		isHeader := isHeader([]byte(tc.in))
		if isHeader != tc.expected {
			t.Errorf("isHeader(%s) = %T\nwants: %T", tc.in, isHeader, tc.expected)
		}
	}
}

func TestIsComment(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"# this is a comment", true},
		{"#-this is not a comment", false},
		{"#+TITLE", false},
		{"This is not a comment", false},
	}

	for _, tc := range testCases {
		isComment := isComment([]byte(tc.in))
		if isComment != tc.expected {
			t.Errorf("isComment(%s) = %T\nwants: %T", tc.in, isComment, tc.expected)
		}
	}
}

func TestGenerateComment(t *testing.T) {
	p := NewParser(blackfriday.HtmlRenderer(blackfriday.HTML_USE_XHTML, "", ""))
	var out bytes.Buffer
	text := "This is a comment and we expect it to look a certain way."
	orgComment := []byte("# " + text)
	expected := "<!-- " + text + " -->\n"
	p.generateComment(&out, orgComment)
	if out.String() != expected {
		t.Errorf("generateComment(%s) = %s\nwants: %s", text, out.String(), expected)
	}
}

func TestOrgCommonFromFile(t *testing.T) {
	source := "./testdata/test.org"
	golden := "./testdata/test.html.golden"
	contents, err := ioutil.ReadFile(source)
	if err != nil {
		t.Fatalf("failed to read %s file: %s", source, err)
	}

	out := OrgCommon(contents)

	if *update {
		if err := ioutil.WriteFile(golden, out, 0644); err != nil {
			t.Errorf("failed to write %s file: %s", golden, err)
		}
		return
	}

	gld, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatalf("failed to read %s file: %s", golden, err)
	}

	if !bytes.Equal(out, gld) {
		t.Errorf("OrgCommon() from %s = %s\nwants: %s", source, out, gld)
	}
}
