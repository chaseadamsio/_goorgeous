package transform

import (
	"fmt"
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
)

func TransformToHTML(root *ast.RootNode) string {
	return walk(root.Children())
}

func walk(inAST []ast.Node) string {
	var out []string

	for _, child := range inAST {
		switch node := child.(type) {
		case *ast.HeadlineNode:
			out = append(out, processHeadlineNode(node))
		case *ast.ListNode:
			out = append(out, processListNode(node))
		case *ast.ListItemNode:
			out = append(out, processListItemNode(node))
		case *ast.TextNode:
			out = append(out, node.Val)
		default:
		}
	}

	return strings.Join(out, "\n")
}

func processHeadlineNode(node *ast.HeadlineNode) string {
	return fmt.Sprintf("<h%d>%s</h%d>", node.Depth, node.Children()[0], node.Depth)
}

func processListNode(node *ast.ListNode) string {
	listTyp := ""
	if node.ListType == "UNORDERED" {
		listTyp = "ul"
	}
	if node.ListType == "ORDERED" {
		listTyp = "ol"
	}
	children := walk(node.Children())
	return fmt.Sprintf("<%s>%s</%s>", listTyp, children, listTyp)
}

func processListItemNode(node *ast.ListItemNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<li>%s</li>", children)
}
