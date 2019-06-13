package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

type parser struct {
	depth int
}

func (p *parser) peekToNewLine(items []lex.Item) (end int) {
	end = 0
	itemsLength := len(items)
	for end < itemsLength {
		currItem := items[end]
		if currItem.IsEOF() {
			return end
		}
		if currItem.IsNewline() {
			return end // we don't want to pass the newline character
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
		if prevIsNewline && isHeadline(items, end) {
			depth := headlineDepth(items[end:])
			if p.depth < depth {
				end++
				continue
			} else {
				return end - depth
			}
		}

		end++
	}
	return itemsLength
}

func findClosestSectionNode(parent ast.Node) ast.Node {
	for parent.Type() != "Root" {
		if parent.Type() == "Section" {
			return parent
		}
		parent = parent.Parent()
	}
	if len(parent.Children()) > 0 && parent.Children()[0].Type() == "Section" {
		return parent.Children()[0]
	}
	node := ast.NewSectionNode(parent)
	parent.Append(node)
	return node
}

func appendCurrentItemsToParent(start, current int, parent ast.Node, items []lex.Item) {
	if start < current {
		child := ast.NewTextNode(parent, items[start:current])
		parent.Append(child)
	}
}

func (p *parser) appendToParent(start, current int, parent ast.Node, items []lex.Item,
	findFunc func([]lex.Item) int,
	newNodeFunc func(int, int, ast.Node, []lex.Item) ast.Node) (newCurrent, newStart int) {

	end := current + findFunc(items[current:])
	node := newNodeFunc(current, end, parent, items)
	parent.Append(node)

	p.walkElements(node, items[current+1:end])

	current = end + 1
	start = current
	return current, start
}

func (p *parser) walkElements(parent ast.Node, items []lex.Item) {
	start := 0
	current := 0
	itemsLength := len(items)
	var prevText []lex.Item

	for current < itemsLength {
		if isLink(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newLink(parent, items[current:])
			current = end + 1
			start = current

		} else if isFootnoteDefinition(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newFootnoteDefinition(parent, items[current:])
			current = end + 1
			start = current

		} else if isBold(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newBold(parent, items[current:])
			current = end + 1
			start = current

		} else if isVerbatim(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newVerbatim(parent, items[current:])
			current = end + 1
			start = current

		} else if isItalic(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newItalic(parent, items[current:])
			current = end + 1
			start = current

		} else if isStrikeThrough(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newStrikeThrough(parent, items[current:])
			current = end + 1
			start = current

		} else if isUnderline(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newUnderline(parent, items[current:])
			current = end + 1
			start = current

		} else if isCode(items[current:]) {

			appendCurrentItemsToParent(start, current, parent, items)

			end := current + p.newCode(parent, items[current:])
			current = end + 1
			start = current

		} else {
			prevText = append(prevText, items[current])
			current++
		}
	}

	appendCurrentItemsToParent(start, itemsLength, parent, items)
}

// recursively walk through each token
func (p *parser) walk(parent ast.Node, items []lex.Item) {
	// create top-level paragraph nodes by creating nodes
	// on block level elements
	start := 0
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		token := items[current]

		if token.Type() == lex.ItemAsterisk && isHeadline(items, start) {
			depth := headlineDepth(items[current:])

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
			spaceWidth := 1

			peekStart := current + depth + spaceWidth

			headlineEnd := peekStart + p.peekToNewLine(items[peekStart:])

			node := ast.NewHeadlineNode(current, headlineEnd, depth, parent, items[peekStart:headlineEnd])

			end := headlineEnd + p.peekToNextBlock(items[headlineEnd:])
			parent.Append(node)

			// if headlineEnd = end, nothing left to parse in the headline!
			if headlineEnd != end {
				p.walk(node, items[headlineEnd:end])
				current = end
				start = current
			} else {
				current = end + 1
				start = current
			}
		} else if listTyp, listStart, listEnd, found := maybeList(items[current:]); found {

			p.makeList(listTyp, parent, items[start+listStart:start+listEnd])

			current = start + listEnd
			start = current

		} else if isTable(token, items, current) {
			tableEnd := findTable(items[current:])
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent)
			}
			node := ast.NewTableNode(current, tableEnd, parent, items[current:current+tableEnd])
			parent.Append(node)
			current = current + tableEnd
			start = current
		} else if isKeyword(token, items[current:]) {
			foundGreaterBlock, end := findGreaterBlock(items[current:])
			if foundGreaterBlock {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewGreaterBlockNode(current, end, parent, items[current:current+end])
				parent.Append(node)
				current = current + end
				continue
			}

			if parent.Type() == "Root" {
				node := ast.NewSectionNode(parent)
				parent.Append(node)
				parent = node
			}
			keywordWidth := 2 // #+
			keywordEnd := current + p.peekToNewLine(items[current:])
			current = current + keywordWidth
			node := ast.NewKeywordNode(start, current, parent, items[current:keywordEnd])
			parent.Append(node)
			current = keywordEnd
			start = current
		} else if isFootnoteDefinition(items[current:]) {
			if start < current {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewParagraphNode(start, current-1, parent, items[start:current-1])
				p.walkElements(node, items[start:current])
				parent.Append(node)
			}
			peekStart := current
			end := peekStart + p.peekToNewLine(items[peekStart:])

			node := ast.NewFootnoteDefinitionNode(parent, items[current:end])
			parent.Append(node)
			current = end
			start = current
		} else if token.IsEOF() || current+1 == itemsLength {
			if start < current {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewParagraphNode(start, current, parent, items[start:current])
				p.walkElements(node, items[start:current])
				parent.Append(node)
			}
			current++
			start = current
		} else if token.IsNewline() {
			if start < current && items[current-1].IsNewline() {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewParagraphNode(start, current, parent, items[start:current])
				p.walkElements(node, items[start:current])

				parent.Append(node)
				current++
				start = current
			}
			current++
		} else {
			current++
		}

	}
}

func findParentInTree(n ast.Node, typ ast.NodeType) bool {
	if n.Parent().Type() != "Root" && n.Parent().Type() != typ {
		findParentInTree(n.Parent(), typ)
	}
	return n.Parent().Type() == typ
}

// Parse takes an input string and returns a new Tree
func Parse(input string) *ast.RootNode {
	lexedItems := lex.NewLexer(input)
	p := &parser{depth: 0}
	// items is a slice that contains slices of items that are based on newline.
	var items []lex.Item

	root := ast.NewRootNode()

	// gather all of the tokens
	for item := range lexedItems {
		items = append(items, item)
	}

	p.walk(root, items)

	return root
}

func Append(children []ast.Node, child ast.Node) []ast.Node {
	return append(children, child)
}
