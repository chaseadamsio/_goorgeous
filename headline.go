package goorgeous

import (
	"bytes"

	"github.com/shurcooL/sanitized_anchor_name"
)

func isHeadline(data []byte) bool {
	if data[0] != '*' {
		return false
	}
	level := 0
	for level < 6 && data[level] == '*' {
		level++
	}
	if data[level] != ' ' {
		return false
	}
	return true
}

func (p *parser) generateHeadline(out *bytes.Buffer, data []byte) {
	level := 1

	for level < 6 && data[level] == '*' {
		level++
	}

	start := level
	start = skipChar(data, start, ' ')

	data = data[start:]
	i := 0
	status := ""
	priority := ""

	if bytes.Contains(data[:4], []byte("TODO")) || bytes.Contains(data[:4], []byte("DONE")) {
		status = string(data[:4])
		i += 5 // one extra character for the next whitespace
	}

	if data[i] == '[' {
		priority = string(data[i+1])
		i += 4 // for "[c]" + ' '
	}

	tags := []string{}
	tagOpener := 0
	tagMarker := tagOpener
	for tIdx := i; tIdx < len(data); tIdx++ {
		if tagMarker > 0 && data[tIdx] == ':' {
			tags = append(tags, string(data[tagMarker+1:tIdx]))
			tagMarker = tIdx
		}
		if data[tIdx] == ':' && tagOpener == 0 {
			tagMarker = tIdx
			tagOpener = tIdx
		}
	}

	headlineID := sanitized_anchor_name.Create(string(data[i:]))

	generate := func() bool {
		dataEnd := len(data)
		if tagOpener > 0 {
			dataEnd = tagOpener
		}
		if status != "" {
			out.WriteString("<span class=\"todo " + status + "\">" + status + "</span>")
			out.WriteByte(' ')
		}
		if priority != "" {
			out.WriteString("<span class=\"priority " + priority + "\">[" + priority + "]</span>")
			out.WriteByte(' ')
		}
		out.Write(data[i:dataEnd])
		if tagOpener > 0 {
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
