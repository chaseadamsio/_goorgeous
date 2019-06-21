package parse

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/lex"
)

func TestIsHorizontalRule(t *testing.T) {
	var testCases = []struct {
		source string
		start  int
		end    int
		found  bool
	}{
		{"-----", 0, 5, true},
		{" -----", 1, 6, true},
		{"-", 0, -1, false},
		{"--", 0, -1, false},
		{"---", 0, -1, false},
		{"----", 0, -1, false},
		{"------", 0, -1, false},
		{" ------", 1, -1, false},
		{"not ------", 2, -1, false},
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

			found, end := p.matchesHorizontalRule(tc.start)
			if tc.end != end {
				t.Errorf("expected %s end to be %d. got %d", tc.source, tc.end, end)
			}
			if tc.found != found {
				t.Errorf("expected %s found to be %t", tc.source, tc.found)
			}

		})
	}
}
