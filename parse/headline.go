package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

const maxHeadlineDepth = 6

func (p *parser) makeHeadline(parent ast.Node, items []lex.Item, start, end int) (node ast.Node) {

	depth := headlineDepth(items[start:end])

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

	node = ast.NewHeadlineNode(depth, parent, items[start:end])

	parent.Append(node)

	return node
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
			headlineEnd := peekStart + peekToNewLine(p.items[peekStart:])
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
