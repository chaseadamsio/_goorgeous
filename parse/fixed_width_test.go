package parse

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
)

func TestMatchesFixedWidth(t *testing.T) {
	var testCases = []struct {
		source string
		start  int
		end    int
		found  bool
	}{
		{": Fixed width area", 0, 8, true},
		{": Fixed width area\n: still fixed width area", 0, 18, true},
		{"		: Fixed width area", 2, 10, true},
		{"not		: a Fixed width area", 3, -1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.source, func(t *testing.T) {

			var items []lex.Item
			lexedItems := lex.NewLexer(tc.source)
			for item := range lexedItems {
				items = append(items, item)
			}

			p := &parser{
				items: items,
				depth: 0,
			}

			found, end := p.matchesFixedWidth(tc.start)
			if tc.end != end {
				t.Errorf("expected %s end to be %d. got %d", tc.source, tc.end, end)
			}
			if tc.found != found {
				t.Errorf("expected %s found to be %t", tc.source, tc.found)
			}

		})
	}
}
