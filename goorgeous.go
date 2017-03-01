package goorgeous

type Tree struct {
	Name    string
	Root    *ListNode
	text    string
	funcs   []map[string]interface{}
	lex     *lexer
	treeSet map[string]*Tree
}

type Node interface {
	Type() NodeType
	String() string
	Copy() Node
	Position() Pos
	tree() *Tree
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText = iota
	NodeNewLine
	NodeH1
	NodeH2
	NodeH3
	NodeH4
	NodeH5
	NodeH6
	NodeStatus
	NodePriority
	NodeTags
	NodePropertyDrawer
	NodeProperty
	NodeKeyword
	NodeComment
	NodeHorizontalRule
	NodeEmphasis
	NodeBold
	NodeStrikethrough
	NodeVerbatim
	NodeCode
	NodeUnderline
	NodeDefinitionList
	NodeOrderedList
	NodeUnorderedList
	NodeListItem
	NodeTable
	NodeTH
	NodeTR
	NodeTD
	NODEHTML
)

type ListNode struct {
	NodeType
	Pos
	tr    *Tree
	Nodes []Node
}

func New(name string, funcs ...map[string]interface{}) *Tree {
	return &Tree{
		Name:  name,
		funcs: funcs,
	}
}

func (t *Tree) startParse(funcs []map[string]interface{}, lex *lexer, treeSet map[string]*Tree) {
	t.Root = nil
	t.lex = lex
	t.funcs = funcs
	t.treeSet = treeSet
}

func (t *Tree) Parse(text string, treeSet map[string]*Tree, funcs ...map[string]interface{}) (tree *Tree, err error) {
	t.startParse(funcs, lex(t.Name, text), treeSet)
	t.text = text
	t.parse()
	return t, nil
}

func (t *Tree) parse() (items []item) {
	for {
		item := t.lex.nextItem()
		items = append(items, item)
		if item.typ == itemEOF {
			break
		}
	}
	return
}

// charMatches is a helper function to evaluate if two bytes are equal
func charMatches(a byte, b byte) bool {
	return a == b
}
