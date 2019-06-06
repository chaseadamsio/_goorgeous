package goorgeous

import (
	"encoding/json"
	"strings"
)

// NodeType is the type of a node
type NodeType string

type Node interface {
	Type() NodeType
	Parent() Node
	String() string
	Children() []Node
	Append(Node)
}

func newRootNode() *RootNode {
	node := &RootNode{
		NodeType: "Root",
	}

	return node
}

type RootNode struct {
	NodeType
	parent   Node
	children []Node
}

func (n *RootNode) String() string {
	tree := []treeTyp{
		{
			"Root",
			getChildren(n),
			"",
		},
	}
	out, _ := json.MarshalIndent(tree, "", "  ")
	return string(out)
}

type treeTyp struct {
	NodeType
	Children []treeTyp
	Value    string
}

func getChildren(node Node) []treeTyp {
	var tree []treeTyp

	if len(node.Children()) == 0 {
		return nil
	}

	for _, child := range node.Children() {
		tree = append(tree, treeTyp{
			child.Type(),
			getChildren(child),
			child.String(),
		})
	}

	return tree
}

func (n RootNode) Tree() RootNode {
	return n
}

func (n RootNode) Type() NodeType {
	return n.NodeType
}

func (n RootNode) Children() []Node {
	return n.children
}

func (n *RootNode) Parent() Node {
	return n.parent
}

func (n *RootNode) Append(child Node) {
	n.children = append(n.children, child)
}

func newSectionNode() *SectionNode {
	node := &SectionNode{
		NodeType: "Section",
	}

	return node
}

type SectionNode struct {
	NodeType
	parent   Node
	children []Node
}

func (n SectionNode) Type() NodeType {
	return n.NodeType
}

func (n SectionNode) String() string {
	return ""
}

func (n SectionNode) Children() []Node {
	return n.children
}

func (n *SectionNode) Parent() Node {
	return n.parent
}

func (n *SectionNode) Append(child Node) {
	n.children = append(n.children, child)
}

func newTextNode(start, end int, parent Node, items []item) *TextNode {
	var values []string
	for _, item := range items {
		values = append(values, item.val)
	}
	node := &TextNode{
		NodeType: "Text",
		val:      strings.Join(values, ""),
		parent:   parent,
	}

	return node
}

type Pos struct {
	Start int
	End   int
}

type TextNode struct {
	NodeType
	parent   Node
	children []Node
	val      string
	start    int
	end      int
}

func (n TextNode) Type() NodeType {
	return n.NodeType
}

func (n TextNode) String() string {
	return n.val
}

func (n TextNode) Children() []Node {
	return n.children
}
func (n *TextNode) Parent() Node {
	return n.parent
}
func (n *TextNode) Append(child Node) {
	n.children = append(n.children, child)
}

type ParagraphNode struct {
	NodeType
	parent   Node
	start    int
	end      int
	children []Node
}

func newParagraphNode(start, end int, parent Node, items []item) *ParagraphNode {
	node := &ParagraphNode{
		NodeType: "Paragraph",
		parent:   parent,
		start:    start,
		end:      end,
	}

	return node
}

// Type returns the type of node this is
func (n *ParagraphNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *ParagraphNode) String() string {
	return ""
}

func (n ParagraphNode) Children() []Node {
	return n.children
}
func (n *ParagraphNode) Parent() Node {
	return n.parent
}
func (n *ParagraphNode) Append(child Node) {
	n.children = append(n.children, child)
}

type HeadlineNode struct {
	NodeType
	start, end int
	depth      int
	parent     Node
	rawvalue   string
	children   []Node
	Keyword    string
}

func newHeadlineNode(start, end, depth int, parent Node, items []item) *HeadlineNode {
	node := &HeadlineNode{
		NodeType: "Headline",
		depth:    depth,
		parent:   parent,
		start:    start,
		end:      end,
	}

	node.parse(items)

	return node
}

func (n *HeadlineNode) parse(items []item) {
	var headlineVal []string
	for idx, item := range items {
		if hasKeyword(idx, items) {
			n.Keyword = item.val
		}
		headlineVal = append(headlineVal, item.val)
	}
	n.rawvalue = strings.Join(headlineVal, "")
}

var keywords = map[string]struct{}{
	"TODO": struct{}{},
	"DONE": struct{}{},
}

func hasKeyword(idx int, items []item) bool {
	// keywords will only _ever_ occur in the first space
	if idx != 0 {
		return false
	}
	if _, ok := keywords[items[idx].val]; ok {
		return true
	}
	return false
}

// Type returns the type of node this is
func (n *HeadlineNode) Type() NodeType {
	return n.NodeType
}

// Type returns the type of node this is
func (n *HeadlineNode) String() string {
	return n.rawvalue
}

func (n HeadlineNode) Children() []Node {
	return n.children
}

func (n *HeadlineNode) Parent() Node {
	return n.parent
}
func (n HeadlineNode) Depth() int {
	return n.depth
}

func (n *HeadlineNode) Append(child Node) {
	n.children = append(n.children, child)
}
