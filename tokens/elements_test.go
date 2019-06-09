package tokens

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
)

func TestIsBold(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"*this is bold.*", true},
		{"*_this is bold._*", true},
		{"*this is b*old, with a star in the middle*", true},
		{"_*this is bold._*", false},
		{"*this is \n not bold.*", false},
		{"* this is not bold.*", false},
		{"*this is not bold. *", false},
		{"*this is not *bold. *", false},
		{"*this is not b*old. ", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsBold(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestIsItalic(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"/this is italic./", true},
		{"/this is ita/lic, with a star in the middle/", true},
		{"/this is \n not italic./", false},
		{"/ this is not italic./", false},
		{"/this is not italic. /", false},
		{"/this is not /italic. /", false},
		{"/this is not ita/lic. ", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsItalic(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestIsVerbatim(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"=this is verbatim.=", true},
		{"=this is verbat=im, with an equal in the middle=", true},
		{"=this is \n not verbatim.=", false},
		{"= this is not verbatim.=", false},
		{"=this is not verbatim. =", false},
		{"=this is not =verbatim. =", false},
		{"=this is not verbat=im. ", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsVerbatim(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestIsStrikeThrough(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"+this is strike through.+", true},
		{"+this is strike+through, with a plus in the middle+", true},
		{"+this is \n not strike through.+", false},
		{"+ this is not strike through.+", false},
		{"+this is not strike through. +", false},
		{"+this is not +strike through. +", false},
		{"+this is not verbat+im. ", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsStrikeThrough(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestIsUnderline(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"_this is underline._", true},
		{"_this is under_line, with a plus in the middle_", true},
		{"_this is \n not underline._", false},
		{"_ this is not underline._", false},
		{"_this is not underline. _", false},
		{"_this is not _underline. _", false},
		{"_this is not verbat_im. ", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsUnderline(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}

func TestIsCode(t *testing.T) {
	testCases := []struct {
		value    string
		expected bool
	}{
		{"~this is Code.~", true},
		{"~this is Co~de, with a plus in the middle~", true},
		{"~this is \n not Code.~", false},
		{"~ this is not Code.~", false},
		{"~this is not Code. ~", false},
		{"~this is not ~Code. ~", false},
		{"~this is not Cod~e. ", false},
	}

	for _, tc := range testCases {
		var items []lex.Item
		lexedItems := lex.NewLexer(tc.value)
		for item := range lexedItems {
			items = append(items, item)
		}
		if IsCode(items) != tc.expected {
			t.Errorf("expected \"%s\" to be %t", tc.value, tc.expected)
		}
	}
}
