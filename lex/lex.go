package lex

// Item represents a parsed token
type Item interface {
	checker
	Type() itemType
	Value() string
	Column() int
	Offset() int
	Line() int
}

type checker interface {
	asteriskChecker
	newlineChecker
	eofChecker
}

type asteriskChecker interface {
	IsAsterisk() bool
}

type newlineChecker interface {
	IsNewline() bool
}

type eofChecker interface {
	IsEOF() bool
}

type itemType int

func (i itemType) String() string {
	var val string
	switch i {
	case ItemNewLine:
		val = "NewLine"
	case ItemAsterisk:
		val = "Asterisk"
	case ItemTilde:
		val = "Tilde"
	case ItemForwardSlash:
		val = "ForwardSlash"
	case ItemUnderscore:
		val = "Underscore"
	case ItemPlus:
		val = "Plus"
	case ItemColon:
		val = "Colon"
	case ItemSpace:
		val = "Space"
	case ItemBracket:
		val = "Bracket"
	case ItemBacktick:
		val = "Backtick"
	case ItemParenthesis:
		val = "Parenthesis"
	case ItemEqual:
		val = "Equal"
	case ItemPipe:
		val = "Pipe"
	case ItemDash:
		val = "Dash"
	case ItemHash:
		val = "Hash"
	case ItemEOF:
		val = "EOF"
	}
	return val
}

const (
	// ItemNewLine is a New Line Item
	ItemNewLine itemType = iota
	// ItemAsterisk is an Asterisk Item
	ItemAsterisk
	// ItemTilde is a Tilde Item
	ItemTilde
	// ItemForwardSlash is a ForwardSlash Item
	ItemForwardSlash
	// ItemUnderscore is a Underscore Item
	ItemUnderscore
	// ItemPlus is a Plus Item
	ItemPlus
	// ItemColon is a Colon Item
	ItemColon
	// ItemSpace is a Space Item
	ItemSpace
	// ItemBracket is a Bracket Item
	ItemBracket
	// ItemBacktick is a Backtick Item
	ItemBacktick
	// ItemParenthesis is a Parenthesis Item
	ItemParenthesis
	// ItemEqual is a Equal Item
	ItemEqual
	// ItemPipe is a Pipe Item
	ItemPipe
	// ItemDash is a Dash Item
	ItemDash
	// ItemHash is a Hash Item
	ItemHash
	// ItemText is a Text Item
	ItemText
	// ItemEOF is a EOF Item
	ItemEOF
)

var charToItem = map[rune]itemType{
	'\n': ItemNewLine,
	'*':  ItemAsterisk,
	'~':  ItemTilde,
	'/':  ItemForwardSlash,
	'_':  ItemUnderscore,
	'+':  ItemPlus,
	':':  ItemColon,
	' ':  ItemSpace,
	'`':  ItemBacktick,
	'[':  ItemBracket,
	']':  ItemBracket,
	'(':  ItemParenthesis,
	')':  ItemParenthesis,
	'=':  ItemEqual,
	'|':  ItemPipe,
	'-':  ItemDash,
	'#':  ItemHash,
}

type item struct {
	typ    itemType
	val    string
	column int
	offset int
	line   int
}

func (i item) Type() itemType {
	return i.typ
}

func (i item) Value() string {
	return i.val
}

func (i item) Column() int {
	return i.column
}

func (i item) Offset() int {
	return i.offset
}

func (i item) Line() int {
	return i.line
}

func (i item) IsAsterisk() bool {
	return i.typ == ItemAsterisk
}

func (i item) IsNewline() bool {
	return i.typ == ItemNewLine
}

func (i item) IsEOF() bool {
	return i.typ == ItemEOF
}

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
