package goorgeous

import (
	"bytes"
	"regexp"
)

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

		handler := p.inlineCallback[data[end]]

		if consumed := handler(p, out, data, i); consumed == 0 {
			end = i + 1
		} else {
			i += consumed
			end = i
		}
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
	if isSpace(data[i]) || data[i] == c {
		return 0
	}

	for i < len(data) {
		if data[i] == c {
			if c == '/' {
				if len(data) > i+1 && data[i+1] == '/' {
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

var reLink = regexp.MustCompile(`\[\[(.+?)\]\[?(.*?)\]?\]`)

func generateLink(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	data = data[offset+1:]
	start := 1
	i := start
	var hyperlink []byte
	closedLink := false
	hasContent := false

	if data[0] != '[' {
		return 0
	}

	for i < len(data) {
		if data[i] == ']' && closedLink == false {
			hyperlink = data[start:i]
			closedLink = true
		} else if data[i] == '[' {
			start = i + 1
			hasContent = true
		} else if data[i] == ']' && closedLink == true && hasContent == true {
			p.r.Link(out, hyperlink, data[start:i], data[start:i])
			return i + 3
		} else if data[i] == ']' && closedLink == true && hasContent == false {
			p.r.Link(out, hyperlink, hyperlink, hyperlink)
			return i + 2
		}
		i++
	}

	return 0
}
