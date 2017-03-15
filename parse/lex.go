package parse

type item struct {
	typ itemType
	pos int
	val string
}

type itemType int

const (
	itemNewLine itemType = iota
	itemEOF
	itemText
)

type stateFn func(*lexer) stateFn

type lexer struct {
	pos, start int
	input      string
	state      stateFn
	items      chan item
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}

	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

const (
	delimNewLine = '\n'
)

func lexText(l *lexer) stateFn {
	if l.pos < len(l.input) {
		switch l.input[l.pos] {
		case delimNewLine:
			lexNewLine(l)
			return lexText
		default:
			l.pos += 1
			return lexText
		}
	}
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexNewLine(l *lexer) {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.pos += 1
	l.emit(itemNewLine)
}
