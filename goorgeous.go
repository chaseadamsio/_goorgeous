package goorgeous

import (
	"bufio"
	"bytes"

	"github.com/russross/blackfriday"
)

type inlineParser func(p *parser, out *bytes.Buffer, data []byte, offset int) int

type parser struct {
	r              blackfriday.Renderer
	inlineCallback [256]inlineParser
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
	var output bytes.Buffer

	p := new(parser)
	p.r = renderer

	p.inlineCallback['='] = generateVerbatim
	p.inlineCallback['~'] = generateCode
	p.inlineCallback['/'] = generateEmphasis
	p.inlineCallback['*'] = generateBold
	p.inlineCallback['+'] = generateStrikethrough
	p.inlineCallback['['] = generateLink

	scanner := bufio.NewScanner(bytes.NewReader(input))
	// used to capture code blocks
	marker := ""
	syntax := ""
	listType := ""
	inList := false
	var tmpBlock bytes.Buffer

	for scanner.Scan() {
		data := scanner.Bytes()

		switch {
		case isEmpty(data):
			if inList == true {
				generateList := func() bool {
					output.Write(tmpBlock.Bytes())
					return true
				}
				switch listType {
				case "ul":
					p.r.List(&output, generateList, 0)
				case "ol":
					p.r.List(&output, generateList, blackfriday.LIST_TYPE_ORDERED)
				}
				inList = false
				listType = ""
				tmpBlock.Reset()
			}
			continue
		case isBlock(data) || marker != "":
			matches := reBlock.FindSubmatch(data)
			if len(matches) > 0 {
				if string(matches[1]) == "END" {
					p.r.BlockCode(&output, tmpBlock.Bytes(), syntax)
					marker = ""
					tmpBlock.Reset()
					continue
				}

			}
			if marker != "" {
				tmpBlock.Write(data)
				tmpBlock.WriteByte('\n')
			} else {
				marker = string(matches[2])
				syntax = string(matches[3])
			}
		case isComment(data):
			p.generateComment(&output, data)
		case isHeadline(data):
			p.generateHeadline(&output, data)
		case isDefinitionList(data):
			p.generateDefinitionList(&output, data)
		case isUnorderedList(data):
			if inList != true {
				listType = "ul"
				inList = true
			}
			matches := reUnorderedList.FindSubmatch(data)
			var work bytes.Buffer
			p.inline(&work, matches[1])
			p.r.ListItem(&tmpBlock, work.Bytes(), 0)
		case isOrderedList(data):
			if inList != true {
				listType = "ol"
				inList = true
			}
			matches := reOrderedList.FindSubmatch(data)
			var work bytes.Buffer
			p.inline(&work, matches[1])
			p.r.ListItem(&tmpBlock, work.Bytes(), 0)
		default:
			p.generateParagraph(&output, data)
		}
	}

	return output.Bytes()
}

func isEmpty(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	var i int
	for i = 0; i < len(data) && data[i] != '\n'; i++ {
		if data[i] != ' ' && data[i] != '\t' {
			return false
		}
	}
	return false
}

func isComment(data []byte) bool {
	return data[0] == '#' && data[1] == ' '
}

func (p *parser) generateComment(out *bytes.Buffer, data []byte) {
	var work bytes.Buffer
	work.WriteString("<!-- ")
	work.Write(data[2:])
	work.WriteString(" -->")
	work.WriteByte('\n')
	out.Write(work.Bytes())
}

func skipChar(data []byte, start int, char byte) int {
	i := start
	for i < len(data) && data[i] == char {
		i++
	}
	return i
}

func isSpace(char byte) bool {
	return char == ' '
}
