package goorgeous

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	for _, tc := range lexTests {
		if !strings.HasPrefix(tc.name, "headline - deep") {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("children for", tc.name)
			// _ = Parse(tc.input)
			ast := Parse(tc.input)
			fmt.Println(ast)
			// if tc.ast != nil {
			// 	for _, child := range ast.Children() {
			// 		var path []NodeType
			// 	}
			// }
		})
	}
}

func getPathForNode(n Node, path []NodeType) []NodeType {
	path = append(path, n.Type())
	if n.Parent().Type() != "Root" {
		path = getPathForNode(n.Parent(), path)
	} else {
		var reversedPath []NodeType
		path = append(path, "Root")
		for idx := len(path) - 1; idx >= 0; idx-- {
			reversedPath = append(reversedPath, path[idx])
		}
		return reversedPath
	}
	return path
}
