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
	itemProperty

	itemKeyword
	itemComment
	itemHorizontalRule

	itemEmphasis
	itemBold
	itemStrikethrough
	itemVerbatim
	itemCode
	itemUnderline

	itemImgOrLinkOpen
	itemImgOrLinkOpenSingle
	itemImgPre
	itemImgOrLinkClose
	itemImgOrLinkCloseSingle

	itemDefinitionTerm
	itemDefinitionDescription
	itemOrderedList
	// [@<d>] to set the value of an ordered list item
	itemOrderedListNumber
	itemUnorderedList
	itemListItem

	itemTable

	itemBlock

	itemHTML

	itemText
)

const eof = -1

type Pos int

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
	isOpen map[string]bool
}

func (l *lexer) resetInlineOpeners() {
	for k := range l.isOpen {
		l.isOpen[k] = false
	}
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
		line:  1,
		isOpen: map[string]bool{
			"emphasis":      false,
			"bold":          false,
			"strikethrough": false,
			"verbatim":      false,
			"code":          false,
			"underline":     false,
		},
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

func (l *lexer) skip() {
	l.start = l.pos
}

func (l *lexer) skipSpace() {
	if int(l.pos) < len(l.input) {
		if l.input[l.pos] == ' ' {
			l.start += 1
			l.pos += 1
		}
	}
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

	delimTags = ':'

	delimEmphasis      = '/'
	delimBold          = '*'
	delimStrikethrough = '+'
	delimVerbatim      = '='
	delimCode          = '~'
	delimUnderline     = '_'

	delimPropertyBegin = ":PROPERTIES:"
	delimPropertyEnd   = ":END:"
)

func lexText(l *lexer) stateFn {
	for int(l.pos) < len(l.input) {
		switch {
		case l.isNewLine():
			return lexNewLine
		case l.isCommentCandidate():
			return lexComment
		case l.isBlockCandidate():
			return lexBlock
		case l.isTableCandidate():
			return lexTable
		case l.isOrderedListCandidate():
			return lexOrderedList
		case l.isDefinitionListCandidate():
			return lexDefinitionList
		case l.isUnorderedListCandidate():
			return lexUnorderedList
		case l.isImgOrLinkCandidate():
			return lexImgOrLink
		case l.isBoldCandidate():
			return lexBold
		case l.isEmphasisCandidate():
			return lexEmphasis(l)
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
		case l.isPropertyDrawerCandidate():
			return lexPropertyDrawer
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

func (l *lexer) isCommentCandidate() bool {
	if !(l.input[l.pos] == '#' && l.input[l.pos+1] == ' ') {
		return false
	}

	for idx := int(l.pos - 1); idx >= 0 && l.input[idx] != '\n'; idx-- {
		if l.input[idx] != ' ' {
			return false
		}
	}

	return true
}

func lexComment(l *lexer) stateFn {
	idx := int(l.pos)
	for ; idx < len(l.input) && l.input[idx] != '\n'; idx++ {
		if l.input[idx] == '#' {
			break
		}
	}
	l.pos = Pos(idx + 1)
	l.emit(itemComment)
	return lexText
}

func (l *lexer) isBlockCandidate() bool {
	if !(l.input[l.pos] == '#' && l.input[l.pos+1] == '+') {
		return false
	}

	for idx := int(l.pos - 1); idx >= 0 && l.input[idx] != '\n'; idx-- {
		if l.input[idx] != ' ' {
			return false
		}
	}

	return true
}

func lexBlock(l *lexer) stateFn {
	idx := int(l.pos)
	for ; idx < len(l.input) && l.input[idx] != '\n'; idx++ {
		if l.input[idx] == ' ' {
			break
		}
	}
	l.pos = Pos(idx)
	l.emit(itemBlock)
	return lexText
}

func (l *lexer) isTableCandidate() bool {
	if l.input[l.pos] != '|' {
		return false
	}
	if l.pos-1 >= 0 && l.input[l.pos-1] == ' ' {
		return false
	}
	return true
}

func lexTable(l *lexer) stateFn {
	for idx := int(l.pos); idx < len(l.input); idx++ {
		if l.input[idx] == '|' {
			l.pos = Pos(idx)
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.next()
			l.emit(itemTable)
		}
		if l.input[idx] == '\n' {
			break
		}
	}
	return lexText
}

// numeral followed by either a period or a right parenthesis2, such as ‘1.’ or ‘1)`
func (l *lexer) isOrderedListCandidate() bool {
	if !isDigit(l.input[l.pos]) {
		return false
	}

	for idx := int(l.pos - 1); idx >= 0 && l.input[idx] != '\n'; idx-- {
		if l.input[idx] != ' ' {
			return false
		}
	}

	for idx := int(l.pos); idx < len(l.input); idx++ {
		if isDigit(l.input[idx]) {
			continue
		}
		if !(l.input[idx] == '.' || l.input[idx] == ')') {
			return false
		}
		break
	}
	return true
}

func isDigit(char byte) bool {
	return char <= '9' && char >= '0'
}

func lexOrderedList(l *lexer) stateFn {
	for idx := int(l.pos); idx < len(l.input) && l.input[idx] != '\n'; idx++ {
		if isDigit(l.input[idx]) {
			continue
		}
		if !(l.input[idx] == '.' || l.input[idx] == ')') {
			l.pos = Pos(idx)
			l.emit(itemOrderedList)
			if l.input[l.pos:l.pos+3] == " [@" {
				foundDigits := false
				for valIdx := int(l.pos) + 3; valIdx < len(l.input); valIdx++ {
					if l.input[valIdx] == ']' && foundDigits {
						l.pos = Pos(valIdx + 1)
						l.emit(itemOrderedListNumber)
						break
					}
					if isDigit(l.input[valIdx]) {
						foundDigits = true
						continue
					}
					return lexText
				}
			}
			break
		}
	}
	return lexText
}

func (l *lexer) isDefinitionListCandidate() bool {
	if l.input[l.pos] != '-' {
		return false
	}

	for idx := int(l.pos - 1); idx >= 0 && l.input[idx] != '\n'; idx-- {
		if l.input[idx] != ' ' {
			return false
		}
	}

	for idx := int(l.pos + 1); idx+4 < len(l.input) && l.input[idx] != '\n'; idx++ {
		if l.input[idx:idx+4] == " :: " {
			return true
		}
	}

	return false
}

func lexDefinitionList(l *lexer) stateFn {
	l.next()
	l.emit(itemDefinitionTerm)

	for idx := int(l.pos + 1); idx+2 < len(l.input) && l.input[idx] != '\n'; idx++ {
		if l.input[idx:idx+2] == "::" {
			l.pos = Pos(idx)
			l.emit(itemText)
			l.pos += 2
			l.emit(itemDefinitionDescription)
		}
	}
	return lexText
}

func (l *lexer) isUnorderedListCandidate() bool {
	if !(l.input[l.pos] == '-' || l.input[l.pos] == '+') {
		return false
	}

	for idx := int(l.pos - 1); idx >= 0 && l.input[idx] != '\n'; idx-- {
		if l.input[idx] != ' ' {
			return false
		}
	}

	if l.input[l.pos+1] != ' ' {
		return false
	}

	return true
}

func lexUnorderedList(l *lexer) stateFn {
	l.next()
	l.emit(itemUnorderedList)
	return lexText
}

func (l *lexer) isImgOrLinkCandidate() bool {
	if int(l.pos)+2 < len(l.input) {
		if l.input[l.pos:l.pos+2] != "[[" {
			return false
		}

		if testImgOrLink(int(l.pos+2), l.input) {
			return true
		}
	}
	return false
}

func testImgOrLink(start int, in string) bool {
	foundFirstCloseTag := false
	foundSecondOpenTag := false
	for idx := start; idx < len(in); idx++ {
		if in[idx] == ']' {
			if in[idx+1] == ']' {
				return true
			}

			if !foundFirstCloseTag {
				foundFirstCloseTag = true
				continue
			}

			if foundFirstCloseTag && foundSecondOpenTag {
				return true
			}
			return false
		}

		if in[idx] == '[' && in[idx-1] != ']' {
			return false
		}

		if in[idx] == '[' && in[idx-1] == ']' {
			if !foundFirstCloseTag {
				return false
			}
			foundSecondOpenTag = true
			continue
		}
	}
	return false
}

func lexImgOrLink(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.pos += 2
	l.emit(itemImgOrLinkOpen)
	if int(l.pos)+5 < len(l.input) && l.input[l.pos:l.pos+5] == "file:" {
		l.pos += 5
		l.emit(itemImgPre)
	}

	for idx := int(l.pos); idx < len(l.input); {
		if l.input[idx] == ']' {
			if l.input[idx+1] == ']' {
				l.pos = Pos(idx)
				if l.pos > l.start {
					l.emit(itemText)
					l.pos += 2
					l.emit(itemImgOrLinkClose)
				} else {
					l.pos += 2
					l.emit(itemImgOrLinkClose)
				}
				idx += 2
				continue
			} else {
				l.pos = Pos(idx)
				if l.pos > l.start {
					l.emit(itemText)
					l.next()
					l.emit(itemImgOrLinkCloseSingle)
				} else {
					l.next()
					l.emit(itemImgOrLinkCloseSingle)
				}

			}
		}

		if l.input[idx] == '[' {
			l.pos = Pos(idx)
			l.next()
			l.emit(itemImgOrLinkOpenSingle)
		}

		idx++
	}

	return lexText
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
	l.skipSpace()
	return lexInsideHeadline
}

func lexInsideHeadline(l *lexer) stateFn {
	// is there a todo status?
	if l.isStatus() {
		lexStatus(l)
	}
	if l.isPriority() {
		lexPriority(l)
	}
	if l.findTags() > int(l.pos) {
		lexTags(l)
	}
	return lexText
}

func (l *lexer) isStatus() bool {
	if int(l.pos)+4 < len(l.input) {
		return l.input[l.pos:l.pos+4] == "TODO" || l.input[l.pos:l.pos+4] == "DONE"
	}
	return false
}

func lexStatus(l *lexer) {
	l.pos += 4
	l.emit(itemStatus)
	l.skipSpace()
}

func (l *lexer) isPriority() bool {
	if int(l.pos)+3 < len(l.input) {
		return l.input[l.pos] == '[' && isPriorityLetter(l.input[l.pos+1]) && l.input[l.pos+2] == ']'
	}
	return false
}

func isPriorityLetter(char byte) bool {
	return (charMatches(char, 'A') || charMatches(char, 'B') || charMatches(char, 'C'))
}

func lexPriority(l *lexer) {
	l.pos += 3
	l.emit(itemPriority)
	l.skipSpace()
}

func (l *lexer) findTags() int {
	idx := int(l.pos)
	for idx < len(l.input) {
		char := l.input[idx]

		if char == '\n' {
			idx = int(l.pos)
			break
		}
		if char == delimTags {
			if l.hasMatch(delimTags) {
				return idx
			}
		}
		idx++
	}

	if idx == len(l.input) {
		idx = int(l.pos)
	}

	return idx
}

func lexTags(l *lexer) {
	for int(l.pos) < len(l.input) {
		if l.input[l.pos] == delimTags {
			l.skipSpace()
			l.emit(itemText)
			l.next()
			l.emit(itemTags)
		} else {
			l.next()
		}
	}
	l.skipSpace()
}

func (l *lexer) isEmphasisCandidate() bool {
	return l.checkDelimCandidate(delimEmphasis, "emphasis")
}

func lexEmphasis(l *lexer) stateFn {
	return lexInsideInline(l, itemEmphasis)
}

func (l *lexer) isBoldCandidate() bool {
	return l.checkDelimCandidate(delimBold, "bold")
}

func lexBold(l *lexer) stateFn {
	return lexInsideInline(l, itemBold)
}

func (l *lexer) isStrikethroughCandidate() bool {
	return l.checkDelimCandidate(delimStrikethrough, "strikethrough")
}

func lexStrikethrough(l *lexer) stateFn {
	return lexInsideInline(l, itemStrikethrough)
}

func (l *lexer) isVerbatimCandidate() bool {
	return l.checkDelimCandidate(delimVerbatim, "verbatim")
}

func lexVerbatim(l *lexer) stateFn {
	return lexInsideInline(l, itemVerbatim)
}

func (l *lexer) isCodeCandidate() bool {
	return l.checkDelimCandidate(delimCode, "code")
}

func lexCode(l *lexer) stateFn {
	return lexInsideInline(l, itemCode)
}

func (l *lexer) isUnderlineCandidate() bool {
	return l.checkDelimCandidate(delimUnderline, "underline")
}

func lexUnderline(l *lexer) stateFn {
	return lexInsideInline(l, itemUnderline)
}

func (l *lexer) checkDelimCandidate(delim byte, t string) bool {
	if l.input[l.pos] != delim {
		return false
	}

	// bold is the only one that might have a collision with another character
	// which is a headline
	if t == "bold" {
		if !l.isOpen[t] && int(l.pos)+1 < len(l.input) {
			if l.input[l.pos+1] == ' ' || l.input[l.pos+1] == '*' {
				return false
			}
		}

	}

	if !l.isOpen[t] && !l.hasMatch(delim) {
		return false
	}

	if !l.isOpen[t] && l.isInlinePreChar() {
		l.isOpen[t] = true
		return true
	}

	if l.isOpen[t] && l.isInlineTerminatingChar() {
		l.isOpen[t] = false
		return true
	}

	return false

}

func lexInsideInline(l *lexer, it itemType) stateFn {
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.next()
	l.emit(it)
	return lexText
}

func (l *lexer) isPropertyDrawerCandidate() bool {
	if (int(l.pos)+len(":PROPERTIES:") <= len(l.input) && l.input[l.pos:int(l.pos)+len(":PROPERTIES:")] == ":PROPERTIES:") || (int(l.pos)+len(":END:") <= len(l.input) && l.input[l.pos:int(l.pos)+len(":END:")] == ":END:") {
		return true
	}
	return false
}

func lexPropertyDrawer(l *lexer) stateFn {
	if l.input[l.pos:int(l.pos)+len(":END:")] == ":END:" {
		l.pos += 5
	} else if l.input[l.pos:int(l.pos)+len(":PROPERTIES:")] == ":PROPERTIES:" {
		l.pos += 12
	}
	l.emit(itemPropertyDrawer)
	return lexText
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
