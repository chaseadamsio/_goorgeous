package goorgeous

// NodeType is the type of a node
type NodeType int

// Type returns itself and provides an easy default implementation
// for embedding in a Node. Embedded in all non-trivial Nodes.
func (t NodeType) Type() NodeType {
	return t
}

const (
	// NodeRoot is the type for a root node
	NodeRoot NodeType = iota
	// NodeParagraph is the type for a paragraph node
	NodeParagraph
	// NodeHeadline is the type for a headline node
	NodeHeadline
	// NodeText is the type for a text node
	NodeText
)

// Node is an interface
type Node interface {
	String() string
}

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos struct {
	Line   int
	Column int
	Offset int
}

// Position is a method to get the position of an element
func (p Pos) Position() Pos {
	return p
}

// ParagraphNode is a node paragraph
type ParagraphNode struct {
	NodeType
	Start Pos
	End   Pos
	tr    *Tree
	items []item
}

func newParagraphNode(start, end Pos, items []item) *ParagraphNode {
	return &ParagraphNode{
		NodeType: NodeParagraph,
		Start:    start,
		End:      end,
		items:    items,
	}
}

// Type returns the type of node this is
func (n *ParagraphNode) Type() NodeType {
	return n.NodeType
}

// TODO: implement real String methods
func (n *ParagraphNode) String() string {
	return ""
}

// HeadlineNode is a headline node
type HeadlineNode struct {
	NodeType
	Start Pos
	End   Pos
	tr    *Tree
	items []item
}

func newHeadlineNode(start, end Pos, items []item) *HeadlineNode {
	return &HeadlineNode{
		NodeType: NodeHeadline,
		Start:    start,
		End:      end,
		items:    items,
	}
}

func (n *HeadlineNode) String() string {
	return ""
}

// Type returns the type of node this is
func (n *HeadlineNode) Type() NodeType {
	return n.NodeType
}

// TextNode is a text node
type TextNode struct {
	NodeType
	Start Pos
	End   Pos
	tr    *Tree
	items []item
}

func newTextNode(start, end Pos, items []item) *TextNode {
	return &TextNode{
		NodeType: NodeText,
		Start:    start,
		End:      end,
		items:    items,
	}
}

func (n *TextNode) String() string {
	return ""
}

// Type returns the type of node this is
func (n *TextNode) Type() NodeType {
	return n.NodeType
}
