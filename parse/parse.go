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
	listMatchers    []foundMatchers
}

func newParser(input string) *parser {
	p := &parser{
		input: input,
		depth: 0,
	}

	p.elementMatchers = []foundMatchers{
		{p.matchesLink, p.newLink},
		{p.matchesFootnoteReference, p.makeFootnoteReference},
		{p.matchesBold, p.newBold},
		{p.matchesItalic, p.newItalic},
		{p.matchesVerbatim, p.newVerbatim},
		{p.matchesStrikeThrough, p.newStrikeThrough},
		{p.matchesUnderline, p.newUnderline},
		{p.matchesCode, p.newCode},
		{p.matchesEnDash, p.newEnDash},
		{p.matchesMDash, p.newMDash},
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

	root := ast.NewRootNode(0, p.items[len(p.items)-1].End())

	p.walk(root, 0, len(p.items))

	return root
}

func (p *parser) peekToNewLine(start int) (end int) {
	current := start
	itemsLength := len(p.items)
	end = itemsLength
	for current < end {
		currItem := p.items[current]
		if currItem.IsEOF() {
			return current
		}
		if currItem.IsNewline() {
			return current // we don't want to pass the newline character
		}
		current++
	}
	return current
}

func (p *parser) peekToNextBlock(current int) int {
	itemsLength := len(p.items)
	for current < itemsLength {
		prevIsNewline := current > 0 && p.items[current-1].IsNewline()
		if found, foundEnd := p.matchesHeadline(current); prevIsNewline && found {
			depth := headlineDepth(p.items[current:foundEnd])
			if p.depth < depth {
				current++
				continue
			} else {
				return current - 1
			}
		} else if p.items[current].IsNewline() && current+1 < itemsLength && p.items[current+1].IsEOF() {
			return current

		}
		current++
	}
	return current
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

func (p *parser) appendCurrentItemsToParent(parent ast.Node, start, end int) {
	if start != end {
		child := ast.NewTextNode(parent, p.items[start:end])
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
				p.appendCurrentItemsToParent(parent, start, current)

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

	p.appendCurrentItemsToParent(parent, start, current)
}

// recursively walk through each token
func (p *parser) walk(parent ast.Node, current, stop int) {
	// create top-level paragraph nodes by creating nodes
	// on block level elements
	start := current
	for current < stop {
		token := p.items[current]

		if found, end := p.matchesHeadline(current); found {

			blockEnd := p.makeHeadline(parent, current, end)
			current = blockEnd
			start = current

		} else if found, end := p.matchesOrderedList(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeOrderedList(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesUnorderedList(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeUnorderedList(parent, current, end)
			current = end
			start = current
		} else if found, end := p.matchesFootnoteDefinition(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeFootnoteDefinition(parent, current, end)
			current = end
			start = current
		} else if found, end := p.matchesHorizontalRule(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeHorizontalRule(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesTable(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeTable(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesGreaterBlock(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeGreaterBlock(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesKeyword(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeKeyword(parent, current, end)
			current = end
			start = current

		} else if found, end := p.matchesFootnoteDefinition(current); found {
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent, p.items[current:])
			}
			p.makeFootnoteDefinition(parent, current, end)
			current = end
			start = current

		} else if token.IsNewline() {

			if current+1 < stop && (p.items[current+1].IsNewline() || p.items[current+1].IsEOF()) {
				if start == current {
					current++
					start = current
					continue
				}

				if p.items[start].IsNewline() && p.items[start+1].IsEOF() {
					break
				}

				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent, p.items[current:])
				}
				node := ast.NewParagraphNode(start, current, parent, p.items[start:current])
				p.walkElements(node, start, current)

				parent.Append(node)
				current++
				start = current
				continue
			}

			if 0 < current && p.items[current-1].IsNewline() {
				current++
				start = current
				continue
			}

			current++

		} else if token.IsEOF() || current+1 == stop {
			if p.items[start].IsNewline() && p.items[start+1].IsEOF() {
				break
			}

			if current+1 == stop && stop != len(p.items) {
				current++ // we've reached the end of an inner element's parse
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
