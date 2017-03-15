package parse

import "bytes"

type NodeType int

const (
	NodeParagraph NodeType = iota
)

type ListNode struct {
	NodeType
	val   string
	pos   int
	Nodes []ListNode
}

func RenderAsHTML(ln []*ListNode) string {
	var out bytes.Buffer
	render(&out, ln)
	return out.String()
}

func render(out *bytes.Buffer, ln []*ListNode) {
	for _, node := range ln {
		if node.NodeType == NodeParagraph {
			out.WriteString("<p>" + node.val + "</p>")
		}
	}

}

func Parse(text string) ([]*ListNode, error) {
	l := lex(text)
	nodes := make([]*ListNode, 0)

	for {
		var typ NodeType
		item := l.nextItem()

		if item.typ == itemEOF {
			break
		}

		if item.typ == itemText {
			typ = NodeParagraph
		}
		nodes = append(nodes, &ListNode{NodeType: typ, val: item.val})

	}

	return nodes, nil
}
