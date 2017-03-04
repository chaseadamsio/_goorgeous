package goorgeous

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/shurcooL/sanitized_anchor_name"
)

func OrgCommon(input []byte) []byte {
	var out bytes.Buffer
	tr := NewTree()
	tr.text = string(input)
	tr.startParse(lex("", tr.text), make([]*Tree, 1))
	tr.parse()
	for _, n := range tr.Root.Nodes {
		out.Write(n.generateHTML())
	}
	return out.Bytes()
}

type HTML struct {
	isInElement   bool
	foundNewLines int
	tmp           *bytes.Buffer
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
	setAttr(string, interface{})
	getAttr(string) interface{}
	generateHTML() []byte
	String() string
	Nodes() []Node
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

	NodeBlock

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
			if t.peek().typ == itemImgOrLinkClose {
				t.next()
				break
			}
			switch ne := t.next(); ne.typ {
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
	case itemCode:
		var text string
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			switch ne := t.next(); ne.typ {
			case itemText:
				text = ne.val
			case itemCode:
				return t.newVerbatim(text)
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
	case itemUnderline:
		var text string
		for t.peek().typ != itemNewLine && t.peek().typ != itemEOF {
			switch ne := t.next(); ne.typ {
			case itemText:
				text = ne.val
			case itemUnderline:
				return t.newUnderline(text)
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
	case itemBlock:
		el := token.val[8:]
		syntax := ""
		if el == "SRC" || el == "EXAMPLE" {
			syntax = t.next().val
		}
		var tmp bytes.Buffer
		for t.peek().typ != itemBlock && t.peek().typ != itemEOF {
			n := t.next()
			tmp.WriteString(n.val)
		}
		t.next()
		return t.newBlock(el, syntax, tmp.Bytes())

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
}

// charMatches is a helper function to evaluate if two bytes are equal
func charMatches(a byte, b byte) bool {
	return a == b
}

// NODES
type NewLineNode struct {
	NodeType
	tr    *Tree
	Text  []byte
	attrs map[string]interface{}
}

func (t *Tree) newNewLine(text string) *NewLineNode {
	return &NewLineNode{tr: t, NodeType: NodeNewLine, Text: []byte(text)}
}

func (t *NewLineNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *NewLineNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func setAttr(attrs map[string]interface{}, k string, v interface{}) {
	attrs[k] = v
}

func getAttr(attrs map[string]interface{}, k string) interface{} {
	if attrs[k] != nil {
		return attrs[k]
	}
	return nil
}

func (t *NewLineNode) generateHTML() []byte {
	return []byte("\n")
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
	tr    *Tree
	Text  []byte
	attrs map[string]interface{}
}

func (t *Tree) newEmphasis(text string) *EmphasisNode {
	return &EmphasisNode{tr: t, NodeType: NodeEmphasis, Text: []byte(text)}
}

func (t *EmphasisNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *EmphasisNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *EmphasisNode) generateHTML() []byte {
	var tmp bytes.Buffer

	el := t.Element()

	tmp.Write(t.Text)

	for _, n := range t.Nodes() {
		tmp.Write(n.generateHTML())
	}

	return []byte("<" + el + ">" + tmp.String() + "</" + el + ">")
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
	tr    *Tree
	Text  []byte
	attrs map[string]interface{}
}

func (t *Tree) newBold(text string) *BoldNode {
	return &BoldNode{tr: t, NodeType: NodeBold, Text: []byte(text)}
}

func (t *BoldNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *BoldNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *BoldNode) generateHTML() []byte {
	var tmp bytes.Buffer

	el := t.Element()

	tmp.Write(t.Text)

	for _, n := range t.Nodes() {
		tmp.Write(n.generateHTML())
	}

	return []byte("<" + el + ">" + tmp.String() + "</" + el + ">")
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

type UnderlineNode struct {
	NodeType
	tr    *Tree
	Text  []byte
	attrs map[string]interface{}
}

func (t *Tree) newUnderline(text string) *UnderlineNode {
	return &UnderlineNode{tr: t, NodeType: NodeUnderline, Text: []byte(text)}
}

func (t *UnderlineNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *UnderlineNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *UnderlineNode) generateHTML() []byte {
	var tmp bytes.Buffer

	el := t.Element()

	tmp.Write(t.Text)

	for _, n := range t.Nodes() {
		tmp.Write(n.generateHTML())
	}

	return []byte("<" + el + " style=\"text-decoration: underline;\">" + tmp.String() + "</" + el + ">")
}

func (t *UnderlineNode) Element() string {
	return "span"
}

func (t *UnderlineNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *UnderlineNode) String() string {
	return fmt.Sprintf(textFormat, t.Text)
}

func (t *UnderlineNode) Nodes() []Node {
	return nil
}

type VerbatimNode struct {
	NodeType
	tr    *Tree
	Text  []byte
	attrs map[string]interface{}
}

func (t *Tree) newVerbatim(text string) *VerbatimNode {
	return &VerbatimNode{tr: t, NodeType: NodeVerbatim, Text: []byte(text)}
}

func (t *VerbatimNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *VerbatimNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *VerbatimNode) generateHTML() []byte {
	var tmp bytes.Buffer

	el := t.Element()

	tmp.Write(t.Text)

	for _, n := range t.Nodes() {
		tmp.Write(n.generateHTML())
	}

	return []byte("<" + el + ">" + tmp.String() + "</" + el + ">")
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
	attrs    map[string]interface{}
}

func (t *Tree) newImgOrLink(link, text, typ string) *ImgOrLinkNode {
	return &ImgOrLinkNode{tr: t, NodeType: NodeImgOrLink, Link: link, Text: text, typ: typ}
}

func (t *ImgOrLinkNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *ImgOrLinkNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *ImgOrLinkNode) Element() string {
	return fmt.Sprintf(textFormat, "p")
}

func (t *ImgOrLinkNode) ID() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *ImgOrLinkNode) String() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *ImgOrLinkNode) generateHTML() []byte {
	switch t.typ {
	case "IMG":
		return []byte(fmt.Sprintf("<img src=\"%s\" />", t.Link))
	default:
		return []byte(fmt.Sprintf("<a href=\"%s\" title=\"%s\">%s</a>", t.Link, t.Text, t.Text))
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
	attrs    map[string]interface{}
}

func (t *Tree) newText(text string) *TextNode {
	return &TextNode{tr: t, NodeType: NodeText, Text: []byte(text)}
}

func (t *TextNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *TextNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *TextNode) String() string {
	return fmt.Sprintf("%s", t.Text)
}

func (t *TextNode) Element() string {
	return ""
}

func (t *TextNode) generateHTML() []byte {
	return []byte(t.Text)
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
	attrs    map[string]interface{}
}

func (t *Tree) newHeadline(text []byte, level int, children []Node) *HeadlineNode {
	return &HeadlineNode{tr: t, NodeType: NodeHeadline, Text: text, Level: level, children: children}
}

func (t *HeadlineNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *HeadlineNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *HeadlineNode) generateHTML() []byte {
	var tmp bytes.Buffer

	el := t.Element()
	for _, n := range t.Nodes() {
		tmp.Write(n.generateHTML())
	}

	return []byte("<" + el + " id=\"" +
		sanitized_anchor_name.Create(t.ID()) + "\">" +
		tmp.String() + "</" + el + ">")
}

func (t *HeadlineNode) ID() string {
	var id []byte
	for _, c := range t.children {
		id = append(id, []byte(c.String())...)
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
	attrs    map[string]interface{}
}

func (t *Tree) newParagraph(children []Node) *ParagraphNode {
	return &ParagraphNode{tr: t, NodeType: NodeParagraph, children: children}
}

func (t *ParagraphNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *ParagraphNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *ParagraphNode) generateHTML() []byte {
	var tmp bytes.Buffer

	for _, n := range t.Nodes() {
		tmp.Write(n.generateHTML())
	}

	return []byte("<p>" + tmp.String() + "</p>")
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

type BlockNode struct {
	NodeType
	tr       *Tree
	el       string
	syntax   string
	Text     []byte
	children []Node
	attrs    map[string]interface{}
}

func (t *Tree) newBlock(element, syntax string, text []byte) *BlockNode {
	return &BlockNode{tr: t, el: element, NodeType: NodeBlock, Text: text, syntax: syntax}
}

func (t *BlockNode) setAttr(k string, v interface{}) {
	setAttr(t.attrs, k, v)
}

func (t *BlockNode) getAttr(k string) interface{} {
	return getAttr(t.attrs, k)
}

func (t *BlockNode) generateHTML() []byte {
	var open, close string
	switch t.el {
	case "SRC":
		open = "<pre><code"
		if t.syntax != "" {
			open += " class=\"language-" + strings.Trim(string(t.syntax), " ") + "\""
		}
		open += ">"
		close = "</code></pre>"
	case "EXAMPLE":
		open = "<pre><code"
		if t.syntax != "" {
			open += " class=\"language-" + strings.Trim(string(t.syntax), " ") + "\""
		}
		open += ">"
		close = "</code></pre>"
	case "QUOTE":
		open = "<blockquote>"
		close = "</blockquote>"
	case "CENTER":
		open = "<center>"
		close = "</center>"
	default:
		open = "<div>"
		close = "</div>"
	}
	return []byte(open + t.String() + close)
}

func (t *BlockNode) ID() string {
	return ""
}

func (t *BlockNode) String() string {
	return string(t.Text)
}

func (t *BlockNode) Element() string {
	return t.el
}

func (t *BlockNode) Nodes() []Node {
	return t.children
}
