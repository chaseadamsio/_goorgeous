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

	itemEmphasis
	itemBold
	itemStrikethrough
	itemVerbatim
	itemCode
	itemUnderline

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
	// set when we find a open character so we know if we've found a closing character
	emphasisOpen      bool
	boldOpen          bool
	strikethroughOpen bool
	verbatimOpen      bool
	codeOpen          bool
	underlineOpen     bool
	linkOrImgOpen     bool
}

func (l *lexer) resetInlineOpeners() {
	l.emphasisOpen = false
	l.boldOpen = false
	l.strikethroughOpen = false
	l.codeOpen = false
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
	i := l.pos + 1
	for int(i) < len(l.input) && l.input[i] != '\n' {
		if delim == l.input[i] {
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
	delimNewLine = '\n'

	delimH1 = "*"
	delimH2 = "**"
	delimH3 = "***"
	delimH4 = "****"
	delimH5 = "*****"
	delimH6 = "******"

	delimEmphasis      = '/'
	delimBold          = '*'
	delimStrikethrough = '+'
	delimVerbatim      = '='
	delimCode          = '~'
	delimUnderline     = '_'
)

func lexText(l *lexer) stateFn {
	for int(l.pos) < len(l.input) {
		switch {
		case l.isNewLine():
			return lexNewLine
		case l.isBoldCandidate():
			return lexBold
		case l.isEmphasisCandidate():
			return lexEmphasis
		case l.isStrikethroughCandidate():
			return lexStrikethrough
		case l.isVerbatimCandidate():
			return lexVerbatim
		case l.isCodeCandidate():
			return lexCode
		case l.isUnderlineCandidate():
			return lexUnderline
		case l.isHeadlineCandidate():
			return lexHeadline
		default:
			l.next()
		}
	}

	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func (l *lexer) isNewLine() bool {
	return l.input[l.pos] == delimNewLine
}

func lexNewLine(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.resetInlineOpeners()
	l.next()
	l.emit(itemNewLine)
	return lexText
}

func isOtherDelim(char byte) bool {
	return char == delimBold || char == delimEmphasis || char == delimStrikethrough
}

func (l *lexer) isInlinePreChar() bool {
	if l.pos-1 >= 0 {
		delim := l.input[l.pos]
		char := l.input[l.pos-1]

		if delim == char {
			return false
		}

		if isOtherDelim(char) {
			return true
		}
		return charMatches(char, ' ') || charMatches(char, '>') || charMatches(char, '(') || charMatches(char, '{') || charMatches(char, '[')
	}
	return true
}

func (l *lexer) isInlineTerminatingChar() bool {
	if int(l.pos)+1 < len(l.input) {
		delim := l.input[l.pos]
		char := l.input[l.pos+1]

		if delim == char {
			return false
		}

		if isOtherDelim(char) {
			return true
		}
		return charMatches(char, ' ') || charMatches(char, '.') || charMatches(char, ',') || charMatches(char, '?') || charMatches(char, '!') || charMatches(char, ')') || charMatches(char, '}') || charMatches(char, ']')
	}
	return true
}

func (l *lexer) isHeadlineCandidate() bool {

	if l.input[l.pos] != '*' {
		return false
	}

	// headlines don't start with a space & shouldn't has a previous character
	// with a * if they've made it to this point
	if l.pos-1 >= 0 && (l.input[l.pos-1] == ' ' || l.input[l.pos-1] == '*') {
		return false
	}

	var i = 0
	for i <= 6 && int(l.pos)+i < len(l.input) && l.input[int(l.pos)+i] == '*' {
		i++
	}

	if int(l.pos)+i < len(l.input) {
		return l.input[int(l.pos)+i] == ' '
	}

	return false
}

func lexHeadline(l *lexer) stateFn {
	var i = 0
	for i <= 6 && int(l.pos)+i < len(l.input) && l.input[int(l.pos)+i] == '*' {
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

func (l *lexer) isEmphasisCandidate() bool {
	if l.input[l.pos] != delimEmphasis {
		return false
	}

	if !l.emphasisOpen && !l.hasMatch(delimEmphasis) {
		return false
	}

	if !l.emphasisOpen && l.isInlinePreChar() {
		l.emphasisOpen = true
		return true
	}

	if l.emphasisOpen && l.isInlineTerminatingChar() {
		l.emphasisOpen = false
		return true
	}

	return false
}

func lexEmphasis(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(itemEmphasis)
	return lexText
}

func (l *lexer) isBoldCandidate() bool {
	if l.input[l.pos] != delimBold {
		return false
	}

	if !l.boldOpen && int(l.pos)+1 < len(l.input) {
		if l.input[l.pos+1] == ' ' || l.input[l.pos+1] == '*' {
			return false
		}
	}

	if !l.boldOpen && !l.hasMatch(delimBold) {
		return false
	}

	if !l.boldOpen && l.isInlinePreChar() {
		l.boldOpen = true
		return true
	}

	if l.boldOpen && l.isInlineTerminatingChar() {
		l.boldOpen = false
		return true
	}

	return false
}

func lexBold(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(itemBold)
	return lexText
}

func (l *lexer) isStrikethroughCandidate() bool {
	if l.input[l.pos] != delimStrikethrough {
		return false
	}

	if !l.strikethroughOpen && !l.hasMatch(delimStrikethrough) {
		return false
	}

	if !l.strikethroughOpen && l.isInlinePreChar() {
		l.strikethroughOpen = true
		return true
	}

	if l.strikethroughOpen && l.isInlineTerminatingChar() {
		l.strikethroughOpen = false
		return true
	}

	return false
}

func lexStrikethrough(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(itemStrikethrough)
	return lexText
}

func (l *lexer) isVerbatimCandidate() bool {
	if l.input[l.pos] != delimVerbatim {
		return false
	}

	if !l.verbatimOpen && !l.hasMatch(delimVerbatim) {
		return false
	}

	if !l.verbatimOpen && l.isInlinePreChar() {
		l.verbatimOpen = true
		return true
	}

	if l.verbatimOpen && l.isInlineTerminatingChar() {
		l.verbatimOpen = false
		return true
	}

	return false
}

func lexVerbatim(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(itemVerbatim)
	return lexText
}

func (l *lexer) isCodeCandidate() bool {
	if l.input[l.pos] != delimCode {
		return false
	}

	if !l.codeOpen && !l.hasMatch(delimCode) {
		return false
	}

	if !l.codeOpen && l.isInlinePreChar() {
		l.codeOpen = true
		return true
	}

	if l.codeOpen && l.isInlineTerminatingChar() {
		l.codeOpen = false
		return true
	}

	return false
}

func lexCode(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(itemCode)
	return lexText
}

func (l *lexer) isUnderlineCandidate() bool {
	if l.input[l.pos] != delimUnderline {
		return false
	}

	if !l.underlineOpen && !l.hasMatch(delimUnderline) {
		return false
	}

	if !l.underlineOpen && l.isInlinePreChar() {
		l.underlineOpen = true
		return true
	}

	if l.underlineOpen && l.isInlineTerminatingChar() {
		l.underlineOpen = false
		return true
	}

	return false
}

func lexUnderline(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(itemUnderline)
	return lexText
}
