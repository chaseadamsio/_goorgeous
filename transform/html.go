package transform

import (
	"fmt"

	"github.com/chaseadamsio/goorgeous/ast"
)

func TransformToHTML(inAST *ast.RootNode) {
	fmt.Println(inAST)
	// fmt.Printf("found headline node: <h%d>%s</h%d>", hl.Depth, hl, hl.Depth)
}

func flattenTree(inAST []ast.Node) []ast.Node {
	return nil

	// for _, child := range inAST {
	// 	if hl, ok := child.(*ast.HeadlineNode); ok {
	// 		out := hl.ToJSON()
	// 		// fmt.Printf("found headline node: <h%d>%s</h%d>", hl.Depth, hl, hl.Depth)
	// 	}
	// }
	// var outAST []HTMLNode
	// for _, child := range inAST {
	// 	if _, found := nonHTML[string(child.Type())]; !found {
	// 		outAST = append(outAST, HTMLNode{
	// 			Type:  string(child.Type()),
	// 			Value: child.String(),
	// 		})
	// 	}
	// 	if len(child.Children()) > 0 {
	// 		outAST = append(outAST, flattenTree(child.Children())...)
	// 	}
	// }
	// return outAST
}
