package goorgeous

import (
	"bufio"
	"bytes"
	"regexp"

	"github.com/russross/blackfriday"
	"github.com/shurcooL/sanitized_anchor_name"
)

type inlineParser func(p *parser, out *bytes.Buffer, data []byte, offset int) int

type parser struct {
	r              blackfriday.Renderer
	inlineCallback [256]inlineParser
}

// NewParser returns a new parser with the inlineCallbacks required for org content
func NewParser(renderer blackfriday.Renderer) *parser {
	p := new(parser)
	p.r = renderer

	p.inlineCallback['='] = generateVerbatim
	p.inlineCallback['~'] = generateCode
	p.inlineCallback['/'] = generateEmphasis
	p.inlineCallback['*'] = generateBold
	p.inlineCallback['+'] = generateStrikethrough
	p.inlineCallback['['] = generateLinkOrImg

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
			switch {
			case inList:
				p.generateList(&output, tmpBlock.Bytes(), listType)
				inList = false
				listType = ""
				tmpBlock.Reset()
			case inTable:
				p.generateTable(&output, tmpBlock.Bytes())
				inTable = false
				listType = ""
				tmpBlock.Reset()
			default:
				continue
			}
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
		case isKeyword(data):
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
		case isHorizontalRule(data):
			p.r.HRule(&output)
		default:
			p.generateParagraph(&output, data)
		}
	}

	return output.Bytes()
}

// Org Syntax has been broken up into 4 distinct sections based on
// the org-syntax draft (http://orgmode.org/worg/dev/org-syntax.html):
// - Headlines
// - Greater Elements
// - Elements
// - Objects

// Headlines
func isHeadline(data []byte) bool {
	if !charMatches(data[0], '*') {
		return false
	}
	level := 0
	for level < 6 && charMatches(data[level], '*') {
		level++
	}
	return charMatches(data[level], ' ')
}

func (p *parser) generateHeadline(out *bytes.Buffer, data []byte) {
	level := 1
	status := ""
	priority := ""

	for level < 6 && data[level] == '*' {
		level++
	}

	start := skipChar(data, level, ' ')

	data = data[start:]
	i := 0

	// Check if has a status so it can be rendered as a separate span that can be hidden or
	// modified with CSS classes
	if hasStatus(data[i:4]) {
		status = string(data[i:4])
		i += 5 // one extra character for the next whitespace
	}

	// Check if the next byte is a priority marker
	if data[i] == '[' && hasPriority(data[i+1]) {
		priority = string(data[i+1])
		i += 4 // for "[c]" + ' '
	}

	tags, tagsFound := findTags(data, i)

	headlineID := sanitized_anchor_name.Create(string(data[i:]))

	generate := func() bool {
		dataEnd := len(data)
		if tagsFound > 0 {
			dataEnd = tagsFound
		}

		headline := bytes.TrimRight(data[i:dataEnd], " \t")

		if status != "" {
			out.WriteString("<span class=\"todo " + status + "\">" + status + "</span>")
			out.WriteByte(' ')
		}

		if priority != "" {
			out.WriteString("<span class=\"priority " + priority + "\">[" + priority + "]</span>")
			out.WriteByte(' ')
		}

		out.Write(headline)

		if tagsFound > 0 {
			for _, tag := range tags {
				out.WriteByte(' ')
				out.WriteString("<span class=\"tags\">" + tag + "</span>")
				out.WriteByte(' ')
			}
		}
		return true
	}

	p.r.Header(out, generate, level, headlineID)
}

func hasStatus(data []byte) bool {
	return bytes.Contains(data, []byte("TODO")) || bytes.Contains(data, []byte("DONE"))
}

func hasPriority(char byte) bool {
	return (charMatches(char, 'A') || charMatches(char, 'B') || charMatches(char, 'C'))
}

func findTags(data []byte, start int) ([]string, int) {
	tags := []string{}
	tagOpener := 0
	tagMarker := tagOpener
	for tIdx := start; tIdx < len(data); tIdx++ {
		if tagMarker > 0 && data[tIdx] == ':' {
			tags = append(tags, string(data[tagMarker+1:tIdx]))
			tagMarker = tIdx
		}
		if data[tIdx] == ':' && tagOpener == 0 {
			tagMarker = tIdx
			tagOpener = tIdx
		}
	}
	return tags, tagOpener
}

// Greater Elements
// ~~ Definition Lists
var reDefinitionList = regexp.MustCompile(`^\s*-\s+(.+?)\s+::\s+(.*)`)

func isDefinitionList(data []byte) bool {
	return reDefinitionList.Match(data)
}

func (p *parser) generateDefinitionList(out *bytes.Buffer, data []byte) {
	flags := blackfriday.LIST_TYPE_DEFINITION
	flags |= blackfriday.LIST_ITEM_BEGINNING_OF_LIST
	matches := reDefinitionList.FindSubmatch(data)
	generate := func() bool {
		flags |= blackfriday.LIST_TYPE_TERM
		p.r.ListItem(out, matches[1], flags)
		flags &= ^blackfriday.LIST_TYPE_TERM
		p.r.ListItem(out, matches[2], flags)
		flags &= ^blackfriday.LIST_ITEM_END_OF_LIST
		return true
	}
	p.r.List(out, generate, flags)
}

// ~~ Ordered Lists
var reOrderedList = regexp.MustCompile(`^\s*[0-9]+.\s+(.+)`)

func isOrderedList(data []byte) bool {
	return reOrderedList.Match(data)
}

// ~~ Unordered Lists
var reUnorderedList = regexp.MustCompile(`^\s*-\s+(.+)`)

func isUnorderedList(data []byte) bool {
	return reUnorderedList.Match(data)
}

// ~~ Tables
var reTableHeaders = regexp.MustCompile(`^[|+-]*$`)

func isTable(data []byte) bool {
	return charMatches(data[0], '|')
}

func (p *parser) generateTable(output *bytes.Buffer, data []byte) {
	var table bytes.Buffer
	rows := bytes.Split(bytes.Trim(data, "\n"), []byte("\n"))
	hasTableHeaders := reTableHeaders.Match(rows[1])
	tbodySet := false

	for idx, row := range rows {
		var rowBuff bytes.Buffer
		if hasTableHeaders && idx == 0 {
			table.WriteString("<thead>")
			for _, cell := range bytes.Split(row[1:len(row)-1], []byte("|")) {
				p.r.TableHeaderCell(&rowBuff, bytes.Trim(cell, " \t"), 0)
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
				p.inline(&cellBuff, bytes.Trim(cell, " \t"))
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

// ~~ Dynamic Blocks
var reBlock = regexp.MustCompile(`^#\+(BEGIN|END)_(\w+)\s*([0-9A-Za-z_\-]*)?`)

func isBlock(data []byte) bool {
	return reBlock.Match(data)
}

// Elements
// ~~ Keywords
func isKeyword(data []byte) bool {
	return charMatches(data[0], '#') && charMatches(data[1], '+') && !charMatches(data[2], ' ')
}

// ~~ Comments
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

// ~~ Horizontal Rules
var reHorizontalRule = regexp.MustCompile(`^\s*?-----\s?$`)

func isHorizontalRule(data []byte) bool {
	return reHorizontalRule.Match(data)
}

// ~~ Paragraphs
func (p *parser) generateParagraph(out *bytes.Buffer, data []byte) {
	generate := func() bool {
		p.inline(out, bytes.Trim(data, " "))
		return true
	}
	p.r.Paragraph(out, generate)
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

// Objects

func (p *parser) inline(out *bytes.Buffer, data []byte) {
	i, end := 0, 0

	for i < len(data) {
		for end < len(data) && p.inlineCallback[data[end]] == nil {
			end++
		}

		p.r.NormalText(out, data[i:end])

		if end >= len(data) {
			break
		}
		i = end

		handler := p.inlineCallback[data[i]]

		if consumed := handler(p, out, data, i); consumed > 0 {
			i += consumed
			end = i
			continue
		}

		end = i + 1
	}
}

func generator(p *parser, out *bytes.Buffer, data []byte, offset int, char byte, doInline bool, renderer func(*bytes.Buffer, []byte)) int {
	data = data[offset:]
	c := byte(char)
	start := 1
	i := start
	if len(data) <= 1 {
		return 0
	}

	// Org mode spec says a non-whitespace character must immediately follow.
	// if the current char is the marker, then there's no text between, not a candidate
	if isSpace(data[i]) || charMatches(data[i], c) {
		return 0
	}

	for i < len(data) {
		if charMatches(data[i], c) {
			if c == '/' {
				if len(data) > i+1 && charMatches(data[i+1], '/') {
					return 0
				}
			}
			var work bytes.Buffer
			if doInline {
				p.inline(&work, data[start:i])
				renderer(out, work.Bytes())
			} else {
				renderer(out, data[start:i])
			}
			return i + 1
		}
		i++
	}
	return 0
}

// ~~ Images and Links
var reLinkOrImg = regexp.MustCompile(`\[\[(.+?)\]\[?(.*?)\]?\]`)

func generateLinkOrImg(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	data = data[offset+1:]
	start := 1
	i := start
	var hyperlink []byte
	isImage := false
	closedLink := false
	hasContent := false

	if data[0] != '[' {
		return 0
	}

	if bytes.Equal(data[1:6], []byte("file:")) {
		isImage = true
	}

	for i < len(data) {
		currChar := data[i]
		switch {
		case charMatches(currChar, ']') && closedLink == false:
			if isImage {
				hyperlink = data[start+5 : i]
			} else {
				hyperlink = data[start:i]
			}
			closedLink = true
		case charMatches(currChar, '['):
			start = i + 1
			hasContent = true
		case charMatches(currChar, ']') && closedLink == true && hasContent == true && isImage == true:
			p.r.Image(out, hyperlink, data[start:i], data[start:i])
			return i + 3
		case charMatches(currChar, ']') && closedLink == true && hasContent == true:
			p.r.Link(out, hyperlink, data[start:i], data[start:i])
			return i + 3
		case charMatches(currChar, ']') && closedLink == true && hasContent == false && isImage == true:
			p.r.Image(out, hyperlink, hyperlink, hyperlink)
			return i + 2
		case charMatches(currChar, ']') && closedLink == true && hasContent == false:
			p.r.Link(out, hyperlink, hyperlink, hyperlink)
			return i + 2
		}
		i++
	}

	return 0
}

// ~~ Text Markup

func generateVerbatim(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return generator(p, out, data, offset, '=', false, p.r.CodeSpan)
}

func generateCode(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return generator(p, out, data, offset, '~', false, p.r.CodeSpan)
}

func generateEmphasis(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return generator(p, out, data, offset, '/', true, p.r.Emphasis)
}

func generateBold(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return generator(p, out, data, offset, '*', true, p.r.DoubleEmphasis)
}

func generateStrikethrough(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return generator(p, out, data, offset, '+', true, p.r.StrikeThrough)
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