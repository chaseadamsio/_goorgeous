package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

type matcherFunc func(current int) (found bool, end int)
type createElementFunc func(parent ast.Node, start, end int)

type foundMatchers struct {
	matcher       matcherFunc
	createElement createElementFunc
}

type parser struct {
	items           []lex.Item
	input           string
	depth           int
	elementMatchers []foundMatchers
	blockMatchers   []foundMatchers
	listMatchers    []foundMatchers
}

func newParser(input string) *parser {
	p := &parser{
		input: input,
		depth: 0,
	}

	p.elementMatchers = []foundMatchers{
		{p.matchesLink, p.newLink},
		{p.matchesFootnoteDefinition, p.makeFootnoteDefinition},
		{p.matchesBold, p.newBold},
		{p.matchesItalic, p.newItalic},
		{p.matchesVerbatim, p.newVerbatim},
		{p.matchesStrikeThrough, p.newStrikeThrough},
		{p.matchesUnderline, p.newUnderline},
		{p.matchesCode, p.newCode},
	}

	p.blockMatchers = []foundMatchers{
		{p.matchesGreaterBlock, p.makeGreaterBlock},
		{p.matchesFootnoteDefinition, p.makeFootnoteDefinition},
		{p.matchesTable, p.makeTable},
	}

	return p
}

// Parse takes an input string and returns a new Tree
func Parse(input string) *ast.RootNode {
	lexedItems := lex.NewLexer(input)

	p := newParser(input)

	// gather all of the tokens
	for item := range lexedItems {
		p.items = append(p.items, item)
	}

	root := ast.NewRootNode(p.items)

	p.walk(root, 0, len(p.items))

	return root
}

func peekToNewLine(items []lex.Item) (end int) {
	end = 0
	itemsLength := len(items)
	for end < itemsLength {
		currItem := items[end]
		if currItem.IsEOF() {
			return end
		}
		if currItem.IsNewline() {
			if end+1 < itemsLength && items[end+1].IsEOF() {
				return end + 2
			}
			return end + 1 // we don't want to pass the newline character
		}
		end++
	}
	return itemsLength
}

func (p *parser) peekToNextBlock(items []lex.Item) (end int) {
	end = 0
	itemsLength := len(items)
	for end < itemsLength {
		currItem := items[end]
		prevIsNewline := end > 0 && items[end-1].IsNewline()
		if currItem.IsEOF() {
			return end
		}
		if found, foundEnd := p.matchesHeadline(end); prevIsNewline && found {
			depth := headlineDepth(items[end:foundEnd])
			if p.depth < depth {
				end++
				continue
			} else {
				return end
			}
		} else if items[end].IsNewline() && end+1 < itemsLength && items[end+1].IsEOF() {
			return end + 1
		}

		end++
	}
	return itemsLength
}

func findClosestSectionNode(parent ast.Node, items []lex.Item) ast.Node {
	for parent.Type() != "Root" {
		if parent.Type() == "Section" {
			return parent
		}
		parent = parent.Parent()
	}
	if len(parent.Children()) > 0 && parent.Children()[0].Type() == "Section" {
		return parent.Children()[0]
	}
	node := ast.NewSectionNode(parent, items)
	parent.Append(node)
	return node
}

func appendCurrentItemsToParent(parent ast.Node, items []lex.Item) {
	if 0 < len(items) {
		child := ast.NewTextNode(parent, items)
		parent.Append(child)
	}
}

func (p *parser) walkElements(parent ast.Node, current, end int) {
	start := current

	for current < end {
		foundMatch := false
	MatcherLoop:
		for _, f := range p.elementMatchers {
			if found, end := f.matcher(current); found {
				appendCurrentItemsToParent(parent, p.items[start:current])

				f.createElement(parent, current, end)
				current = end
				start = current
				foundMatch = true
				break MatcherLoop
			}
		}
		if !foundMatch {
			current++
		}
	}

	appendCurrentItemsToParent(parent, p.items[start:current])
}

// recursively walk through each token
func (p *parser) walk(parent ast.Node, current, end int) {
	// create top-level paragraph nodes by creating nodes
	// on block level elements
	start := current
	for current < end {
		token := p.items[current]

		if found, end := p.matchesHeadline(current); found {
			if p.depth == 0 {
				p.depth = headlineDepth(p.items[current:end])
			}

			blockEnd := end + p.peekToNextBlock(p.items[end:])
			node := p.makeHeadline(parent, p.items, current, blockEnd)
			current = blockEnd
			start = current

			// if foundEnd = end, nothing left to parse in the headline!
			if end != blockEnd {
				p.walk(node, end, blockEnd)
			}

		} else if found, end := p.matchesOrderedList(current); found {

			p.makeList("ORDERED", parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesUnorderedList(current); found {

			p.makeList("UNORDERED", parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesHorizontalRule(current); found {
			p.makeHorizontalRule(parent, p.items[current:end])
			current = end
			start = current

		} else if found, end := p.matchesTable(current); found {

			p.makeTable(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesGreaterBlock(current); found {

			p.makeGreaterBlock(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesKeyword(current); found {

			p.makeKeyword(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesFootnoteDefinition(current); found {

			p.makeFootnoteDefinition(parent, current, end)
			current = end
			start = current

		} else if token.IsNewline() {
			if start < current && p.items[current-1].IsNewline() {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent, p.items[current:])
				}
				node := ast.NewParagraphNode(start, current, parent, p.items[start:current])
				p.walkElements(node, start, current)

				parent.Append(node)
				current++
				start = current
			}

			current++
		} else if token.IsEOF() || current+1 == end {
			if current-start == 1 && p.items[current-1].IsNewline() {
				break
			}
			if start < current {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent, p.items[current:])
				}
				node := ast.NewParagraphNode(start, current, parent, p.items[start:current])
				p.walkElements(node, start, current)
				parent.Append(node)
			}
			current++
			start = current
		} else {
			current++
		}

	}
}
