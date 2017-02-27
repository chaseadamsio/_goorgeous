package goorgeous

import (
	"fmt"
	"unicode/utf8"
)

// itemType is a lexical token of Org Mode
type itemType int

const (
	itemEOF itemType = iota
	itemWS
	itemNewLine

	itemHeadline
	itemH1
	itemH2
	itemH3
	itemH4
	itemH5
	itemH6

	itemStatus
	itemPriority
	itemTags
	itemPropertyDrawer

	itemKeyword
	itemComment
	itemHorizontalRule

	itemVerbatim
	itemCode
	itemEmphasis
	itemUpderline
	itemBold
	itemStrikeThrough

	itemLink
	itemImage

	itemDefinitionList
	itemOrderedList
	itemUnorderedList

	itemTable

	itemHTML

	itemText
)

const eof = -1

const (
	spaceChars = " \t\r\n"
)

type Pos int

func (p Pos) Position() Pos {
	return p
}

type item struct {
	typ  itemType
	pos  Pos
	val  string
	line int
}

func (i item) String() string {
	return fmt.Sprintf("%q", i.val)
}

type stateFn func(*lexer) stateFn

type lexer struct {
	name    string
	input   string
	start   Pos
	pos     Pos
	width   Pos
	lastPos Pos
	line    int
	items   chan item
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
		line:  1,
	}

	go l.run()
	return l
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) hasMatch(delim byte) bool {
	i := l.pos
	for i != '\n' && int(i) < len(l.input) {
		if delim == l.input[i] {
			if l.pos+1 == i || l.pos == i {
				return false
			}
			return true
		}
		i++
	}
	return false
}

func (l *lexer) backup() {
	l.pos -= l.width
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos], l.line}
	l.start = l.pos
}

func (l *lexer) next() rune {
	var r rune

	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width

	if r == '\n' {
		l.line++
	}
	return r
}

func (l *lexer) current() rune {
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

const (
	delimNewLine = "\n"

	delimH1       = "*"
	delimH2       = "**"
	delimH3       = "***"
	delimH4       = "****"
	delimH5       = "*****"
	delimH6       = "******"
	delimEmphasis = "/"
)

func lexText(l *lexer) stateFn {

	for int(l.pos) < len(l.input) {
		switch r := l.next(); {
		case r == '*':
			l.backup()
			// only true for headlines
			if l.pos == 0 && l.input[l.pos] != ' ' {
				return lexHeadline
			}
			l.next()
			if l.hasMatch('*') {
				return lexBold
			}
			l.next()
			continue
		case r == '/':
			return lexEmphasis
		case r == '\n':
			l.backup()
			l.emit(itemText)
			l.next()
			return lexNewLine
		}
	}

	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexNewLine(l *lexer) stateFn {
	l.emit(itemNewLine)
	if int(l.pos) < len(l.input) {
		l.pos += 1
	}
	return lexText
}

func lexHeadline(l *lexer) stateFn {
	var i = 0
	for i <= 6 && l.input[int(l.pos)+i] == '*' {
		i++
	}
	l.pos += Pos(i)
	if l.peek() != ' ' {
		return lexText
	}
	l.emit(itemHeadline)
	return lexInsideHeadline
}

func lexInsideHeadline(l *lexer) stateFn {
	return lexText
}

func lexEmphasis(l *lexer) stateFn {
	if l.pos > 1 {
		l.backup()
		l.emit(itemText)
		l.next()
	}
	l.emit(itemEmphasis)
	return lexText
}

func lexBold(l *lexer) stateFn {
	if l.pos > 1 {
		l.backup()
		l.emit(itemText)
		l.next()
	}
	l.emit(itemBold)
	return lexText
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}
