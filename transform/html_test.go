package transform

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/parse"
)

type testCase struct {
	input string
	// expected string
}

var testCases = []testCase{
	{
		"#+title: basic-happy-path-new-content-after\n#+author: Chase Adams\n#+description: This is my description!\n* This starts the content!",
	},
	{
		"- apples\n- oranges\n- bananas\nsomething else",
	},
}

func TestTransform(t *testing.T) {
	for _, tc := range testCases {
		ast := parse.Parse(tc.input)
		out := TransformToHTML(ast)
		flatOut := []string{"\n"}

		for _, node := range out {
			flatOut = append(flatOut, "type: "+string(node.Type)+"\n\t"+node.Value+"\n")
		}
		// fmt.Println(strings.Join(flatOut, ""))
	}
}
