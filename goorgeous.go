package goorgeous

import (
	"bytes"
	"fmt"

	"github.com/shurcooL/sanitized_anchor_name"
)

func OrgCommon(input []byte) []byte {
	var out bytes.Buffer
	tr := NewTree()
	tr.text = string(input)
	tr.startParse(lex("", tr.text), make([]*Tree, 1))
	tr.parse()
	h := &HTML{
		tmp: new(bytes.Buffer),
	}
	h.render(tr.Root.Nodes, &out)
	return out.Bytes()
}

type HTML struct {
	isInElement   bool
	foundNewLines int
	tmp           *bytes.Buffer
}

func (h *HTML) render(lNodes []Node, output *bytes.Buffer) {
	for _, n := range lNodes {
		switch n.Type() {
		case NodeParagraph:
			var tmp bytes.Buffer
			h.render(n.Nodes(), &tmp)
			output.WriteString("<p>" + tmp.String() + "</p>")
		case NodeHeadline:
			h.isInElement = true
			el := n.Element()
			var tmp bytes.Buffer
			h.render(n.Nodes(), &tmp)
			output.WriteString("<" + el + " id=\"" +
				sanitized_anchor_name.Create(n.ID()) + "\">")
			output.WriteString(tmp.String() + "</" + el + ">")
			h.isInElement = false
		case NodeNewLine:
			output.WriteByte('\n')
		case NodeImgOrLink:
			output.WriteString(n.String())
		case NodeText:
			output.WriteString(n.String())
		default:
			if n != nil {
				el := n.Element()
				output.WriteString("<" + el + ">" + n.String() + "</" + el + ">")
			}
		}
	}
}

var textFormat = "%s"

type Tree struct {
	Name      string
	Root      *ListNode
	text      string
	lex       *lexer
	treeSet   []*Tree
	token     [3]item
	peekCount int
}

type Node interface {
	Type() NodeType
	Element() string
	String() string
	ID() string
	Nodes() []Node
	tree() *Tree
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeText NodeType = iota // Plain text
	NodeList                 // A list of Nodes
	NodeNewLine

	NodeParagraph

	NodeHeadline
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

	NodeImgOrLink

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
	tr    *Tree
	Nodes []Node
}

func (t *Tree) newList() *ListNode {
	return &ListNode{tr: t, NodeType: NodeList}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func NewTree() *Tree {
	return new(Tree)
}

func (t *Tree) backup() {
	t.peekCount++
}

func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextItem()
	return t.token[0]
}

func (t *Tree) add() {
	t.treeSet = append(t.treeSet, t)
}

func (t *Tree) startParse(lex *lexer, treeSet []*Tree) {
	t.Root = nil
	t.lex = lex
	t.treeSet = treeSet
}

func (t *Tree) Parse(text string, treeSet []*Tree) (tree *Tree, err error) {
	t.startParse(lex(t.Name, text), treeSet)
	t.text = text
	t.parse()
	return t, nil
}

func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}
	return t.token[t.peekCount]
}

func (t *Tree) parse() {
	t.Root = t.newList()
	for t.peek().typ != itemEOF {
		n := t.getBlockNode()
		t.Root.append(n)
	}
}

func (t *Tree) getInlineNode() Node {
	switch token := t.next(); token.typ {
	case itemImgOrLinkOpen:
		var link, text string
		typ := "ANCHOR"
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			switch ne := t.next(); ne.typ {
			case itemImgOrLinkClose:
				break
			case itemImgOrLinkURL:
				link = ne.val
			case itemImgOrLinkText:
				text = ne.val
			case itemImgPre:
				typ = "IMG"
			}
		}
		if text == "" {
			text = link
		}
		return t.newImgOrLink(link, text, typ)
	case itemEmphasis:
		var text string
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			switch ne := t.next(); ne.typ {
			case itemText:
				text = ne.val
			case itemEmphasis:
				return t.newEmphasis(text)
			default:

				return t.newText(text)
			}
		}
	case itemVerbatim:
		var text string
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			switch ne := t.next(); ne.typ {
			case itemText:
				text = ne.val
			case itemVerbatim:
				return t.newVerbatim(text)
			default:

				return t.newText(text)
			}
		}
	case itemBold:
		var text string
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			switch ne := t.next(); ne.typ {
			case itemText:
				text = ne.val
			case itemBold:
				return t.newBold(text)
			default:
				return t.newText(text)
			}
		}
	default:
		return t.newText(token.val)
	}
	return nil
}

func (t *Tree) getBlockNode() Node {

	switch token := t.next(); token.typ {
	case itemHeadline:
		headingLevel := len(token.val)
		// for generating an id based on itemText
		var text []byte
		var children []Node
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			if t.peek().typ == itemText {
				text = append([]byte(t.peek().String()))
			}
			children = append(children, t.getInlineNode())
		}

		return t.newHeadline(text, headingLevel, children)

	case itemNewLine:
		return t.newNewLine(token.val)
	default:
		t.backup()
		var children []Node
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			children = append(children, t.getInlineNode())
		}
		return t.newParagraph(children)
	}
	return nil
}

// charMatches is a helper function to evaluate if two bytes are equal
func charMatches(a byte, b byte) bool {
	return a == b
}

// NODES
type NewLineNode struct {
	NodeType
	tr   *Tree
	Text []byte
}

func (t *Tree) newNewLine(text string) *NewLineNode {
	return &NewLineNode{tr: t, NodeType: NodeNewLine, Text: []byte(text)}
}

func (t *NewLineNode) Element() string {
	return ""
}

func (t *NewLineNode) ID() string {
	return ""
}

func (t *NewLineNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *NewLineNode) Nodes() []Node {
	return nil
}

func (t *NewLineNode) tree() *Tree {
	return t.tr
}

// NODES
type EmphasisNode struct {
	NodeType
	tr   *Tree
	Text []byte
}

func (t *Tree) newEmphasis(text string) *EmphasisNode {
	return &EmphasisNode{tr: t, NodeType: NodeEmphasis, Text: []byte(text)}
}

func (t *EmphasisNode) Element() string {
	return "em"
}

func (t *EmphasisNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *EmphasisNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *EmphasisNode) Nodes() []Node {
	return nil
}

func (t *EmphasisNode) tree() *Tree {
	return t.tr
}

type BoldNode struct {
	NodeType
	tr   *Tree
	Text []byte
}

func (t *Tree) newBold(text string) *BoldNode {
	return &BoldNode{tr: t, NodeType: NodeBold, Text: []byte(text)}
}

func (t *BoldNode) Element() string {
	return "strong"
}

func (t *BoldNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *BoldNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *BoldNode) Nodes() []Node {
	return nil
}

func (t *BoldNode) tree() *Tree {
	return t.tr
}

type VerbatimNode struct {
	NodeType
	tr   *Tree
	Text []byte
}

func (t *Tree) newVerbatim(text string) *VerbatimNode {
	return &VerbatimNode{tr: t, NodeType: NodeVerbatim, Text: []byte(text)}
}

func (t *VerbatimNode) Element() string {
	return "code"
}

func (t *VerbatimNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *VerbatimNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *VerbatimNode) Nodes() []Node {
	return nil
}

func (t *VerbatimNode) tree() *Tree {
	return t.tr
}

type ImgOrLinkNode struct {
	NodeType
	tr       *Tree
	Link     string
	Text     string
	typ      string
	children []Node
}

func (t *Tree) newImgOrLink(link, text, typ string) *ImgOrLinkNode {
	return &ImgOrLinkNode{tr: t, NodeType: NodeImgOrLink, Link: link, Text: text, typ: typ}
}

func (t *ImgOrLinkNode) Element() string {
	return fmt.Sprintf(textFormat, "p")
}

func (t *ImgOrLinkNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *ImgOrLinkNode) String() string {
	switch t.typ {
	case "IMG":
		return fmt.Sprintf("<img src=\"%s\" />", t.Link)
	default:
		return fmt.Sprintf("<a href=\"%s\" title=\"%s\">%s</a>", t.Link, t.Text, t.Text)
	}
}

func (t *ImgOrLinkNode) Nodes() []Node {
	return t.children
}

func (t *ImgOrLinkNode) tree() *Tree {
	return t.tr
}

type TextNode struct {
	NodeType
	tr       *Tree
	Text     []byte
	children []Node
}

func (t *Tree) newText(text string) *TextNode {
	return &TextNode{tr: t, NodeType: NodeText, Text: []byte(text)}
}

func (t *TextNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *TextNode) Element() string {
	return ""
}

func (t *TextNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *TextNode) Nodes() []Node {
	return t.children
}

func (t *TextNode) tree() *Tree {
	return t.tr
}

type HeadlineNode struct {
	NodeType
	Level    int
	tr       *Tree
	Text     []byte
	children []Node
}

func (t *Tree) newHeadline(text []byte, level int, children []Node) *HeadlineNode {
	return &HeadlineNode{tr: t, NodeType: NodeHeadline, Text: text, Level: level, children: children}
}

func (t *HeadlineNode) ID() string {
	var id []byte
	for _, c := range t.children {
		id = append(id, []byte(c.ID())...)
	}
	return string(id)
}

func (t *HeadlineNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *HeadlineNode) Element() string {
	return fmt.Sprintf("h%d", t.Level)
}

func (t *HeadlineNode) Nodes() []Node {
	return t.children
}

func (t *HeadlineNode) tree() *Tree {
	return t.tr
}

type ParagraphNode struct {
	NodeType
	tr       *Tree
	Text     []byte
	children []Node
}

func (t *Tree) newParagraph(children []Node) *ParagraphNode {
	return &ParagraphNode{tr: t, NodeType: NodeParagraph, children: children}
}

func (t *ParagraphNode) ID() string {
	return ""
}

func (t *ParagraphNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *ParagraphNode) Element() string {
	return fmt.Sprintf("p")
}

func (t *ParagraphNode) Nodes() []Node {
	return t.children
}

func (t *ParagraphNode) tree() *Tree {
	return t.tr
}
