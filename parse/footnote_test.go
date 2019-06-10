package parse

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
)

func TestIsFootnoteDefinition(t *testing.T) {
	testCases := []struct {
		input    string
		start    int
		expected bool
	}{
		{"[fn:1] The link is: https://orgmode.org", 0, true},
		{"[fn::This is the inline definition of this footnote]", 0, true},
		{"foo\n[fn::This is the inline definition of this footnote]", 2, true},
		{"[fn:name:a definition]", 0, true},
		{"foo[fn::This is the inline definition of this footnote]", 1, true},
		{"foo[fn::This is the inline definition of this footnote]", 0, false},
	}
	for _, tc := range testCases {
		var items []lex.Item
		itemsChan := lex.NewLexer(tc.input)
		for item := range itemsChan {
			items = append(items, item)
		}
		if isFootnoteDefinition(items[tc.start:]) != tc.expected {
			t.Errorf("Expected %s to be %t", tc.input, tc.expected)
		}
	}
}
