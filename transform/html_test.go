package transform

import (
	"fmt"
	"testing"

	"github.com/chaseadamsio/goorgeous/parse"
)

type testCase struct {
	input string
	// expected string
}

var testCases = []testCase{
	{
		"* This is a headline",
	},
}

func TestTransform(t *testing.T) {
	for _, tc := range testCases {
		ast := parse.Parse(tc.input)
		fmt.Println(ast)
		// TransformToHTML(ast)
		// flatOut := []string{"\n"}

		// for _, node := range out {
		// 	flatOut = append(flatOut, "type: "+string(node.Type)+"\n\t"+node.Value+"\n")
		// }
		// fmt.Println(strings.Join(flatOut, ""))
	}
}
