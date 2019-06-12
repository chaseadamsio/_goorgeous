package parse

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
	"github.com/chaseadamsio/goorgeous/testdata"
)

func TestIsUnorderedList(t *testing.T) {
	var testCases = []struct {
		source string
		start  int
		found  bool
	}{
		{testdata.UnorderedListBasic, 0, true},
		{testdata.UnorderedListNotAList, -1, false},
		{testdata.UnorderedListWithStartingNewline, -1, false},
		{testdata.Headline1, -1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.source, func(t *testing.T) {
			value := testdata.GetOrgStr(tc.source)

			var items []lex.Item
			lexedItems := lex.NewLexer(value)
			for item := range lexedItems {
				items = append(items, item)
			}

			start, found := isUnorderedList(items)
			if tc.start != start {
				t.Errorf("expected %s start to be %d. got %d", tc.source, tc.start, start)
			}
			if tc.found != found {
				t.Errorf("expected %s found to be %t", tc.source, tc.found)
			}

		})
	}
}

func TestIsOrderedList(t *testing.T) {

	testCases := []struct {
		source string
		start  int
		found  bool
	}{
		{testdata.OrderedListBasic, 0, true},
		{testdata.OrderedListNotAList, -1, false},
		{testdata.OrderedListWithStartingNewline, -1, false},
		{testdata.Headline1, -1, false},
	}

	for _, tc := range testCases {
		value := testdata.GetOrgStr(tc.source)

		var items []lex.Item
		lexedItems := lex.NewLexer(value)
		for item := range lexedItems {
			items = append(items, item)
		}

		start, found := isOrderedList(items)
		if tc.start != start {
			t.Errorf("expected %s start to be %d. got %d", tc.source, tc.start, start)
		}
		if tc.found != found {
			t.Errorf("expected %s found to be %t", tc.source, tc.found)
		}
	}
}

func TestFindUnorderedListBoundaries(t *testing.T) {
	testCases := []struct {
		source     string
		start, end int
	}{
		{testdata.UnorderedListBasic, 0, 16},
		{testdata.UnorderedListFollowParagraph, 0, 16},
		{testdata.UnorderedListFollowDashNotList, 0, 12},
		{testdata.UnorderedListFollowAsteriskHeading, 0, 12},
		{testdata.UnorderedListWithFollowOrderedList, 0, 16},
		{testdata.UnorderedListWithNestedOrderedList, 0, 56},
		{testdata.UnorderedListWithNestedContent, 0, 30},
		{testdata.UnorderedListWithDeepNestedChildren, 0, 78},
	}

	for _, tc := range testCases {
		value := testdata.GetOrgStr(tc.source)

		var items []lex.Item
		lexedItems := lex.NewLexer(value)
		for item := range lexedItems {
			items = append(items, item)
		}
		start, end := findUnorderedListBoundaries(items)
		if tc.start != start {
			t.Errorf("expected \"%s\" start to be %d. got %d", tc.source, tc.start, start)
		}
		if tc.end != end {
			t.Errorf("expected \"%s\" end to be %d. got %d", tc.source, tc.end, end)
		}
	}
}

func TestFindOrderedList(t *testing.T) {
	testCases := []struct {
		source     string
		start, end int
	}{
		{testdata.OrderedListBasic, 0, 16},
		{testdata.OrderedListFollowParagraph, 0, 16},
		{testdata.OrderedListFollowAsteriskHeading, 0, 16},
		{testdata.OrderedListFollowNumberNotList, 0, 16},
		{testdata.OrderedListWithFollowUnOrderedList, 0, 16},
		{testdata.OrderedListWithNestedOrderedList, 0, 56},
		{testdata.OrderedListWithNestedContent, 0, 30},
	}

	for _, tc := range testCases {
		value := testdata.GetOrgStr(tc.source)

		var items []lex.Item
		lexedItems := lex.NewLexer(value)
		for item := range lexedItems {
			items = append(items, item)
		}

		start, end := findOrderedListBoundaries(items)
		if tc.start != start {
			t.Errorf("expected \"%s\" start to be %d. got %d", tc.source, tc.start, start)
		}
		if tc.end != end {
			t.Errorf("expected \"%s\" end to be %d. got %d", tc.source, tc.end, end)
		}
	}
}
