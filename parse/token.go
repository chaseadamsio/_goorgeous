package parse

import "strconv"

// tokenType identifies the type of lex items
type tokenType int

const (
	tokenError tokenType = iota
	tokenNewline
	tokenSpace
	tokenEOF
	tokenWord
	tokenAsterisk // "*" indicates either a headline or a bold token
	tokenHash     // "#  " indicates a comment token
	tokenPlus
	tokenSlash
	tokenEqual
	tokenTilde
	tokenUnderscore
	tokenDash
	tokenColon
	tokenBracketLeft
	tokenBracketRight
	tokenPipe
)

var elTypes = [...]string{
	tokenError:        "tokenError",
	tokenNewline:      "tokenNewline",
	tokenEOF:          "tokenEOF",
	tokenWord:         "tokenText",
	tokenAsterisk:     "*",
	tokenHash:         "#",
	tokenPlus:         "+",
	tokenSlash:        "/",
	tokenEqual:        "=",
	tokenTilde:        "~",
	tokenUnderscore:   "_",
	tokenDash:         "-",
	tokenColon:        ":",
	tokenBracketLeft:  "[",
	tokenBracketRight: "]",
	tokenPipe:         "|",
}

func (typ tokenType) String() string {
	s := ""
	if 0 <= typ && typ < tokenType(len(elTypes)) {
		s = elTypes[typ]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(typ)) + ")"
	}
	return s
}
