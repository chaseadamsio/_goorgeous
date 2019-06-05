package goorgeous

// Tree is the general shape of a node
type Tree struct {
	typ        NodeType
	items      []item
	attributes map[string]interface{}
	children   []Tree
	currPos    int
	start      Pos
	end        Pos
}

// recursively walk through each token and return a node
func walk() {

}

// Parse takes an input string and returns a new Tree
func Parse(input string) *Tree {
	l := lex(input)

	root := &Tree{
		typ: NodeRoot,
		start: Pos{
			Column: 1, Offset: 0, Line: 1,
		},
	}

	for item := range l.items {
		root.items = append(root.items, item)
	}

	root.end = Pos{
		Column: root.items[len(root.items)-1].Column,
		Offset: root.items[len(root.items)-1].Offset,
		Line:   root.items[len(root.items)-1].Line,
	}

	// create top-level paragraph nodes by creating nodes
	// on block level elements
	start := 0
	for idx, token := range root.items {
		root.currPos = idx
		if token.typ == itemNewLine {
			// previous token was a newlne, don't create a new block
			if len(root.items) > idx-1 && root.items[idx-1].typ == itemNewLine {
				start = idx
				continue
			}

			child := processNode(root, start)

			child.parse()
			root.children = append(root.children, child)
			start = root.currPos
		}
	}

	if root.currPos > start {
		child := processNode(root, start)
		child.parse()
		root.children = append(root.children, child)
	}

	return root
}

func processNode(t *Tree, start int) Tree {
	if isHeadline(t.items) {
		t.headline(start)
	}
	return Tree{
		typ:   NodeParagraph,
		items: t.items[start : t.currPos+1],
	}

}

func (t *Tree) headline(start int) Tree {
	depth := headlineDepth(t.items)
	itemStart := start + depth + 1 // current + depth + 1 for trailing space
	return Tree{
		typ: NodeHeadline,
		attributes: map[string]interface{}{
			"depth": depth,
		},
		items: t.items[itemStart : t.currPos-1],
		start: Pos{
			Column: t.items[0].Column,
			Offset: t.items[0].Offset,
			Line:   t.items[0].Line,
		},
	}
}

func (t *Tree) parse() {
	// start := 0
	// inText := false
	// text := ""
	// for idx, _ := range t.items {
	// 	t.currPos = idx
	// if token.typ == itemText || token.typ == itemSpace {
	// 	inText = true
	// 	text = text + token.val
	// } else {
	// 	inText = false
	// }

	// if !inText && text != "" {
	// 	child := Tree{
	// 		typ:   NodeText,
	// 		items: nil,
	// 	}
	// 	n.children = append(n.children, child)
	// }
	// start = idx
	// }

	// if n.currPos > start {
	// 	child := Tree{
	// 		typ:   NodeParagraph,
	// 		items: n.items[start:n.currPos],
	// 	}
	// 	// child.parse()
	// 	n.children = append(n.children, child)
	// }
}
