package goorgeous

import (
	"bytes"
	"regexp"
)

var reBlock = regexp.MustCompile(`^#\+(BEGIN|END)_(\w+)\s*([0-9A-Za-z_\-]*)?`)

func isBlock(data []byte) bool {
	return reBlock.Match(data)
}

func (p *parser) generateParagraph(out *bytes.Buffer, data []byte) {
	generate := func() bool {
		p.inline(out, bytes.Trim(data, " "))
		return true
	}
	p.r.Paragraph(out, generate)
}
