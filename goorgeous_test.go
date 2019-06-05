package goorgeous

func mkItem(typ itemType, text string, start, offset, line int) item {
	return item{
		typ, text, start, offset, line,
	}
}

func mkNewline(pos, offset, line int) item {
	return mkItem(itemNewLine, "\n", pos, offset, line)
}

func mkAsterisk(pos, offset, line int) item {
	return mkItem(itemAsterisk, "*", pos, offset, line)
}

func mkTilde(pos, offset, line int) item {
	return mkItem(itemTilde, "~", pos, offset, line)
}

func mkForwardSlash(pos, offset, line int) item {
	return mkItem(itemForwardSlash, "/", pos, offset, line)
}

func mkUnderscore(pos, offset, line int) item {
	return mkItem(itemUnderscore, "_", pos, offset, line)
}

func mkPlus(pos, offset, line int) item {
	return mkItem(itemPlus, "+", pos, offset, line)
}

func mkColon(pos, offset, line int) item {
	return mkItem(itemColon, ":", pos, offset, line)
}

func mkSpace(pos, offset, line int) item {
	return mkItem(itemSpace, " ", pos, offset, line)
}

func mkBacktick(pos, offset, line int) item {
	return mkItem(itemBacktick, "`", pos, offset, line)
}

func mkOpenBracket(pos, offset, line int) item {
	return mkItem(itemBracket, "[", pos, offset, line)
}

func mkCloseBracket(pos, offset, line int) item {
	return mkItem(itemBracket, "]", pos, offset, line)
}

func mkOpenParenthesis(pos, offset, line int) item {
	return mkItem(itemParenthesis, "(", pos, offset, line)
}

func mkCloseParenthesis(pos, offset, line int) item {
	return mkItem(itemParenthesis, ")", pos, offset, line)
}

func mkEqual(pos, offset, line int) item {
	return mkItem(itemEqual, "=", pos, offset, line)
}

func mkDash(pos, offset, line int) item {
	return mkItem(itemDash, "-", pos, offset, line)
}

func mkHash(pos, offset, line int) item {
	return mkItem(itemHash, "#", pos, offset, line)
}

func mkPipe(pos, offset, line int) item {
	return mkItem(itemPipe, "|", pos, offset, line)
}

func mkText(val string, pos, offset, line int) item {
	return mkItem(itemText, val, pos, offset, line)
}

func mkEOF(pos, offset, line int) item {
	return mkItem(itemEOF, "", pos, offset, line)
}

type testCase struct {
	name  string
	input string
	items []item
}

var lexTests = []testCase{
	{
		"headline - level 1",
		"* this is a headline",
		[]item{
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
		[]item{
			mkAsterisk(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("this", 3, 2, 1),
			mkSpace(7, 6, 1),
			mkText("is", 8, 7, 1),
			mkSpace(10, 9, 1),
			mkText("a", 11, 10, 1),
			mkSpace(12, 11, 1),
			mkText("headline", 13, 12, 1),
			mkNewline(21, 20, 1),
			mkEOF(21, 21, 1),
		},
	},
	{
		"link",
		"[this is a link](https://github.com)",
		[]item{
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
		[]item{
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
			mkNewline(37, 36, 1),
			mkEOF(37, 37, 1),
		},
	},
	{
		"complex",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC",
		[]item{
			mkAsterisk(1, 0, 1),
			mkAsterisk(2, 1, 1),
			mkSpace(3, 2, 1),
			mkText("hello", 4, 3, 1),
			mkNewline(9, 8, 1),
			mkText("this", 1, 9, 2),
			mkSpace(5, 13, 2),
			mkText("is", 6, 14, 2),
			mkSpace(8, 16, 2),
			mkText("some", 9, 17, 2),
			mkSpace(13, 21, 2),
			mkText("text", 14, 22, 2),
			mkNewline(18, 26, 2),
			mkHash(1, 27, 3),
			mkPlus(2, 28, 3),
			mkText("BEGIN", 3, 29, 3),
			mkUnderscore(8, 34, 3),
			mkText("SRC", 9, 35, 3),
			mkSpace(12, 38, 3),
			mkText("javascript", 13, 39, 3),
			mkNewline(23, 49, 3),
			mkText("console.log", 1, 50, 4),
			mkOpenParenthesis(12, 61, 4),
			mkText("\"hello", 13, 62, 4),
			mkSpace(19, 68, 4),
			mkText("world\"", 20, 69, 4),
			mkCloseParenthesis(26, 75, 4),
			mkText(";", 27, 76, 4),
			mkNewline(28, 77, 4),
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
		[]item{
			mkAsterisk(1, 0, 1),
			mkAsterisk(2, 1, 1),
			mkSpace(3, 2, 1),
			mkText("hello", 4, 3, 1),
			mkNewline(9, 8, 1),
			mkText("this", 1, 9, 2),
			mkSpace(5, 13, 2),
			mkText("is", 6, 14, 2),
			mkSpace(8, 16, 2),
			mkText("some", 9, 17, 2),
			mkSpace(13, 21, 2),
			mkText("text", 14, 22, 2),
			mkNewline(18, 26, 2),
			mkHash(1, 27, 3),
			mkPlus(2, 28, 3),
			mkText("BEGIN", 3, 29, 3),
			mkUnderscore(8, 34, 3),
			mkText("SRC", 9, 35, 3),
			mkSpace(12, 38, 3),
			mkText("javascript", 13, 39, 3),
			mkNewline(23, 49, 3),
			mkText("console.log", 1, 50, 4),
			mkOpenParenthesis(12, 61, 4),
			mkText("\"hello", 13, 62, 4),
			mkSpace(19, 68, 4),
			mkText("world\"", 20, 69, 4),
			mkCloseParenthesis(26, 75, 4),
			mkText(";", 27, 76, 4),
			mkNewline(28, 77, 4),
			mkHash(1, 78, 5),
			mkPlus(2, 79, 5),
			mkText("END", 3, 80, 5),
			mkUnderscore(6, 83, 5),
			mkText("SRC", 7, 84, 5),
			mkNewline(10, 87, 5),
			mkEOF(88, 88, 5),
		},
	},
	{
		"complex w/ trailing text",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\nhello",
		[]item{
			mkAsterisk(1, 0, 1),
			mkAsterisk(2, 1, 1),
			mkSpace(3, 2, 1),
			mkText("hello", 4, 3, 1),
			mkNewline(9, 8, 1),
			mkText("this", 1, 9, 2),
			mkSpace(5, 13, 2),
			mkText("is", 6, 14, 2),
			mkSpace(8, 16, 2),
			mkText("some", 9, 17, 2),
			mkSpace(13, 21, 2),
			mkText("text", 14, 22, 2),
			mkNewline(18, 26, 2),
			mkHash(1, 27, 3),
			mkPlus(2, 28, 3),
			mkText("BEGIN", 3, 29, 3),
			mkUnderscore(8, 34, 3),
			mkText("SRC", 9, 35, 3),
			mkSpace(12, 38, 3),
			mkText("javascript", 13, 39, 3),
			mkNewline(23, 49, 3),
			mkText("console.log", 1, 50, 4),
			mkOpenParenthesis(12, 61, 4),
			mkText("\"hello", 13, 62, 4),
			mkSpace(19, 68, 4),
			mkText("world\"", 20, 69, 4),
			mkCloseParenthesis(26, 75, 4),
			mkText(";", 27, 76, 4),
			mkNewline(28, 77, 4),
			mkHash(1, 78, 5),
			mkPlus(2, 79, 5),
			mkText("END", 3, 80, 5),
			mkUnderscore(6, 83, 5),
			mkText("SRC", 7, 84, 5),
			mkNewline(10, 87, 5),
			mkText("hello", 1, 88, 6),
			mkEOF(93, 93, 6),
		},
	},
	{
		"list",
		"- apples\n- oranges\n- bananas\nsomething else",
		[]item{
			mkDash(1, 0, 1),
			mkSpace(2, 1, 1),
			mkText("apples", 3, 2, 1),
			mkNewline(9, 8, 1),
			mkDash(1, 9, 2),
			mkSpace(2, 10, 2),
			mkText("oranges", 3, 11, 2),
			mkNewline(10, 18, 2),
			mkDash(1, 19, 3),
			mkSpace(2, 20, 3),
			mkText("bananas", 3, 21, 3),
			mkNewline(10, 28, 3),
			mkText("something", 1, 29, 4),
			mkSpace(10, 38, 4),
			mkText("else", 11, 39, 4),
			mkEOF(43, 43, 4),
		},
	},
	{
		"table",
		"| Name  | Phone | Age |\n|-------+-------+-----|\n| Peter |  1234 |  17 |\n| Anna  |  4321 |  25 |\n",
		[]item{
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
			mkNewline(24, 23, 1),
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
			mkNewline(24, 47, 2),
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
			mkNewline(24, 71, 3),
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
			mkNewline(24, 95, 4),
			mkEOF(96, 96, 4),
		},
	},
}
