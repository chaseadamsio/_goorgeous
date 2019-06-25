package transform

import (
	"encoding/json"
	"fmt"

	"github.com/chaseadamsio/goorgeous/ast"
)

type HTMLTree struct {
	TOC                 []ast.Node
	footnoteLabelToID   map[string]int
	FootnoteDefinitions map[int]ast.Node
	footnoteCount       int
	Root                *ast.RootNode
}

func (tree *HTMLTree) String() string {
	out, _ := json.MarshalIndent(tree, "", "  ")
	return string(out)
}

func TransformToHTML(root *ast.RootNode) (tree *HTMLTree) {

	tree = &HTMLTree{
		TOC:                 make([]ast.Node, 1),
		footnoteLabelToID:   make(map[string]int, 1),
		FootnoteDefinitions: make(map[int]ast.Node, 1),
		Root:                ast.NewRootNode(root.Start, root.End),
		footnoteCount:       0,
	}

	tree.walk(root, root.Children())

	return tree
}

func (tree *HTMLTree) walk(parent ast.Node, children []ast.Node) {

	for _, child := range children {
		switch node := child.(type) {
		case *ast.HeadlineNode:
			processHeadlineNode(tree, node)
		case *ast.SectionNode:
			processSectionNode(tree, node)
		case *ast.FootnoteDefinitionNode:
			processFootnoteDefinitionNode(tree, node)
		default:
			parent.Append(child)
			findInlineFootnoteReferences(tree, child)
		}
	}
}

func findInlineFootnoteReferences(tree *HTMLTree, parent ast.Node) {
	for _, child := range parent.Children() {
		if node, ok := child.(*ast.FootnoteReferenceNode); ok {
			tree.footnoteCount++
			label := node.Label
			if label == "" {
				label = fmt.Sprintf("%d", tree.footnoteCount)
			}

			tree.footnoteLabelToID[label] = tree.footnoteCount

			node.Label = label
			node.ID = tree.footnoteCount

			if node.ReferenceType == "INLINE" {
				if refDef, ok := node.Children()[0].(*ast.FootnoteDefinitionNode); ok {
					refDef.Label = label
					refDef.ID = tree.footnoteCount
					processFootnoteDefinitionNode(tree, refDef)
				}
			}

		}
	}
}

func processHeadlineNode(tree *HTMLTree, node *ast.HeadlineNode) {
	newNode := node.Copy()
	tree.Root.Append(newNode)
	newNode.ChildrenNodes = make([]ast.Node, 1)
	tree.walk(newNode, node.Children())
}

func processHorizontalNode(tree *HTMLTree, node *ast.HorizontalRuleNode) {
	tree.Root.Append(node)
}

func processSectionNode(tree *HTMLTree, node *ast.SectionNode) {
	newNode := node.Copy()
	tree.Root.Append(newNode)
	newNode.ChildrenNodes = make([]ast.Node, 1)
	tree.walk(newNode, node.Children())
}

func processListNode(tree *HTMLTree, node *ast.ListNode) {
	newNode := node.Copy()
	tree.Root.Append(newNode)
	newNode.ChildrenNodes = make([]ast.Node, 1)
	tree.walk(newNode, node.Children())
}

func processTableNode(tree *HTMLTree, node *ast.TableNode) {
	newNode := node.Copy()
	tree.Root.Append(newNode)
	newNode.ChildrenNodes = make([]ast.Node, 1)
	tree.walk(newNode, node.Children())
}

func processGreaterBlockNode(tree *HTMLTree, node *ast.GreaterBlockNode) {
	newNode := node.Copy()
	tree.Root.Append(newNode)
	tree.walk(newNode, node.Children())
}

func processFootnoteDefinitionNode(tree *HTMLTree, node *ast.FootnoteDefinitionNode) {
	newNode := node.Copy()
	tree.Root.Append(newNode)
	newNode.ChildrenNodes = make([]ast.Node, 1)
	tree.walk(newNode, node.Children())

	if id, found := tree.footnoteLabelToID[node.Label]; found {
		newNode.ID = id
		tree.FootnoteDefinitions[id] = newNode
	}
}
