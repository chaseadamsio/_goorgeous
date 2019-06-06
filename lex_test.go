package goorgeous

import (
	"fmt"
	"testing"
)

func equal(i1, i2 []item, checkPos bool, t *testing.T) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			t.Logf("types not equal: %s: %T, %s: %T", i1[k].val, i1[k].typ, i2[k].val, i2[k].typ)
			return false
		}
		if i1[k].val != i2[k].val {
			t.Logf("vals not equal: %s, %s", i1[k].val, i2[k].val)
			return false
		}
		if checkPos {
			if i1[k].Column != i2[k].Column {
				t.Logf("Column not equal: %d, %d", i1[k].Column, i2[k].Column)
				return false
			}
			if i1[k].Offset != i2[k].Offset {
				t.Logf("Offset not equal: %d, %d", i1[k].Offset, i2[k].Offset)
				return false
			}
			if i1[k].Line != i2[k].Line {
				t.Logf("Line not equal: %d, %d", i1[k].Line, i2[k].Line)
				return false
			}
		}
	}
	return true
}
func TestLex(t *testing.T) {
	for _, tc := range lexTests {
		t.Run(tc.name, func(t *testing.T) {
			var foundItems []item
			l := lex(tc.input)
			for item := range l.items {
				foundItems = append(foundItems, item)
			}
			if !equal(foundItems, tc.items, true, t) {
				t.Errorf("got\n\t%v\nexpected\n\t%v", foundItems, tc.items)

				mkItemsText(foundItems)
			}

		})
	}
}

func mkItemsText(items []item) {
	fmt.Println("[]item{")
	for _, item := range items {
		if item.typ == itemText {
			fmt.Printf("mkText(\"%s\", %d, %d, %d),\n", item.val, item.Column, item.Offset, item.Line)
		} else {
			fmt.Printf("mk%s(%d, %d, %d),\n", item.typ, item.Column, item.Offset, item.Line)
		}
	}
	fmt.Println("},")
}
