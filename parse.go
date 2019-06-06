package goorgeous

type parser struct {
	depth int
}

func (p *parser) peekToNewLine(items []item) (end int) {
	end = 0
	itemsLength := len(items)
	for end < itemsLength {
		if items[end].typ == itemEOF {
			return end
		}
		if items[end].typ == itemNewLine {
			return end
		}
		end++
	}
	return itemsLength
}

func (p *parser) peekToNextBlock(items []item) (end int) {
	end = 0
	itemsLength := len(items)
	for end < itemsLength {
		if items[end].typ == itemEOF {
			return end
		}
		if end > 0 && items[end-1].typ == itemNewLine && isHeadline(items[end:], items[end]) {
			depth := headlineDepth(items[end:])
			if p.depth < depth {
				end++
				continue
			} else {
				return end - 2
			}
		}

		end++
	}
	return itemsLength
}

func walkElements(parent Node, items []item) {

}

// recursively walk through each token
func (p *parser) walk(parent Node, items []item) {
	// create top-level paragraph nodes by creating nodes
	// on block level elements
	start := 0
	current := 0
	itemsLength := len(items)
	for current < itemsLength {
		token := items[current]

		if token.typ == itemAsterisk && isHeadline(items[current:], token) {
			depth := headlineDepth(items[current:])
			if p.depth < depth {
				p.depth = depth
			}

			if p.depth > depth {
				for depth <= p.depth {
					parent = parent.Parent()
					p.depth--
				}
			}
			spaceWidth := 1

			peekStart := current + depth + spaceWidth

			headlineEnd := p.peekToNewLine(items)

			// probably looked at this too long but there's a time when a headline is getting
			// parsed wrond so we still get "* " and have an empty headline. I think this is
			// probably because of the way we're accessing the slice in the previous call
			// (as an example, if you comment this out and run the headline - deep test you'll)
			// see the extra headline that I'm pretty sure is a result of off by a few errors in
			// the items sent to walk in the "in depth 2" section
			if len(items[peekStart:peekStart+headlineEnd]) == 0 {
				current++
				continue
			}
			node := newHeadlineNode(current, peekStart+headlineEnd, depth, parent, items[peekStart:peekStart+headlineEnd])

			end := p.peekToNextBlock(items[peekStart+headlineEnd:])
			parent.Append(node)

			afterHeadlineNewLine := peekStart + headlineEnd

			if afterHeadlineNewLine < peekStart+end {
				p.walk(node, items[afterHeadlineNewLine:peekStart+end])
				current = peekStart + end
				start = current
			} else {
				current = afterHeadlineNewLine
				start = current
			}

		} else if token.typ == itemNewLine {
			if start < current {
				// node := newSectionNode()
				node := newTextNode(start, current, parent, items)
				// node.Append(textNode)
				parent.Append(node)
				start = current
			}
			current++
			start = current
		} else if token.typ == itemEOF {
			if start < current {
				node := newTextNode(start, current, parent, items)
				parent.Append(node)
			}
			return
			// current++
			// start = current
		} else {
			current++
		}

	}
}

func findParentInTree(n Node, typ NodeType) bool {
	if n.Parent().Type() != "Root" && n.Parent().Type() != typ {
		findParentInTree(n.Parent(), typ)
	}
	return n.Parent().Type() == typ
}

// Parse takes an input string and returns a new Tree
func Parse(input string) *RootNode {
	l := lex(input)
	p := &parser{depth: 0}
	// items is a slice that contains slices of items that are based on newline.
	var items []item

	root := newRootNode()

	// gather all of the tokens
	for item := range l.items {
		items = append(items, item)
	}

	p.walk(root, items)

	return root
}

func Append(children []Node, child Node) []Node {
	return append(children, child)
}
