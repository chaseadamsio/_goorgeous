package goorgeous

import (
	"bytes"
	"fmt"

	"github.com/russross/blackfriday"
)

type Parser struct {
	r blackfriday.Renderer
}

// NewParser returns a new parser with the inlineCallbacks required for org content
func NewParser(renderer blackfriday.Renderer) *Parser {
	return &Parser{r: renderer}
}

// OrgCommon is the easiest way to parse a byte slice of org content and makes assumptions
// that the caller wants to use blackfriday's HTMLRenderer with XHTML
func OrgCommon(input []byte) []byte {
	renderer := blackfriday.HtmlRenderer(blackfriday.HTML_USE_XHTML, "", "")
	return OrgOptions(input, renderer)
}

// Org is a convenience name for OrgOptions
func Org(input []byte, renderer blackfriday.Renderer) []byte {
	return OrgOptions(input, renderer)
}

// OrgOptions takes an org content byte slice and a renderer to use
func OrgOptions(input []byte, renderer blackfriday.Renderer) []byte {
	// in the case that we need to render something in isEmpty but there isn't a new line char
	input = append(input, '\n')
	var output bytes.Buffer

	p := NewParser(renderer)
	fmt.Println(p)

	return output.Bytes()
}

// Helpers
func skipChar(data []byte, start int, char byte) int {
	i := start
	for i < len(data) && charMatches(data[i], char) {
		i++
	}
	return i
}

func isSpace(char byte) bool {
	return charMatches(char, ' ')
}

func isEmpty(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	for i := 0; i < len(data) && !charMatches(data[i], '\n'); i++ {
		if !charMatches(data[i], ' ') && !charMatches(data[i], '\t') {
			return false
		}
	}
	return true
}

func charMatches(a byte, b byte) bool {
	return a == b
}
