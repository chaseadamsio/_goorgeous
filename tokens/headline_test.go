package tokens

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
)

func TestIsHeadline(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"* is a headline", true},
		{"** is a headline", true},
		{"*** is a headline", true},
		{"**** is a headline", true},
		{"***** is a headline", true},
		{"****** is a headline", true},
		{"******* is NOT a headline", false},
		{"*is NOT a headline", false},
		{"**is NOT a headline", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsHeadline(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestHeadlineDepth(t *testing.T) {
	testCases := []struct {
		value    string
		expected int
	}{
		{"* is a headline", 1},
		{"** is a headline", 2},
		{"*** is a headline", 3},
		{"**** is a headline", 4},
		{"***** is a headline", 5},
		{"****** is a headline", 6},
	}

	for _, tc := range testCases {
		var items []lex.Item

		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}

		depth := HeadlineDepth(items)
		if depth != tc.expected {
			t.Errorf("expected depth of \"%s\" to be %d. Instead, got %d", tc.value, tc.expected, depth)
		}
	}
}
