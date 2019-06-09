package transform

import (
	"github.com/chaseadamsio/goorgeous/ast"
)

type HTMLNode struct {
	Type  string
	Value string
}

func TransformToHTML(inAST *ast.RootNode) []HTMLNode {
	out := flattenTree(inAST.Children())
	return out
}

var nonHTML = map[string]struct{}{
	"Keyword": struct{}{},
	"Section": struct{}{},
}

func flattenTree(inAST []ast.Node) []HTMLNode {
	var outAST []HTMLNode
	for _, child := range inAST {
		if _, found := nonHTML[string(child.Type())]; !found {
			outAST = append(outAST, HTMLNode{
				Type:  string(child.Type()),
				Value: child.String(),
			})
		}
		if len(child.Children()) > 0 {
			outAST = append(outAST, flattenTree(child.Children())...)
		}
	}
	return outAST
}
