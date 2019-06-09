package lex

// NewLexer returns a new lexer and runs the lexer
func NewLexer(input string) chan Item {
	l := &lexer{
		input:  input,
		column: 1,
		offset: 0,
		line:   1,
		pos:    0,
		items:  make(chan Item),
	}

	go l.run()

	return l.items
}

type lexer struct {
	input  string
	column int
	offset int
	line   int
	width  int
	pos    int
	items  chan Item
}

func (l lexer) emit(i itemType, value string) {
	l.items <- item{i, value, l.column, l.offset, l.line}
}

func (l lexer) emitText() {
	l.emit(ItemText, l.input[l.offset:l.pos])
}

func (l lexer) emitNewLine() {
	l.emit(ItemNewLine, "\n")
}

func (l lexer) emitEOF() {
	l.items <- item{ItemEOF, "", l.pos, l.pos, l.line}
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
				if val == ItemNewLine {
					if inTextBlock {
						inTextBlock = false
						l.emitText()
						l.column = l.column + l.width
						l.width = 0
					}
					l.offset = idx
					l.emitNewLine()
					l.column = 1
					// if this is the last item, don't increase the line count
					if idx+1 != len(l.input) {
						l.line = l.line + 1
					}
					continue
				}

				if inTextBlock {
					inTextBlock = false
					l.emitText()
					l.column = l.column + l.width
					l.width = 0
				}

				l.offset = idx
				l.emit(val, string(char))
				l.column = l.column + 1
			} else if !inTextBlock {
				if l.width == 0 {
					l.width = 1
				}
				l.offset = idx
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

func (l *lexer) nextItem() Item {
	return <-l.items
}
