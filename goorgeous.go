package goorgeous

import (
	"bufio"
	"bytes"
	"regexp"

	"github.com/russross/blackfriday"
)

type inlineParser func(p *parser, out *bytes.Buffer, data []byte, offset int) int

type parser struct {
	r              blackfriday.Renderer
	inlineCallback [256]inlineParser
}

func NewParser(renderer blackfriday.Renderer) *parser {
	p := new(parser)
	p.r = renderer

	p.inlineCallback['='] = generateVerbatim
	p.inlineCallback['~'] = generateCode
	p.inlineCallback['/'] = generateEmphasis
	p.inlineCallback['*'] = generateBold
	p.inlineCallback['+'] = generateStrikethrough
	p.inlineCallback['['] = generateLink

	return p
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

	scanner := bufio.NewScanner(bytes.NewReader(input))
	// used to capture code blocks
	marker := ""
	syntax := ""
	listType := ""
	inList := false
	inTable := false
	var tmpBlock bytes.Buffer

	for scanner.Scan() {
		data := scanner.Bytes()

		switch {
		case isEmpty(data):
			if inList == true {
				p.generateList(&output, tmpBlock.Bytes(), listType)
				inList = false
				listType = ""
				tmpBlock.Reset()
			}
			if inTable == true {
				p.generateTable(&output, tmpBlock.Bytes())

				inTable = false
				listType = ""
				tmpBlock.Reset()
			}
			continue
		case isBlock(data) || marker != "":
			matches := reBlock.FindSubmatch(data)
			if len(matches) > 0 {
				if string(matches[1]) == "END" {
					switch marker {
					case "SRC":
						p.r.BlockCode(&output, tmpBlock.Bytes(), syntax)
					case "QUOTE":
						var tmpBuf bytes.Buffer
						p.inline(&tmpBuf, tmpBlock.Bytes())
						p.r.BlockQuote(&output, tmpBuf.Bytes())
					case "CENTER":
						var tmpBuf bytes.Buffer
						output.WriteString("<center>\n")
						p.inline(&tmpBuf, tmpBlock.Bytes())
						output.Write(tmpBuf.Bytes())
						output.WriteString("</center>\n")
					}
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
		case isTable(data):
			if inTable != true {
				inTable = true
			}
			tmpBlock.Write(data)
			tmpBlock.WriteByte('\n')
		case isHeader(data):
			continue
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

func isTable(data []byte) bool {
	return charMatches(data[0], '|')
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
	return true
}

func charMatches(a byte, b byte) bool {
	return a == b
}

func isHeader(data []byte) bool {
	return charMatches(data[0], '#') && charMatches(data[1], '+') && !charMatches(data[2], ' ')
}

func isComment(data []byte) bool {
	return charMatches(data[0], '#') && charMatches(data[1], ' ')
}

func (p *parser) generateComment(out *bytes.Buffer, data []byte) {
	var work bytes.Buffer
	work.WriteString("<!-- ")
	work.Write(data[2:])
	work.WriteString(" -->")
	work.WriteByte('\n')
	out.Write(work.Bytes())
}

func (p *parser) generateList(output *bytes.Buffer, data []byte, listType string) {
	generateList := func() bool {
		output.Write(data)
		return true
	}
	switch listType {
	case "ul":
		p.r.List(output, generateList, 0)
	case "ol":
		p.r.List(output, generateList, blackfriday.LIST_TYPE_ORDERED)
	}
}

func (p *parser) generateTable(output *bytes.Buffer, data []byte) {
	var table bytes.Buffer
	rows := bytes.Split(bytes.Trim(data, "\n"), []byte("\n"))
	reTableHeaders := regexp.MustCompile(`^[|+-]*$`)
	hasTableHeaders := reTableHeaders.Match(rows[1])
	tbodySet := false

	for idx, row := range rows {
		var rowBuff bytes.Buffer
		if hasTableHeaders && idx == 0 {
			table.WriteString("<thead>")
			for _, cell := range bytes.Split(row[1:len(row)-1], []byte("|")) {
				p.r.TableHeaderCell(&rowBuff, cell, 0)
			}
			p.r.TableRow(&table, rowBuff.Bytes())
			table.WriteString("</thead>\n")
		} else if hasTableHeaders && idx == 1 {
			continue
		} else {
			if !tbodySet {
				table.WriteString("<tbody>")
				tbodySet = true
			}
			for _, cell := range bytes.Split(row[1:len(row)-1], []byte("|")) {
				var cellBuff bytes.Buffer
				p.inline(&cellBuff, cell)
				p.r.TableCell(&rowBuff, cellBuff.Bytes(), 0)
			}
			p.r.TableRow(&table, rowBuff.Bytes())
			if tbodySet && idx == len(rows)-1 {
				table.WriteString("</tbody>\n")
				tbodySet = false
			}

		}
	}

	output.WriteString("\n<table>\n")
	output.Write(table.Bytes())
	output.WriteString("</table>\n")
}
