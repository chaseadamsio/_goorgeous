package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeGreaterBlock(parent ast.Node, items []lex.Item) {
	node := ast.NewGreaterBlockNode(parent, items)

	parent.Append(node)
	var key string
	var val []string
	for idx, item := range items {
		if item.Type() == lex.ItemColon {
			key = items[idx-1].Value()
			continue
		} else if key != "" {
			val = append(val, item.Value())
		}
	}
	node.Key = key
	node.Value = strings.Join(val, "")
}
