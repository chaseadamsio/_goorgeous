package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

const maxHeadlineDepth = 6
const space = 1

func (p *parser) makeHeadline(parent ast.Node, start, end int) (blockEnd int) {

	depth := headlineDepth(p.items[start:end])

	if p.depth == 0 {
		p.depth = depth
	}

	blockEnd = p.peekToNextBlock(end)

	if p.depth < depth {
		// There will always be the chance that content occurs
		// before a headline, so there's a possiblity that we could be
		// in a Section. If that's the case, we traverse up to the Root
		// and reset so that we process the rest of the document properly
		if p.depth == 0 && parent.Type() != "Root" {
			for parent.Type() != "Root" {
				parent = parent.Parent()
			}
		}
		// set the new depth
		p.depth = depth
	}

	node := ast.NewHeadlineNode(depth, parent, p.items[start:end])

	p.walkHeadline(node, start+depth+space, end)

	parent.Append(node)

	// if foundEnd = end, nothing left to parse in the headline!
	if end != blockEnd {
		p.walk(node, end+1, blockEnd)
	}

	return blockEnd
}

func (p *parser) walkHeadline(node *ast.HeadlineNode, start, end int) {
	// var headlineText []string
	current := start

	if _, ok := keywords[p.items[current].Value()]; ok {
		node.Keyword = p.items[current].Value()
		current = current + 1 + space // account for keyword + space
	}

	if found, priority := hasPriority(p.items[current:end]); found {
		node.Priority = priority
		current = current + 4 + space
	}

	if found, foundStart, tags := p.findTags(current, end); found {
		node.Tags = tags
		end = foundStart
	}

	p.walkElements(node, current, end)
}

var keywords = map[string]struct{}{
	"TODO": struct{}{},
	"DONE": struct{}{},
}

func hasPriority(items []lex.Item) (found bool, priority string) {
	if len(items) < 4 {
		return false, ""
	}
	if !(items[0].IsBracket() && items[0].Value() == "[") {
		return false, ""
	}
	if !items[1].IsHash() {
		return false, ""
	}
	if !(items[2].IsText() && len(items[2].Value()) == 1) {
		return false, ""
	}
	if !(items[3].IsBracket() && items[3].Value() == "]") {
		return false, ""
	}
	return true, items[2].Value()
}

func (p *parser) findTags(current, end int) (hasTags bool, start int, tags []string) {

	var currTag []string

	for current < end {
		if p.items[current].IsColon() && !hasTags {

			hasTags = true
			start = current

		} else if p.items[current].IsColon() && hasTags {

			tags = append(tags, strings.Join(currTag, ""))
			currTag = make([]string, 1)

		} else if hasTags && p.items[current].IsText() ||
			p.items[current].IsUnderscore() || p.items[current].IsHash() {

			currTag = append(currTag, p.items[current].Value())

		} else if hasTags {

			hasTags = false

		}
		current++
	}

	return hasTags, start, tags
}

func hasKeyword(idx int, items []lex.Item) bool {
	// keywords will only _ever_ occur in the first space
	if idx != 0 {
		return false
	}
	if _, ok := keywords[items[idx].Value()]; ok {
		return true
	}
	return false
}

func (p parser) matchesHeadline(start int) (found bool, end int) {
	itemsLength := len(p.items)
	token := p.items[start]
	// first item has to be an asterisk
	if !token.IsAsterisk() {
		return false, -1
	}

	if 0 < start {
		reverseSearch := start - 1 // start with the previous character
		for 0 < reverseSearch {
			currItem := p.items[reverseSearch]

			if currItem.IsSpace() || currItem.IsTab() {
				reverseSearch--
				continue
			}
			if !currItem.IsNewline() {
				return false, -1
			}
			break
		}
	}

	current := start
	currHeadlineDepth := 0
	for currHeadlineDepth <= maxHeadlineDepth && current < itemsLength {
		currItem := p.items[current]
		// it's still a potential heading
		if currItem.IsAsterisk() {
			current++
			currHeadlineDepth++
			continue
		}
		// space terminates the headline "stars"
		if currItem.IsSpace() {
			depth := headlineDepth(p.items[current:])
			spaceWidth := 1
			peekStart := start + depth + spaceWidth
			headlineEnd := p.peekToNewLine(peekStart)
			return true, headlineEnd
		}
		return false, -1
	}
	return false, -1
}

// HeadlineDepth determines the depth of a headline
func headlineDepth(items []lex.Item) int {
	depth := 0
	itemsLength := len(items)
	for depth <= maxHeadlineDepth {
		hasNextItem := itemsLength > depth
		currItem := items[depth]
		if hasNextItem && currItem.IsAsterisk() {
			depth++
			continue
		}
		if hasNextItem && currItem.IsSpace() {
			return depth
		}
	}
	return depth
}
