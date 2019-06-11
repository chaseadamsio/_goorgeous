package lex

import (
	"fmt"
	"testing"
)

func mkItem(typ itemType, text string, start, offset, line int) item {
	return item{
		typ, text, start, offset, line,
	}
}

func mkNewLine(pos, offset, line int) item {
	return mkItem(ItemNewLine, "\n", pos, offset, line)
}

func mkAsterisk(pos, offset, line int) item {
	return mkItem(ItemAsterisk, "*", pos, offset, line)
}

func mkTilde(pos, offset, line int) item {
	return mkItem(ItemTilde, "~", pos, offset, line)
}

func mkForwardSlash(pos, offset, line int) item {
	return mkItem(ItemForwardSlash, "/", pos, offset, line)
}

func mkUnderscore(pos, offset, line int) item {
	return mkItem(ItemUnderscore, "_", pos, offset, line)
}

func mkPlus(pos, offset, line int) item {
	return mkItem(ItemPlus, "+", pos, offset, line)
}

func mkColon(pos, offset, line int) item {
	return mkItem(ItemColon, ":", pos, offset, line)
}

func mkSpace(pos, offset, line int) item {
	return mkItem(ItemSpace, " ", pos, offset, line)
}

func mkBacktick(pos, offset, line int) item {
	return mkItem(ItemBacktick, "`", pos, offset, line)
}

func mkOpenBracket(pos, offset, line int) item {
	return mkItem(ItemBracket, "[", pos, offset, line)
}

func mkCloseBracket(pos, offset, line int) item {
	return mkItem(ItemBracket, "]", pos, offset, line)
}

func mkOpenParenthesis(pos, offset, line int) item {
	return mkItem(ItemParenthesis, "(", pos, offset, line)
}

func mkCloseParenthesis(pos, offset, line int) item {
	return mkItem(ItemParenthesis, ")", pos, offset, line)
}

func mkEqual(pos, offset, line int) item {
	return mkItem(ItemEqual, "=", pos, offset, line)
}

func mkDash(pos, offset, line int) item {
	return mkItem(ItemDash, "-", pos, offset, line)
}

func mkHash(pos, offset, line int) item {
	return mkItem(ItemHash, "#", pos, offset, line)
}

func mkPipe(pos, offset, line int) item {
	return mkItem(ItemPipe, "|", pos, offset, line)
}

func mkText(val string, pos, offset, line int) item {
	return mkItem(ItemText, val, pos, offset, line)
}

func mkTab(pos, offset, line int) item {
	return mkItem(ItemTab, "	", pos, offset, line)
}

func mkEOF(pos, offset, line int) item {
	return mkItem(ItemEOF, "", pos, offset, line)
}

func mkItemsText(items []Item) {
	fmt.Println("[]Item{")
	for _, item := range items {
		if item.Type() == ItemText {
			fmt.Printf("mkText(\"%s\", %d, %d, %d),\n", item.Value(), item.Column(), item.Offset(), item.Line())
		} else {
			fmt.Printf("mk%s(%d, %d, %d),\n", item.Type(), item.Column(), item.Offset(), item.Line())
		}
	}
	fmt.Println("},")
}

func equal(i1, i2 []Item, checkPos bool, t *testing.T) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].Type() != i2[k].Type() {
			t.Logf("types not equal: %s: %T, %s: %T", i1[k].Value(), i1[k].Type(), i2[k].Value(), i2[k].Type())
			return false
		}
		if i1[k].Value() != i2[k].Value() {
			t.Logf("vals not equal: %s, %s", i1[k].Value(), i2[k].Value())
			return false
		}
		if checkPos {
			if i1[k].Column() != i2[k].Column() {
				t.Logf("column not equal: %d, %d", i1[k].Column(), i2[k].Column())
				return false
			}
			if i1[k].Offset() != i2[k].Offset() {
				t.Logf("offset not equal: %d, %d", i1[k].Offset(), i2[k].Offset())
				return false
			}
			if i1[k].Line() != i2[k].Line() {
				t.Logf("line not equal: %d, %d", i1[k].Line(), i2[k].Line())
				return false
			}
		}
	}
	return true
}
func TestLex(t *testing.T) {
	for _, tc := range lexTests {
		t.Run(tc.name, func(t *testing.T) {
			var foundItems []Item
			items := NewLexer(tc.input)
			for item := range items {
				foundItems = append(foundItems, item)
			}
			if !equal(foundItems, tc.items, true, t) {
				t.Errorf("got\n\t%v\nexpected\n\t%v", foundItems, tc.items)

				// Output a passing items list for easily adding new tests
				mkItemsText(foundItems)
			}

		})
	}
}

type testCase struct {
	name  string
	input string
	items []Item
}

var lexTests = []testCase{
	{
		"headers",
		"#+title: my org mode content\n#+author: Chase Adams\n#+description: This is my description!",
		[]Item{
			mkHash(1, 0, 1),
			mkPlus(2, 1, 1),
			mkText("title", 3, 2, 1),
			mkColon(8, 7, 1),
			mkSpace(9, 8, 1),
			mkText("my", 10, 9, 1),
			mkSpace(12, 11, 1),
			mkText("org", 13, 12, 1),
			mkSpace(16, 15, 1),
			mkText("mode", 17, 16, 1),
			mkSpace(21, 20, 1),
			mkText("content", 22, 21, 1),
			mkNewLine(29, 28, 1),
			mkHash(1, 29, 2),
			mkPlus(2, 30, 2),
			mkText("author", 3, 31, 2),
			mkColon(9, 37, 2),
			mkSpace(10, 38, 2),
			mkText("Chase", 11, 39, 2),
			mkSpace(16, 44, 2),
			mkText("Adams", 17, 45, 2),
			mkNewLine(22, 50, 2),
			mkHash(1, 51, 3),
			mkPlus(2, 52, 3),
			mkText("description", 3, 53, 3),
			mkColon(14, 64, 3),
			mkSpace(15, 65, 3),
			mkText("This", 16, 66, 3),
			mkSpace(20, 70, 3),
			mkText("is", 21, 71, 3),
			mkSpace(23, 73, 3),
			mkText("my", 24, 74, 3),
			mkSpace(26, 76, 3),
			mkText("description!", 27, 77, 3),
			mkEOF(89, 89, 3),
		},
	},

	{
		"basic-happy-path-new-content-after",
		"#+title: my org mode content\n#+author: Chase Adams\n#+description: This is my description!\n* This starts the content!",
		[]Item{
			mkHash(1, 0, 1),
			mkPlus(2, 1, 1),
			mkText("title", 3, 2, 1),
			mkColon(8, 7, 1),
			mkSpace(9, 8, 1),
			mkText("my", 10, 9, 1),
			mkSpace(12, 11, 1),
			mkText("org", 13, 12, 1),
			mkSpace(16, 15, 1),
			mkText("mode", 17, 16, 1),
			mkSpace(21, 20, 1),
			mkText("content", 22, 21, 1),
			mkNewLine(29, 28, 1),
			mkHash(1, 29, 2),
			mkPlus(2, 30, 2),
			mkText("author", 3, 31, 2),
			mkColon(9, 37, 2),
			mkSpace(10, 38, 2),
			mkText("Chase", 11, 39, 2),
			mkSpace(16, 44, 2),
			mkText("Adams", 17, 45, 2),
			mkNewLine(22, 50, 2),
			mkHash(1, 51, 3),
			mkPlus(2, 52, 3),
			mkText("description", 3, 53, 3),
			mkColon(14, 64, 3),
			mkSpace(15, 65, 3),
			mkText("This", 16, 66, 3),
			mkSpace(20, 70, 3),
			mkText("is", 21, 71, 3),
			mkSpace(23, 73, 3),
			mkText("my", 24, 74, 3),
			mkSpace(26, 76, 3),
			mkText("description!", 27, 77, 3),
			mkNewLine(39, 89, 3),
			mkAsterisk(1, 90, 4),
			mkSpace(2, 91, 4),
			mkText("This", 3, 92, 4),
			mkSpace(7, 96, 4),
			mkText("starts", 8, 97, 4),
			mkSpace(14, 103, 4),
			mkText("the", 15, 104, 4),
			mkSpace(18, 107, 4),
			mkText("content!", 19, 108, 4),
			mkEOF(116, 116, 4),
		},
	},
	{
		"basic-happy-path-with-tags",
		"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+tags: org-content org-mode hugo\n",
		[]Item{
			mkHash(1, 0, 1),
			mkPlus(2, 1, 1),
			mkText("title", 3, 2, 1),
			mkColon(8, 7, 1),
			mkSpace(9, 8, 1),
			mkText("my", 10, 9, 1),
			mkSpace(12, 11, 1),
			mkText("org", 13, 12, 1),
			mkSpace(16, 15, 1),
			mkText("mode", 17, 16, 1),
			mkSpace(21, 20, 1),
			mkText("tags", 22, 21, 1),
			mkSpace(26, 25, 1),
			mkText("content", 27, 26, 1),
			mkNewLine(34, 33, 1),
			mkHash(1, 34, 2),
			mkPlus(2, 35, 2),
			mkText("author", 3, 36, 2),
			mkColon(9, 42, 2),
			mkSpace(10, 43, 2),
			mkText("Chase", 11, 44, 2),
			mkSpace(16, 49, 2),
			mkText("Adams", 17, 50, 2),
			mkNewLine(22, 55, 2),
			mkHash(1, 56, 3),
			mkPlus(2, 57, 3),
			mkText("description", 3, 58, 3),
			mkColon(14, 69, 3),
			mkSpace(15, 70, 3),
			mkText("This", 16, 71, 3),
			mkSpace(20, 75, 3),
			mkText("is", 21, 76, 3),
			mkSpace(23, 78, 3),
			mkText("my", 24, 79, 3),
			mkSpace(26, 81, 3),
			mkText("description!", 27, 82, 3),
			mkNewLine(39, 94, 3),
			mkHash(1, 95, 4),
			mkPlus(2, 96, 4),
			mkText("tags", 3, 97, 4),
			mkColon(7, 101, 4),
			mkSpace(8, 102, 4),
			mkText("org", 9, 103, 4),
			mkDash(12, 106, 4),
			mkText("content", 13, 107, 4),
			mkSpace(20, 114, 4),
			mkText("org", 21, 115, 4),
			mkDash(24, 118, 4),
			mkText("mode", 25, 119, 4),
			mkSpace(29, 123, 4),
			mkText("hugo", 30, 124, 4),
			mkNewLine(34, 128, 4),
			mkEOF(129, 129, 4),
		},
	},

	{
		"basic-happy-path-with-categories",
		"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+categories: org-content org-mode hugo\n",
		[]Item{
			mkHash(1, 0, 1),
			mkPlus(2, 1, 1),
			mkText("title", 3, 2, 1),
			mkColon(8, 7, 1),
			mkSpace(9, 8, 1),
			mkText("my", 10, 9, 1),
			mkSpace(12, 11, 1),
			mkText("org", 13, 12, 1),
			mkSpace(16, 15, 1),
			mkText("mode", 17, 16, 1),
			mkSpace(21, 20, 1),
			mkText("tags", 22, 21, 1),
			mkSpace(26, 25, 1),
			mkText("content", 27, 26, 1),
			mkNewLine(34, 33, 1),
			mkHash(1, 34, 2),
			mkPlus(2, 35, 2),
			mkText("author", 3, 36, 2),
			mkColon(9, 42, 2),
			mkSpace(10, 43, 2),
			mkText("Chase", 11, 44, 2),
			mkSpace(16, 49, 2),
			mkText("Adams", 17, 50, 2),
			mkNewLine(22, 55, 2),
			mkHash(1, 56, 3),
			mkPlus(2, 57, 3),
			mkText("description", 3, 58, 3),
			mkColon(14, 69, 3),
			mkSpace(15, 70, 3),
			mkText("This", 16, 71, 3),
			mkSpace(20, 75, 3),
			mkText("is", 21, 76, 3),
			mkSpace(23, 78, 3),
			mkText("my", 24, 79, 3),
			mkSpace(26, 81, 3),
			mkText("description!", 27, 82, 3),
			mkNewLine(39, 94, 3),
			mkHash(1, 95, 4),
			mkPlus(2, 96, 4),
			mkText("categories", 3, 97, 4),
			mkColon(13, 107, 4),
			mkSpace(14, 108, 4),
			mkText("org", 15, 109, 4),
			mkDash(18, 112, 4),
			mkText("content", 19, 113, 4),
			mkSpace(26, 120, 4),
			mkText("org", 27, 121, 4),
			mkDash(30, 124, 4),
			mkText("mode", 31, 125, 4),
			mkSpace(35, 129, 4),
			mkText("hugo", 36, 130, 4),
			mkNewLine(40, 134, 4),
			mkEOF(135, 135, 4),
		},
	},
	{
		"basic-happy-path-with-aliases",
		"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+aliases: /org/content /org/mode /hugo\n",
		[]Item{
			mkHash(1, 0, 1),
			mkPlus(2, 1, 1),
			mkText("title", 3, 2, 1),
			mkColon(8, 7, 1),
			mkSpace(9, 8, 1),
			mkText("my", 10, 9, 1),
			mkSpace(12, 11, 1),
			mkText("org", 13, 12, 1),
			mkSpace(16, 15, 1),
			mkText("mode", 17, 16, 1),
			mkSpace(21, 20, 1),
			mkText("tags", 22, 21, 1),
			mkSpace(26, 25, 1),
			mkText("content", 27, 26, 1),
			mkNewLine(34, 33, 1),
			mkHash(1, 34, 2),
			mkPlus(2, 35, 2),
			mkText("author", 3, 36, 2),
			mkColon(9, 42, 2),
			mkSpace(10, 43, 2),
			mkText("Chase", 11, 44, 2),
			mkSpace(16, 49, 2),
			mkText("Adams", 17, 50, 2),
			mkNewLine(22, 55, 2),
			mkHash(1, 56, 3),
			mkPlus(2, 57, 3),
			mkText("description", 3, 58, 3),
			mkColon(14, 69, 3),
			mkSpace(15, 70, 3),
			mkText("This", 16, 71, 3),
			mkSpace(20, 75, 3),
			mkText("is", 21, 76, 3),
			mkSpace(23, 78, 3),
			mkText("my", 24, 79, 3),
			mkSpace(26, 81, 3),
			mkText("description!", 27, 82, 3),
			mkNewLine(39, 94, 3),
			mkHash(1, 95, 4),
			mkPlus(2, 96, 4),
			mkText("aliases", 3, 97, 4),
			mkColon(10, 104, 4),
			mkSpace(11, 105, 4),
			mkForwardSlash(12, 106, 4),
			mkText("org", 13, 107, 4),
			mkForwardSlash(16, 110, 4),
			mkText("content", 17, 111, 4),
			mkSpace(24, 118, 4),
			mkForwardSlash(25, 119, 4),
			mkText("org", 26, 120, 4),
			mkForwardSlash(29, 123, 4),
			mkText("mode", 30, 124, 4),
			mkSpace(34, 128, 4),
			mkForwardSlash(35, 129, 4),
			mkText("hugo", 36, 130, 4),
			mkNewLine(40, 134, 4),
			mkEOF(135, 135, 4),
		},
	},
	{
		"basic - text",
		"this is a line.\nthis is a newline.",
		[]Item{
			mkText("this", 1, 0, 1),
			mkSpace(5, 4, 1),
			mkText("is", 6, 5, 1),
			mkSpace(8, 7, 1),
			mkText("a", 9, 8, 1),
			mkSpace(10, 9, 1),
			mkText("line.", 11, 10, 1),
			mkNewLine(16, 15, 1),
			mkText("this", 1, 16, 2),
			mkSpace(5, 20, 2),
			mkText("is", 6, 21, 2),
			mkSpace(8, 23, 2),
			mkText("a", 9, 24, 2),
			mkSpace(10, 25, 2),
			mkText("newline.", 11, 26, 2),
			mkEOF(34, 34, 2),
		},
	},
	{
		"headline - level 1",
		"* this is a headline",
		[]Item{
			mkAsterisk(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("this", 3, 2, 1),
			mkSpace(7, 6, 1),
			mkText("is", 8, 7, 1),
			mkSpace(10, 9, 1),
			mkText("a", 11, 10, 1),
			mkSpace(12, 11, 1),
			mkText("headline", 13, 12, 1),
			mkEOF(20, 20, 1),
		},
	},
	{
		"headline - level 1 w/ newline",
		"* this is a headline\n",
		[]Item{
			mkAsterisk(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("this", 3, 2, 1),
			mkSpace(7, 6, 1),
			mkText("is", 8, 7, 1),
			mkSpace(10, 9, 1),
			mkText("a", 11, 10, 1),
			mkSpace(12, 11, 1),
			mkText("headline", 13, 12, 1),
			mkNewLine(21, 20, 1),
			mkEOF(21, 21, 1),
		},
	},
	{
		"headline - deep",
		"* headline1\n** headline2\n*** headline3\n* headline1-2\n",
		[]Item{
			mkAsterisk(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("headline1", 3, 2, 1),
			mkNewLine(12, 11, 1),
			mkAsterisk(1, 12, 2),
			mkAsterisk(2, 13, 2),
			mkSpace(3, 14, 2),
			mkText("headline2", 4, 15, 2),
			mkNewLine(13, 24, 2),
			mkAsterisk(1, 25, 3),
			mkAsterisk(2, 26, 3),
			mkAsterisk(3, 27, 3),
			mkSpace(4, 28, 3),
			mkText("headline3", 5, 29, 3),
			mkNewLine(14, 38, 3),
			mkAsterisk(1, 39, 4),
			mkSpace(2, 40, 4),
			mkText("headline1", 3, 41, 4),
			mkDash(12, 50, 4),
			mkText("2", 13, 51, 4),
			mkNewLine(14, 52, 4),
			mkEOF(53, 53, 4),
		},
	},
	{
		"link",
		"[this is a link](https://github.com)",
		[]Item{
			mkOpenBracket(1, 0, 1),
			mkText("this", 2, 1, 1),
			mkSpace(6, 5, 1),
			mkText("is", 7, 6, 1),
			mkSpace(9, 8, 1),
			mkText("a", 10, 9, 1),
			mkSpace(11, 10, 1),
			mkText("link", 12, 11, 1),
			mkCloseBracket(16, 15, 1),
			mkOpenParenthesis(17, 16, 1),
			mkText("https", 18, 17, 1),
			mkColon(23, 22, 1),
			mkForwardSlash(24, 23, 1),
			mkForwardSlash(25, 24, 1),
			mkText("github.com", 26, 25, 1),
			mkCloseParenthesis(36, 35, 1),
			mkEOF(36, 36, 1),
		},
	},
	{
		"link w/ newline",
		"[this is a link](https://github.com)\n",
		[]Item{
			mkOpenBracket(1, 0, 1),
			mkText("this", 2, 1, 1),
			mkSpace(6, 5, 1),
			mkText("is", 7, 6, 1),
			mkSpace(9, 8, 1),
			mkText("a", 10, 9, 1),
			mkSpace(11, 10, 1),
			mkText("link", 12, 11, 1),
			mkCloseBracket(16, 15, 1),
			mkOpenParenthesis(17, 16, 1),
			mkText("https", 18, 17, 1),
			mkColon(23, 22, 1),
			mkForwardSlash(24, 23, 1),
			mkForwardSlash(25, 24, 1),
			mkText("github.com", 26, 25, 1),
			mkCloseParenthesis(36, 35, 1),
			mkNewLine(37, 36, 1),
			mkEOF(37, 37, 1),
		},
	},
	{
		"complex",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC",
		[]Item{
			mkAsterisk(1, 0, 1),
			mkAsterisk(2, 1, 1),
			mkSpace(3, 2, 1),
			mkText("hello", 4, 3, 1),
			mkNewLine(9, 8, 1),
			mkText("this", 1, 9, 2),
			mkSpace(5, 13, 2),
			mkText("is", 6, 14, 2),
			mkSpace(8, 16, 2),
			mkText("some", 9, 17, 2),
			mkSpace(13, 21, 2),
			mkText("text", 14, 22, 2),
			mkNewLine(18, 26, 2),
			mkHash(1, 27, 3),
			mkPlus(2, 28, 3),
			mkText("BEGIN", 3, 29, 3),
			mkUnderscore(8, 34, 3),
			mkText("SRC", 9, 35, 3),
			mkSpace(12, 38, 3),
			mkText("javascript", 13, 39, 3),
			mkNewLine(23, 49, 3),
			mkText("console.log", 1, 50, 4),
			mkOpenParenthesis(12, 61, 4),
			mkText("\"hello", 13, 62, 4),
			mkSpace(19, 68, 4),
			mkText("world\"", 20, 69, 4),
			mkCloseParenthesis(26, 75, 4),
			mkText(";", 27, 76, 4),
			mkNewLine(28, 77, 4),
			mkHash(1, 78, 5),
			mkPlus(2, 79, 5),
			mkText("END", 3, 80, 5),
			mkUnderscore(6, 83, 5),
			mkText("SRC", 7, 84, 5),
			mkEOF(87, 87, 5),
		},
	},
	{
		"complex w/ newline",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\n",
		[]Item{
			mkAsterisk(1, 0, 1),
			mkAsterisk(2, 1, 1),
			mkSpace(3, 2, 1),
			mkText("hello", 4, 3, 1),
			mkNewLine(9, 8, 1),
			mkText("this", 1, 9, 2),
			mkSpace(5, 13, 2),
			mkText("is", 6, 14, 2),
			mkSpace(8, 16, 2),
			mkText("some", 9, 17, 2),
			mkSpace(13, 21, 2),
			mkText("text", 14, 22, 2),
			mkNewLine(18, 26, 2),
			mkHash(1, 27, 3),
			mkPlus(2, 28, 3),
			mkText("BEGIN", 3, 29, 3),
			mkUnderscore(8, 34, 3),
			mkText("SRC", 9, 35, 3),
			mkSpace(12, 38, 3),
			mkText("javascript", 13, 39, 3),
			mkNewLine(23, 49, 3),
			mkText("console.log", 1, 50, 4),
			mkOpenParenthesis(12, 61, 4),
			mkText("\"hello", 13, 62, 4),
			mkSpace(19, 68, 4),
			mkText("world\"", 20, 69, 4),
			mkCloseParenthesis(26, 75, 4),
			mkText(";", 27, 76, 4),
			mkNewLine(28, 77, 4),
			mkHash(1, 78, 5),
			mkPlus(2, 79, 5),
			mkText("END", 3, 80, 5),
			mkUnderscore(6, 83, 5),
			mkText("SRC", 7, 84, 5),
			mkNewLine(10, 87, 5),
			mkEOF(88, 88, 5),
		},
	},
	{
		"complex w/ trailing text",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\nhello",
		[]Item{
			mkAsterisk(1, 0, 1),
			mkAsterisk(2, 1, 1),
			mkSpace(3, 2, 1),
			mkText("hello", 4, 3, 1),
			mkNewLine(9, 8, 1),
			mkText("this", 1, 9, 2),
			mkSpace(5, 13, 2),
			mkText("is", 6, 14, 2),
			mkSpace(8, 16, 2),
			mkText("some", 9, 17, 2),
			mkSpace(13, 21, 2),
			mkText("text", 14, 22, 2),
			mkNewLine(18, 26, 2),
			mkHash(1, 27, 3),
			mkPlus(2, 28, 3),
			mkText("BEGIN", 3, 29, 3),
			mkUnderscore(8, 34, 3),
			mkText("SRC", 9, 35, 3),
			mkSpace(12, 38, 3),
			mkText("javascript", 13, 39, 3),
			mkNewLine(23, 49, 3),
			mkText("console.log", 1, 50, 4),
			mkOpenParenthesis(12, 61, 4),
			mkText("\"hello", 13, 62, 4),
			mkSpace(19, 68, 4),
			mkText("world\"", 20, 69, 4),
			mkCloseParenthesis(26, 75, 4),
			mkText(";", 27, 76, 4),
			mkNewLine(28, 77, 4),
			mkHash(1, 78, 5),
			mkPlus(2, 79, 5),
			mkText("END", 3, 80, 5),
			mkUnderscore(6, 83, 5),
			mkText("SRC", 7, 84, 5),
			mkNewLine(10, 87, 5),
			mkText("hello", 1, 88, 6),
			mkEOF(93, 93, 6),
		},
	},
	{
		"unordered-list",
		"- apples\n- oranges\n- bananas\nsomething else",
		[]Item{
			mkDash(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("apples", 3, 2, 1),
			mkNewLine(9, 8, 1),
			mkDash(1, 9, 2),
			mkSpace(2, 10, 2),
			mkText("oranges", 3, 11, 2),
			mkNewLine(10, 18, 2),
			mkDash(1, 19, 3),
			mkSpace(2, 20, 3),
			mkText("bananas", 3, 21, 3),
			mkNewLine(10, 28, 3),
			mkText("something", 1, 29, 4),
			mkSpace(10, 38, 4),
			mkText("else", 11, 39, 4),
			mkEOF(43, 43, 4),
		},
	},
	{
		"unordered-list-with-child-ordered-list",
		"- apples\n\t1. in apples 1\n\t2. in apples 2\n\t3. in apples 3\n- oranges\n- bananas\nsomething else",
		[]Item{
			mkDash(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("apples", 3, 2, 1),
			mkNewLine(9, 8, 1),
			mkTab(1, 9, 2),
			mkText("1.", 2, 10, 2),
			mkSpace(4, 12, 2),
			mkText("in", 5, 13, 2),
			mkSpace(7, 15, 2),
			mkText("apples", 8, 16, 2),
			mkSpace(14, 22, 2),
			mkText("1", 15, 23, 2),
			mkNewLine(16, 24, 2),
			mkTab(1, 25, 3),
			mkText("2.", 2, 26, 3),
			mkSpace(4, 28, 3),
			mkText("in", 5, 29, 3),
			mkSpace(7, 31, 3),
			mkText("apples", 8, 32, 3),
			mkSpace(14, 38, 3),
			mkText("2", 15, 39, 3),
			mkNewLine(16, 40, 3),
			mkTab(1, 41, 4),
			mkText("3.", 2, 42, 4),
			mkSpace(4, 44, 4),
			mkText("in", 5, 45, 4),
			mkSpace(7, 47, 4),
			mkText("apples", 8, 48, 4),
			mkSpace(14, 54, 4),
			mkText("3", 15, 55, 4),
			mkNewLine(16, 56, 4),
			mkDash(1, 57, 5),
			mkSpace(2, 58, 5),
			mkText("oranges", 3, 59, 5),
			mkNewLine(10, 66, 5),
			mkDash(1, 67, 6),
			mkSpace(2, 68, 6),
			mkText("bananas", 3, 69, 6),
			mkNewLine(10, 76, 6),
			mkText("something", 1, 77, 7),
			mkSpace(10, 86, 7),
			mkText("else", 11, 87, 7),
			mkEOF(91, 91, 7),
		},
	},
	{
		"table",
		"| Name  | Phone | Age |\n|-------+-------+-----|\n| Peter |  1234 |  17 |\n| Anna  |  4321 |  25 |\n",
		[]Item{
			mkPipe(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("Name", 3, 2, 1),
			mkSpace(7, 6, 1),
			mkSpace(8, 7, 1),
			mkPipe(9, 8, 1),
			mkSpace(10, 9, 1),
			mkText("Phone", 11, 10, 1),
			mkSpace(16, 15, 1),
			mkPipe(17, 16, 1),
			mkSpace(18, 17, 1),
			mkText("Age", 19, 18, 1),
			mkSpace(22, 21, 1),
			mkPipe(23, 22, 1),
			mkNewLine(24, 23, 1),
			mkPipe(1, 24, 2),
			mkDash(2, 25, 2),
			mkDash(3, 26, 2),
			mkDash(4, 27, 2),
			mkDash(5, 28, 2),
			mkDash(6, 29, 2),
			mkDash(7, 30, 2),
			mkDash(8, 31, 2),
			mkPlus(9, 32, 2),
			mkDash(10, 33, 2),
			mkDash(11, 34, 2),
			mkDash(12, 35, 2),
			mkDash(13, 36, 2),
			mkDash(14, 37, 2),
			mkDash(15, 38, 2),
			mkDash(16, 39, 2),
			mkPlus(17, 40, 2),
			mkDash(18, 41, 2),
			mkDash(19, 42, 2),
			mkDash(20, 43, 2),
			mkDash(21, 44, 2),
			mkDash(22, 45, 2),
			mkPipe(23, 46, 2),
			mkNewLine(24, 47, 2),
			mkPipe(1, 48, 3),
			mkSpace(2, 49, 3),
			mkText("Peter", 3, 50, 3),
			mkSpace(8, 55, 3),
			mkPipe(9, 56, 3),
			mkSpace(10, 57, 3),
			mkSpace(11, 58, 3),
			mkText("1234", 12, 59, 3),
			mkSpace(16, 63, 3),
			mkPipe(17, 64, 3),
			mkSpace(18, 65, 3),
			mkSpace(19, 66, 3),
			mkText("17", 20, 67, 3),
			mkSpace(22, 69, 3),
			mkPipe(23, 70, 3),
			mkNewLine(24, 71, 3),
			mkPipe(1, 72, 4),
			mkSpace(2, 73, 4),
			mkText("Anna", 3, 74, 4),
			mkSpace(7, 78, 4),
			mkSpace(8, 79, 4),
			mkPipe(9, 80, 4),
			mkSpace(10, 81, 4),
			mkSpace(11, 82, 4),
			mkText("4321", 12, 83, 4),
			mkSpace(16, 87, 4),
			mkPipe(17, 88, 4),
			mkSpace(18, 89, 4),
			mkSpace(19, 90, 4),
			mkText("25", 20, 91, 4),
			mkSpace(22, 93, 4),
			mkPipe(23, 94, 4),
			mkNewLine(24, 95, 4),
			mkEOF(96, 96, 4),
		},
	},
}
