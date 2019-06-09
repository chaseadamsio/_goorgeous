package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
	"github.com/chaseadamsio/goorgeous/tokens"
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
		if prevIsNewline && tokens.IsHeadline(items[end:]) {
			depth := tokens.HeadlineDepth(items[end:])
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

// recursively walk through each token
func (p *parser) walk(parent ast.Node, items []lex.Item) {
	// create top-level paragraph nodes by creating nodes
	// on block level elements
	start := 0
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		token := items[current]

		if token.Type() == lex.ItemAsterisk && tokens.IsHeadline(items[current:]) {
			depth := tokens.HeadlineDepth(items[current:])

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

		} else if token.IsNewline() {
			if start < current && items[current-1].IsNewline() {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewParagraphNode(start, current, parent, items[start:current])
				parent.Append(node)
				current++
				start = current
			}
			current++
		} else if tokens.IsOrderedList(token, items, current) {
			// orderedListEnd := tokens.FindOrderedList(items[current:])
			// node := ast.NewOrderedListNode(current, orderedListEnd, parent, items[current:current+orderedListEnd])
			// parent.Append(node)
			// current = current + orderedListEnd
			current++
		} else if tokens.IsUnorderedList(token, items, current) {
			unorderedListEnd := tokens.FindUnorderedList(items[current:])
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent)
			}
			node := ast.NewUnorderedListNode(current, unorderedListEnd, parent, items[current:current+unorderedListEnd])
			parent.Append(node)
			current = current + unorderedListEnd
			start = current
		} else if tokens.IsTable(token, items, current) {
			tableEnd := tokens.FindTable(items[current:])
			if parent.Type() != "Section" {
				parent = findClosestSectionNode(parent)
			}
			node := ast.NewTableNode(current, tableEnd, parent, items[current:current+tableEnd])
			parent.Append(node)
			current = current + tableEnd
			start = current
		} else if tokens.IsKeyword(token, items[current:]) {
			foundGreaterBlock, end := tokens.FindGreaterBlock(items[current:])
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
		} else if tokens.IsFootnoteDefinition(token, items, current) {
			if start < current {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewParagraphNode(start, current-1, parent, items[start:current-1])
				parent.Append(node)
			}
			peekStart := current
			end := peekStart + p.peekToNewLine(items[peekStart:])

			node := ast.NewFootnoteDefinitionNode(start, current, parent, items[current:end])
			parent.Append(node)
			current = end
			start = current
		} else if token.IsEOF() || current+1 == itemsLength {
			if start < current {
				if parent.Type() != "Section" {
					parent = findClosestSectionNode(parent)
				}
				node := ast.NewParagraphNode(start, current, parent, items[start:current])
				parent.Append(node)
			}
			current++
			start = current
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
