package goorgeous

import (
	"testing"
)

func TestParse(t *testing.T) {

	// type parseTest struct {
	// 	name     string
	// 	value    string
	// 	expected *Tree
	// }

	n := Parse("* headline")
	if n.typ != NodeRoot {
		t.Errorf("expected parent to be of type NodeRoot")
	}
	// t.Log("\n\nn is:", fmt.Sprintf("%+v", n.start))
	// t.Log("\n\nn is:", fmt.Sprintf("%+v", n.end))

	// expected := []Tree{
	// 	Tree{
	// 		typ: NodeHeadline,
	// 	},
	// }
	// if !cmp.Equal(n.children, expected) {
	// 	t.Errorf("Actual: %#v\n\nExpected:%#v", n.children, expected)
	// }
	// t.Log("second parse:")
	// n = Parse("* headline\nthis is a new line")
	// t.Log("\n\nn is:", fmt.Sprintf("%+v", n.children))

	// n = Parse("* headline\n\n\nthis is a new line")
	// t.Log("\n\nn is:", fmt.Sprintf("%+v", n.children))

	// t.Log("fourth parse:")
	// n = Parse("* headline\n\n\n\n\n\n\nthis is a new line")
	// t.Log("\n\nn is:", fmt.Sprintf("%+v", n.children))
	// testCases := []parseTest{
	// 	{"basic", "* headline", &Tree{
	// 		Root: &ListNode{
	// 			NodeType: NodeRoot,
	// 			Nodes: []Node{
	// 				&NodeHeadline{
	// 					Nodes: []Node{
	// 						&ListNode{
	// 							NodeType: NodeText,
	// 							value:    " headline",
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}},
	// }

	// for _, _ := range testCases {
	// 	tree := Parse(tc.value)
	// t.Log(tree, "\nexpected:\n", tc.expected)
	// 	t.Log(*tree.Root.Nodes)
	// 	for _, node := range tree.Root.Nodes {
	// 		t.Log(node)
	// 	}
	// }
	// tree := Parse("* headline")

	// tree = Parse("*not a headline*")

	// tree := Parse("this is*not a headline*")

}
