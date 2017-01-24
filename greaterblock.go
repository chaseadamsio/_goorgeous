package goorgeous

import (
	"bytes"
	"regexp"

	"github.com/russross/blackfriday"
)

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

var reOrderedList = regexp.MustCompile(`^\s*[0-9]+.\s+(.+)`)

func isOrderedList(data []byte) bool {
	return reOrderedList.Match(data)
}

var reUnorderedList = regexp.MustCompile(`^\s*-\s+(.+)`)

func isUnorderedList(data []byte) bool {
	return reUnorderedList.Match(data)
}
