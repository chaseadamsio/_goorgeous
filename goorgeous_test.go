package goorgeous

import (
	"flag"
	"fmt"
)

var update = flag.Bool("update", false, "update golden files")

var builtins = map[string]interface{}{
	"printf": fmt.Sprintf,
}

//func TestTreeParse(t *testing.T) {
// 	tr := New("testTree")
// 	tr.text = "*** A H3 Headline"
// 	tr.startParse(nil, lex("test", tr.text), make(map[string]*Tree))
// 	out := tr.parse()
// 	t.Logf("%s", out)
// }
