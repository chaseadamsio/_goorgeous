package lex

import "regexp"

// Item represents a parsed token
type Item interface {
	checker
	Type() itemType
	Value() string
	Column() int
	Offset() int
	End() int
	Line() int
}

type checker interface {
	newlineChecker
	asteriskChecker
	tildeChecker
	forwardSlashChecker
	underscoreChecker
	plusChecker
	colonChecker
	spaceChecker
	bracketChecker
	backtickChecker
	parenthesisChecker
	equalChecker
	pipeChecker
	dashChecker
	textChecker
	hashChecker
	tabChecker
	eofChecker
	whitespaceChecker
	wordChecker
	nonwordChecker
}

type newlineChecker interface {
	IsNewline() bool
}

type asteriskChecker interface {
	IsAsterisk() bool
}

type tildeChecker interface {
	IsTilde() bool
}

type underscoreChecker interface {
	IsUnderscore() bool
}

type plusChecker interface {
	IsPlus() bool
}

type colonChecker interface {
	IsColon() bool
}

type spaceChecker interface {
	IsSpace() bool
}

type bracketChecker interface {
	IsBracket() bool
}

type backtickChecker interface {
	IsBacktick() bool
}

type parenthesisChecker interface {
	IsParenthesis() bool
}

type equalChecker interface {
	IsEqual() bool
}

type pipeChecker interface {
	IsPipe() bool
}

type dashChecker interface {
	IsDash() bool
}

type hashChecker interface {
	IsHash() bool
}

type textChecker interface {
	IsText() bool
}

type forwardSlashChecker interface {
	IsForwardSlash() bool
}

type tabChecker interface {
	IsTab() bool
}

type eofChecker interface {
	IsEOF() bool
}

type whitespaceChecker interface {
	IsWhitespace() bool
}

type wordChecker interface {
	IsWord() bool
}

type nonwordChecker interface {
	IsNonWord() bool
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
	case ItemTab:
		val = "Tab"
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
	// ItemTab is a Tab Item
	ItemTab
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
	'	': ItemTab,
	'`': ItemBacktick,
	'[': ItemBracket,
	']': ItemBracket,
	'(': ItemParenthesis,
	')': ItemParenthesis,
	'=': ItemEqual,
	'|': ItemPipe,
	'-': ItemDash,
	'#': ItemHash,
}

type item struct {
	typ    itemType
	val    string
	column int
	offset int
	width  int
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

func (i item) End() int {
	return i.offset + i.width
}

func (i item) Line() int {
	return i.line
}

func (i item) IsNewline() bool {
	return i.typ == ItemNewLine
}

func (i item) IsAsterisk() bool {
	return i.typ == ItemAsterisk
}

func (i item) IsTilde() bool {
	return i.typ == ItemTilde
}

func (i item) IsForwardSlash() bool {
	return i.typ == ItemForwardSlash
}

func (i item) IsUnderscore() bool {
	return i.typ == ItemUnderscore
}

func (i item) IsPlus() bool {
	return i.typ == ItemPlus
}

func (i item) IsColon() bool {
	return i.typ == ItemColon
}

func (i item) IsSpace() bool {
	return i.typ == ItemSpace
}

func (i item) IsBracket() bool {
	return i.typ == ItemBracket
}

func (i item) IsBacktick() bool {
	return i.typ == ItemBacktick
}

func (i item) IsParenthesis() bool {
	return i.typ == ItemParenthesis
}

func (i item) IsEqual() bool {
	return i.typ == ItemEqual
}

func (i item) IsPipe() bool {
	return i.typ == ItemPipe
}

func (i item) IsDash() bool {
	return i.typ == ItemDash
}

func (i item) IsHash() bool {
	return i.typ == ItemHash
}

func (i item) IsText() bool {
	return i.typ == ItemText
}

func (i item) IsTab() bool {
	return i.typ == ItemTab
}

func (i item) IsEOF() bool {
	return i.typ == ItemEOF
}

func (i item) IsWhitespace() bool {
	return i.typ == ItemSpace
}

func (i item) IsWord() bool {
	matched, err := regexp.Match(`\w`, []byte(i.Value()))
	if err != nil {
		panic(err) // TODO ¯\_(ツ)_/¯ handle this better?
	}
	return matched
}

func (i item) IsNonWord() bool {
	matched, err := regexp.Match(`\W`, []byte(i.Value()))
	if err != nil {
		panic(err) // TODO ¯\_(ツ)_/¯ handle this better?
	}
	return matched
}
