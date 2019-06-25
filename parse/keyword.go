package parse

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeKeyword(parent ast.Node, start, end int) {

	node := ast.NewKeywordNode(parent, p.items[start:end])

	p.parseKeyword(node, start, end)

	parent.Append(node)
}

func (p *parser) parseKeyword(node *ast.KeywordNode, start, end int) {
	var value []string
	current := start
	foundKeywordIdentifier := false
	foundColon := false
	for current < end {
		if !foundKeywordIdentifier && p.items[current].IsHash() && p.items[current+1].IsPlus() {
			foundKeywordIdentifier = true
			current = current + 2
			start = current
			continue
		}

		if foundColon {
			value = append(value, p.items[current].Value())
		}

		if p.items[current].IsColon() && !foundColon {
			keyStr := ""
			keyCurr := start
			for keyCurr < current {
				keyStr = keyStr + p.items[keyCurr].Value()
				keyCurr++
			}
			node.Key = strings.ToUpper(keyStr)

			foundColon = true
		}

		current++
	}

	node.Value = strings.Join(value, "")

}

func (p *parser) matchesKeyword(current int) (found bool, end int) {
	token := p.items[current]
	itemsLength := len(p.items)

	if !token.IsHash() {
		return false, -1
	}
	if current < itemsLength && p.items[current+1].Type() != lex.ItemPlus {
		return false, -1
	}

	for reverseSearch := current - 1; 0 <= reverseSearch; {
		if p.items[reverseSearch].IsSpace() || p.items[reverseSearch].IsTab() {
			reverseSearch--
			continue
		} else if p.items[reverseSearch].IsNewline() {
			break
		}
		return false, -1
	}

	for current < itemsLength {
		if p.items[current].IsNewline() {
			return true, current
		}
		current++
	}
	return false, -1
}
