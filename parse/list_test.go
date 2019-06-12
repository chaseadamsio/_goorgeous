package parse

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
	"github.com/chaseadamsio/goorgeous/testdata"
)

func TestIsUnorderedList(t *testing.T) {

	var testCases = []struct {
		source    string
		startItem int
		expected  bool
	}{
		{testdata.UnorderedListBasic, 0, true},
		{testdata.UnorderedListNotAList, 0, false},
		{testdata.UnorderedListWithStartingNewline, 1, true},
		{testdata.UnorderedListWithStartingNewline, 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.source, func(t *testing.T) {
			value := testdata.GetOrgStr(tc.source)

			var items []lex.Item
			lexedItems := lex.NewLexer(value)
			for item := range lexedItems {
				items = append(items, item)
			}
			if isUnorderedList(items[tc.startItem:]) != tc.expected {
				t.Errorf("expected \"%s\" to be %t", value, tc.expected)
			}
		})
	}
}

func TestIsOrderedList(t *testing.T) {

	testCases := []struct {
		source    string
		startItem int
		expected  bool
	}{
		{testdata.OrderedListBasic, 0, true},
		{testdata.OrderedListNotAList, 0, false},
		{testdata.OrderedListWithStartingNewline, 1, true},
		{testdata.OrderedListWithStartingNewline, 0, false},
	}

	for _, tc := range testCases {
		value := testdata.GetOrgStr(tc.source)

		var items []lex.Item
		lexedItems := lex.NewLexer(value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if isOrderedList(items[tc.startItem:]) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", value, tc.expected)
		}
	}
}

func TestFindIsUnorderedList(t *testing.T) {
	testCases := []struct {
		source   string
		expected int
	}{
		{testdata.UnorderedListBasic, 15},
		{testdata.UnorderedListFollowParagraph, 15},
		{testdata.UnorderedListFollowDashNotList, 11},
		{testdata.UnorderedListFollowAsteriskHeading, 11},
		{testdata.UnorderedListWithFollowOrderedList, 15},
		{testdata.UnorderedListWithNestedOrderedList, 55},
		{testdata.UnorderedListWithNestedContent, 29},
		{testdata.UnorderedListWithDeepNestedChildren, 77},
	}

	for _, tc := range testCases {
		value := testdata.GetOrgStr(tc.source)

		var items []lex.Item
		lexedItems := lex.NewLexer(value)
		for item := range lexedItems {
			items = append(items, item)
		}
		found := findUnorderedList(items)
		if found != tc.expected {
			t.Errorf("expected \"%s\" to be %d. got %d", tc.source, tc.expected, found)
		}
	}
}

func TestFindOrderedList(t *testing.T) {
	testCases := []struct {
		source   string
		expected int
	}{
		{testdata.OrderedListBasic, 15},
		{testdata.OrderedListFollowParagraph, 15},
		{testdata.OrderedListFollowAsteriskHeading, 15},
		{testdata.OrderedListFollowNumberNotList, 15},
		{testdata.OrderedListWithFollowUnOrderedList, 15},
		{testdata.OrderedListWithNestedOrderedList, 55},
		{testdata.OrderedListWithNestedContent, 29},
	}

	for _, tc := range testCases {
		value := testdata.GetOrgStr(tc.source)

		var items []lex.Item
		lexedItems := lex.NewLexer(value)
		for item := range lexedItems {
			items = append(items, item)
		}
		found := findOrderedList(items)
		if found != tc.expected {
			t.Errorf("expected \"%s\" to be %d. got %d", tc.source, tc.expected, found)
		}
	}
}
