package parse

import "testing"

var testCases = []struct {
	name         string
	input        string
	expectedLex  []item
	expectedHTML string
}{
	{
		"basic",
		"hello world!",
		[]item{mkItem(itemText, "hello world!"), tEOF},
		"<p>hello world!</p>",
	},
	{
		"basic w/ multi-line",
		"hello world!\nIt's me, your friend!",
		[]item{mkItem(itemText, "hello world!"), tNewLine, mkItem(itemText, "It's me, your friend!"), tEOF},
		"<p>hello world!\n\nIt's me, your friend!</p>",
	},
	{
		"basic w/ multi-line and empty new line",
		"hello world!\n\nIt's me, your friend!",
		[]item{mkItem(itemText, "hello world!"), tNewLine, tNewLine, mkItem(itemText, "It's me, your friend!"), tEOF},
		"<p>hello world!\n\nIt's me, your friend!</p>",
	},
}

func TestRenderAsHTML(t *testing.T) {
	nl, err := Parse("hello world!")
	if err != nil {
		t.Fatalf("%s failed: %v", "test", err)
	}

	html := RenderAsHTML(nl)
	t.Logf("%v", html)
}

func TestParse(t *testing.T) {
	nl, err := Parse("hello world!")
	if err != nil {
		t.Fatalf("%s failed: %v", "test", err)
	}
	for _, l := range nl {
		t.Logf("%v", l)
	}
}
