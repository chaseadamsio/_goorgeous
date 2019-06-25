package transform

import (
	"fmt"
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/transform"
)

type HTMLOptions struct {
	Minify bool
}

type HTMLDocument struct {
	options             *HTMLOptions
	TOC                 []ast.Node
	FootnoteDefinitions map[int]ast.Node
}

func GenerateHTML(htmlTree *transform.HTMLTree, options *HTMLOptions) string {
	doc := &HTMLDocument{
		TOC:                 htmlTree.TOC,
		FootnoteDefinitions: htmlTree.FootnoteDefinitions,
		options:             options,
	}
	out := doc.walk(htmlTree.Root.Children())

	sortedFootnotes := make([]ast.Node, len(doc.FootnoteDefinitions))
	for place, node := range doc.FootnoteDefinitions {
		sortedFootnotes[int(place)-1] = node
	}

	var footnotes []string
	for _, node := range sortedFootnotes {
		footnotes = append(footnotes, doc.processFootnoteDefinitionNode(node.(*ast.FootnoteDefinitionNode)))
	}

	out = out + strings.Join(footnotes, "")

	return out
}

func (doc *HTMLDocument) walk(inAST []ast.Node) string {
	var out []string
	foundTableHeader := false

	for idx, child := range inAST {
		switch node := child.(type) {
		case *ast.HeadlineNode:
			out = append(out, doc.processHeadlineNode(node))
		case *ast.SectionNode:
			out = append(out, doc.processSectionNode(node))
		case *ast.HorizontalRuleNode:
			out = append(out, doc.processHorizontalNode(node))
		case *ast.ParagraphNode:
			out = append(out, doc.processParagraphNode(node))
		case *ast.ListNode:
			out = append(out, doc.processListNode(node))
		case *ast.ListItemNode:
			out = append(out, doc.processListItemNode(node))
		case *ast.LinkNode:
			out = append(out, doc.processLinkNode(node))
		case *ast.TableNode:
			out = append(out, doc.processTableNode(node))
		case *ast.TableRowNode:
			// foundTableHeader is to account for the test case that can be found
			// in testdata/in/table/basic.org under "multiple table rules:"
			// there are cases where a table can have multiple table rules and
			// this is the best way to account for that so that only one thead
			// is generated
			if !foundTableHeader && idx+1 < len(inAST) && inAST[idx+1].Type() == "TableRule" {
				out = append(out, doc.processTableHeaderNode(node))
				foundTableHeader = true
			} else if node.NodeType != "TableRule" {
				out = append(out, doc.processTableRowNode(node))
			}
		case *ast.TableCellNode:
			out = append(out, doc.processTableCellNode(node))
		case *ast.GreaterBlockNode:
			switch node.Name {
			case "SRC":
				out = append(out, doc.processGreaterBlockNode(node))
			case "EXAMPLE":
				out = append(out, doc.processGreaterBlockNode(node))
			case "QUOTE":
				out = append(out, doc.processQuoteBlockNode(node))
			case "VERSE":
				out = append(out, doc.processVerseBlockNode(node))
			default:
				out = append(out, doc.processSpecialGreaterBlockNode(node))
			}
		case *ast.FixedWidthNode:
			out = append(out, doc.processFixedWidthNode(node))
		case *ast.FootnoteReferenceNode:
			out = append(out, doc.processFootnoteReferenceNode(node))
		case *ast.TextNode:
			switch node.NodeType {
			case "Bold":
				out = append(out, doc.processBoldNode(node))
			case "Italic":
				out = append(out, doc.processItalicNode(node))
			case "Code":
				fallthrough // code and verbatim are processed the same
			case "Verbatim":
				out = append(out, doc.processVerbatimNode(node))
			case "Underline":
				out = append(out, doc.processUnderlineNode(node))
			case "EnDash":
				out = append(out, doc.processEnDashNode(node))
			case "MDash":
				out = append(out, doc.processMDashNode(node))

			default:
				out = append(out, node.Val)
			}
		default:

		}
	}
	return strings.Join(out, "")
}

func (doc *HTMLDocument) processHeadlineNode(node *ast.HeadlineNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<h%d>%s</h%d>", node.Depth, children, node.Depth)
}

func (doc *HTMLDocument) processHorizontalNode(node *ast.HorizontalRuleNode) string {
	return fmt.Sprintf("<hr />")
}

func (doc *HTMLDocument) processSectionNode(node *ast.SectionNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<div>\n%s\n</div>\n", children)
}

func (doc *HTMLDocument) processParagraphNode(node *ast.ParagraphNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<p>\n%s\n</p>\n", children)
}

func (doc *HTMLDocument) processListNode(node *ast.ListNode) string {
	listTyp := ""
	if node.ListType == "UNORDERED" {
		listTyp = "ul"
	}
	if node.ListType == "ORDERED" {
		listTyp = "ol"
	}
	children := doc.walk(node.Children())
	return fmt.Sprintf("<%s>\n\t%s\n\t</%s>\n", listTyp, children, listTyp)
}

func (doc *HTMLDocument) processListItemNode(node *ast.ListItemNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<li>\n%s\n</li>\n", children)
}

func (doc *HTMLDocument) processTableNode(node *ast.TableNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<table>\n%s\n</table>\n", children)
}

func (doc *HTMLDocument) processTableHeaderNode(node *ast.TableRowNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<thead>\n%s\n</thead>\n", children)
}

func (doc *HTMLDocument) processTableRowNode(node *ast.TableRowNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<tr>\n%s\n</tr>\n", children)
}

func (doc *HTMLDocument) processTableCellNode(node *ast.TableCellNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<td>\n%s\n</td>\n", children)
}

func (doc *HTMLDocument) processGreaterBlockNode(node *ast.GreaterBlockNode) string {
	className := strings.ToLower(node.Name)
	if node.Language != "" {
		className = className + " " + node.Language
	}

	return fmt.Sprintf("<pre class=\"%s\">\n%s\n</pre>\n", className, node.Value)
}

func (doc *HTMLDocument) processFootnoteReferenceNode(node *ast.FootnoteReferenceNode) string {

	return fmt.Sprintf("<sup><a id=\"fnr.%d\" href=\"#fn.%d\">%d</a></sup>", node.ID, node.ID, node.ID)
}

func (doc *HTMLDocument) processFootnoteDefinitionNode(node *ast.FootnoteDefinitionNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<div class=\"footdef\">\n<sup><a id=\"fn.%d\" href=\"#fnr.%d\">%d</a></sup><span>%s</span>\n</div>\n", node.ID, node.ID, node.ID, children)
}

func (doc *HTMLDocument) processQuoteBlockNode(node *ast.GreaterBlockNode) string {
	return fmt.Sprintf("<blockquote>\n<p>\n%s\n</p>\n</blockquote>\n", node.Value)
}

func (doc *HTMLDocument) processVerseBlockNode(node *ast.GreaterBlockNode) string {
	children := strings.Split(node.Value, "\n")
	inner := strings.Join(children, "<br />\n")
	return fmt.Sprintf("<div class=\"verse\">%s</div>\n", inner)
}

func (doc *HTMLDocument) processSpecialGreaterBlockNode(node *ast.GreaterBlockNode) string {
	return fmt.Sprintf("<div class=\"%s\">%s</div>\n", node.Name, node.Value)
}

func (doc *HTMLDocument) processFixedWidthNode(node *ast.FixedWidthNode) string {

	return fmt.Sprintf("<pre class=\"example\">\n%s\n</pre>\n", node.Value)
}

func (doc *HTMLDocument) processLinkNode(node *ast.LinkNode) string {
	children := node.Link // fallback link text is the link if no description provided
	if 0 < len(node.ChildrenNodes) {
		children = doc.walk(node.ChildrenNodes)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>", node.Link, children)
}

func (doc *HTMLDocument) processBoldNode(node *ast.TextNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<strong>%s</strong>", children)
}

func (doc *HTMLDocument) processItalicNode(node *ast.TextNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<em>%s</em>", children)
}

func (doc *HTMLDocument) processVerbatimNode(node *ast.TextNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<code>%s</code>", children)
}

func (doc *HTMLDocument) processStrikeThroughNode(node *ast.TextNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<span style=\"text-decoration: line-through\">%s</span>", children)
}

func (doc *HTMLDocument) processUnderlineNode(node *ast.TextNode) string {
	children := doc.walk(node.ChildrenNodes)
	return fmt.Sprintf("<span style=\"text-decoration:underline\">%s</span>", children)
}

func (doc *HTMLDocument) processMDashNode(node *ast.TextNode) string {
	return fmt.Sprintf("&mdash;")
}

func (doc *HTMLDocument) processEnDashNode(node *ast.TextNode) string {
	return fmt.Sprintf("&ndash;")
}
