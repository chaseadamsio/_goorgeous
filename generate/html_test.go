package generator

import (
	"testing"

	"github.com/chaseadamsio/goorgeous/parse"
	"github.com/chaseadamsio/goorgeous/transform"
)

func TestGenerator(t *testing.T) {
	inAST := parse.Parse("- apples\n- oranges\n- bananas\nsomething else")
	outAST := transform.TransformToHTML(inAST)
	_ = Generate(outAST)
	// html := Generate(outAST)
	// t.Errorf(html)
}
