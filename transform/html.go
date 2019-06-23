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

	for idx, child := range inAST {
		switch node := child.(type) {
		case *ast.HeadlineNode:
			out = append(out, processHeadlineNode(node))
		case *ast.SectionNode:
			out = append(out, processSectionNode(node))
		case *ast.HorizontalRuleNode:
			out = append(out, processHorizontalNode(node))
		case *ast.ParagraphNode:
			out = append(out, processParagraphNode(node))
		case *ast.ListNode:
			out = append(out, processListNode(node))
		case *ast.ListItemNode:
			out = append(out, processListItemNode(node))
		case *ast.LinkNode:
			out = append(out, processLinkNode(node))
		case *ast.TableNode:
			out = append(out, processTableNode(node))
		case *ast.TableRowNode:
			if idx+1 < len(inAST) && inAST[idx+1].Type() == "TableRule" {
				out = append(out, processTableHeaderNode(node))
			} else if node.NodeType != "TableRule" {
				out = append(out, processTableRowNode(node))
			}
		case *ast.TableCellNode:
			out = append(out, processTableCellNode(node))
		case *ast.TextNode:
			switch node.NodeType {
			case "Bold":
				out = append(out, processBoldNode(node))
			case "Italic":
				out = append(out, processItalicNode(node))
			case "Verbatim":
				out = append(out, processVerbatimNode(node))
			case "Underline":
				out = append(out, processUnderlineNode(node))
			default:
				out = append(out, node.Val)
			}
		default:

		}
	}

	return strings.Join(out, "")
}

func processHeadlineNode(node *ast.HeadlineNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<h%d>%s</h%d>", node.Depth, children, node.Depth)
}

func processHorizontalNode(node *ast.HorizontalRuleNode) string {
	return fmt.Sprintf("<hr />")
}

func processSectionNode(node *ast.SectionNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<div>\n%s\n</div>\n", children)
}

func processParagraphNode(node *ast.ParagraphNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<p>\n%s\n</p>\n", children)
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
	return fmt.Sprintf("<%s>\n\t%s\n\t</%s>\n", listTyp, children, listTyp)
}

func processListItemNode(node *ast.ListItemNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<li>%s</li>", children)
}

func processTableNode(node *ast.TableNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<table>%s</table>", children)
}

func processTableHeaderNode(node *ast.TableRowNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<thead>%s</thead>", children)
}

func processTableRowNode(node *ast.TableRowNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<tr>%s</tr>", children)
}

func processTableCellNode(node *ast.TableCellNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<td>%s</td>", children)
}

func processLinkNode(node *ast.LinkNode) string {
	children := node.Link // fallback link text is the link if no description provided
	if 0 < len(node.ChildrenNodes) {
		children = walk(node.ChildrenNodes)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", node.Link, children)
}

func processBoldNode(node *ast.TextNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<strong>%s</strong>", children)
}

func processItalicNode(node *ast.TextNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<em>%s</em>", children)
}

func processVerbatimNode(node *ast.TextNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<code>%s</code>", children)
}

func processStrikeThroughNode(node *ast.TextNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<span style=\"text-decoration: line-through\">%s</span>", children)
}

func processUnderlineNode(node *ast.TextNode) string {
	children := walk(node.ChildrenNodes)
	return fmt.Sprintf("<span style=\"text-decoration:underline\">%s</span>", children)
}
