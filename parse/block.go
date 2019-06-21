package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeGreaterBlock(parent ast.Node, start, end int) {
	node := ast.NewGreaterBlockNode(parent, p.items[start:end])

	parent.Append(node)
	var key string
	var val []string
	for idx, item := range p.items[start:end] {
		if item.Type() == lex.ItemColon {
			key = p.items[idx-1].Value()
			continue
		} else if key != "" {
			val = append(val, item.Value())
		}
	}
	node.Key = key
	node.Value = strings.Join(val, "")
}
