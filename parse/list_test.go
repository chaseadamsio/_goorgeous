package parse

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
)

func TestIsUnorderedList(t *testing.T) {
	testCases := []struct {
		value    string
		start    int
		expected bool
	}{
		{"- apples\n- bananas\n- oranges", 0, true},
		{"\n- apples\n- bananas\n- oranges", 1, true},
		{"\n- apples\n- bananas\n- oranges", 0, false},
		{"beforedash- apples\n- bananas\n- oranges", 0, false},
		{"\n-apples\n-bananas\n-oranges", 0, false},
		{"\n-apples\n-bananas\n-oranges", 1, false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if isUnorderedList(items[tc.start:]) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestIsOrderedList(t *testing.T) {
	testCases := []struct {
		value    string
		start    int
		expected bool
	}{
		{"1. apples\n2. bananas\n3. oranges", 0, true},
		{"\n1. apples\n2. bananas\n3. oranges", 1, true},
		{"\n1. apples\n2. bananas\n3. oranges", 0, false},
		{"beforedash1. apples\n2. bananas\n3. oranges", 0, false},
		{"\n1.apples\n2.bananas\n3.oranges", 0, false},
		{"\n1.apples\n2.bananas\n3.oranges", 1, false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if isOrderedList(items[tc.start:]) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestFindIsUnorderedList(t *testing.T) {
	testCases := []struct {
		value    string
		expected int
	}{
		{"- apples\n- bananas\n- oranges", 11},
		// {"- apples\n- bananas\n- oranges\n\tsome text", 17},
		{"- apples\n- bananas\n- oranges\nsome text", 11},
		{"- apples\n- bananas\n- oranges\n* test", 11},
		{"- apples\n- bananas\n- oranges\n1. test", 11},
		// {"\n- apples\n- bananas\n- oranges", 1},
		// {"\n- apples\n- bananas\n- oranges", 0},
		// {"beforedash- apples\n- bananas\n- oranges", 0},
		// {"\n-apples\n-bananas\n-oranges", 0},
		// {"\n-apples\n-bananas\n-oranges", 1},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		found := findUnorderedList(items)
		if found != tc.expected {
			t.Errorf("expected \"%s\" to be %d. got %d", tc.value, tc.expected, found)
		}
	}
}

func TestFindIsOrderedList(t *testing.T) {
	testCases := []struct {
		value    string
		expected int
	}{
		{"1. apples\n2. bananas\n3. oranges", 11},
		// {"- apples\n- bananas\n- oranges\n\tsome text", 17},
		{"1. apples\n2. bananas\n3. oranges\nsome text", 11},
		{"1. apples\n2. bananas\n3. oranges\n* test", 11},
		{"1. apples\n2. bananas\n3. oranges\n- test", 11},
		// {"\n- apples\n- bananas\n- oranges", 1},
		// {"\n- apples\n- bananas\n- oranges", 0},
		// {"beforedash- apples\n- bananas\n- oranges", 0},
		// {"\n-apples\n-bananas\n-oranges", 0},
		// {"\n-apples\n-bananas\n-oranges", 1},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		found := findOrderedList(items)
		if found != tc.expected {
			t.Errorf("expected \"%s\" to be %d. got %d", tc.value, tc.expected, found)
		}
	}
}
