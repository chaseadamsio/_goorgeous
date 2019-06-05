package goorgeous

type itemType int

const (
	itemNewLine itemType = iota
	itemAsterisk
	itemTilde
	itemForwardSlash
	itemUnderscore
	itemPlus
	itemColon
	itemSpace
	itemBracket
	itemBacktick
	itemParenthesis
	itemEqual
	itemPipe
	itemDash
	itemHash
	itemText
	itemEOF
)

var charToItem = map[rune]itemType{
	'\n': itemNewLine,
	'*':  itemAsterisk,
	'~':  itemTilde,
	'/':  itemForwardSlash,
	'_':  itemUnderscore,
	'+':  itemPlus,
	':':  itemColon,
	' ':  itemSpace,
	'`':  itemBacktick,
	'[':  itemBracket,
	']':  itemBracket,
	'(':  itemParenthesis,
	')':  itemParenthesis,
	'=':  itemEqual,
	'|':  itemPipe,
	'-':  itemDash,
	'#':  itemHash,
}

type item struct {
	typ    itemType
	val    string
	Column int
	Offset int
	Line   int
}

type lexer struct {
	input  string
	Column int
	Offset int
	Line   int
	width  int
	pos    int
	items  chan item
}

func (l lexer) emit(i itemType, value string) {
	l.items <- item{i, value, l.Column, l.Offset, l.Line}
}

func (l lexer) emitText() {
	l.emit(itemText, l.input[l.Offset:l.pos])
}

func (l lexer) emitNewLine() {
	l.emit(itemNewLine, "\n")
}

func (l lexer) emitEOF() {
	l.items <- item{itemEOF, "", l.pos, l.pos, l.Line}
}

// TODO maybe rename this to tokenizer
func (l *lexer) run() {
	inTextBlock := false
	var char rune
	for idx := 0; idx <= len(l.input); idx++ {
		l.pos = idx
		if idx != len(l.input) {
			char = rune(l.input[idx])
			if val, found := charToItem[char]; found {
				if val == itemNewLine {
					if inTextBlock {
						inTextBlock = false
						l.emitText()
						l.Column = l.Column + l.width
						l.width = 0
					}
					l.Offset = idx
					l.emitNewLine()
					l.Column = 1
					// if this is the last item, don't increase the line count
					if idx+1 != len(l.input) {
						l.Line = l.Line + 1
					}
					continue
				}

				if inTextBlock {
					inTextBlock = false
					l.emitText()
					l.Column = l.Column + l.width
					l.width = 0
				}

				l.Offset = idx
				l.emit(val, string(char))
				l.Column = l.Column + 1
			} else if !inTextBlock {
				if l.width == 0 {
					l.width = 1
				}
				l.Offset = idx
				inTextBlock = true
			} else {
				l.width = l.width + 1
			}
		} else {
			if inTextBlock {
				inTextBlock = false
				l.emitText()
			}
		}
	}
	l.emitEOF()
	close(l.items)
}

func (l *lexer) nextItem() item {
	return <-l.items
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		Column: 1,
		Offset: 0,
		Line:   1,
		pos:    0,
		items:  make(chan item),
	}

	go l.run()

	return l
}
