package parse

// item is an instance of a token with
// the type of the token as well as the start
// and end position of the the tokenized value
type item struct {
	typ   itemType
	val   []byte
	start int
	end   int
}

// Lexer is an input lexer that
// tokenizes the input for org mode syntax
type Lexer struct {
	input []byte
	state stateFn
	start int
	pos   int
	items chan item
}

const eof = -1

// stateFn is a state function that takes a lexer
// and returns another state function
type stateFn func(*Lexer) stateFn

// NewLexer returns a new Lexer based on
// some input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: []byte(input),
		items: make(chan item),
	}
	go l.run()
	return l
}

// run loops through stateFns until it receives a nil
// stateFn and then closes the items channel
func (l *Lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// emitItemText emits item text if it the current
// pos is greater than the start (as is the case)
// when some text has been passed over before finding
// a token
func emitItemText(l *Lexer) {
	if l.pos > l.start {
		l.emit(itemText)
	}
	if l.pos+1 <= len(l.input) {
		l.pos++
	}
}

// lexText looks for an identifier for a possible
// token and returns the lexer for that token
func lexText(l *Lexer) stateFn {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '\n':
			emitItemText(l)
			return lexNewLine
		default:
			l.pos++
		}
	}
	emitItemText(l)
	l.emit(itemEOF)
	return nil
}

// lexNewline emits a newline and returns
// the lexText stateFn
func lexNewLine(l *Lexer) stateFn {
	l.emit(itemNewline)
	return lexText
}

// emit takes an itemType for the current
// token to emit and sends a new item to the
// items channel with the start and end position
// of the item
func (l *Lexer) emit(typ itemType) {
	l.items <- item{
		typ:   typ,
		val:   l.input[l.start:l.pos],
		start: l.start,
		end:   l.pos,
	}
	l.start = l.pos
}

// nextItem takes values off of the channel and returns
// the items to the caller
func (l *Lexer) nextItem() item {
	item := <-l.items
	return item
}
