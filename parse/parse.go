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

func walkElements(parent ast.Node, items []lex.Item) {

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
				p.depth = depth
			}
			spaceWidth := 1

			peekStart := current + depth + spaceWidth

			headlineEnd := peekStart + p.peekToNewLine(items[peekStart:])

			// probably looked at this too long but there's a time when a headline is getting
			// parsed wrond so we still get "* " and have an empty headline. I think this is
			// probably because of the way we're accessing the slice in the previous call
			// (as an example, if you comment this out and run the headline - deep test you'll)
			// see the extra headline that I'm pretty sure is a result of off by a few errors in
			// the items sent to walk in the "in depth 2" section
			if len(items[peekStart:headlineEnd]) == 0 {
				current++
				continue
			}
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

		} else if token.Type() == lex.ItemNewLine {
			if start < current {
				// node := newSectionNode()
				node := ast.NewTextNode(start, current, parent, items[start:current])
				// ast.Append(textNode)
				parent.Append(node)
				start = current
			}
			current++
			start = current
		} else if token.Type() == lex.ItemEOF {
			if start < current {
				node := ast.NewTextNode(start, current, parent, items[start:current])
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
