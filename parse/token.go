package parse

import "strconv"

// elType identifies the type of lex items
type elType int

const (
	elError elType = iota
	elNewline
	elSpace
	elEOF
	elWord
	elAsterisk // "*" indicates either a headline or a bold token
	elHash     // "#  " indicates a comment token
	elSlash
	elEqual
	elTilde
	elUnderscore
	elColon
	elBracketLeft
	elBracketRight
	elPipe
)

var elTypes = [...]string{
	elError:        "elError",
	elNewline:      "elNewline",
	elEOF:          "elEOF",
	elWord:         "elText",
	elAsterisk:     "*",
	elHash:         "#",
	elSlash:        "/",
	elEqual:        "=",
	elTilde:        "~",
	elUnderscore:   "_",
	elColon:        ":",
	elBracketLeft:  "[",
	elBracketRight: "]",
	elPipe:         "|",
}

func (typ elType) String() string {
	s := ""
	if 0 <= typ && typ < elType(len(elTypes)) {
		s = elTypes[typ]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(typ)) + ")"
	}
	return s
}
