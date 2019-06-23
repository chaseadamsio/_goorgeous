package parse

import (
	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/lex"
)

func (p *parser) makeTable(parent ast.Node, start, end int) {
	node := ast.NewTableNode(parent, p.items[start:end])
	parent.Append(node)

	p.parseTableRows(node, start, end)
}

func (p *parser) parseTableRows(node ast.Node, start, end int) {
	current := start
	rowStart := start
	inRuleRow := false

	for current <= end {
		if p.items[current].IsPipe() && p.items[current+1].IsDash() {
			inRuleRow = true
		}

		if p.items[current].IsNewline() {
			tr := ast.NewTableRowNode(node, p.items[rowStart:current])
			node.Append(tr)

			if inRuleRow {
				tr.NodeType = "TableRule"
				inRuleRow = false
			} else {
				p.parseTableCells(tr, rowStart, current)
			}

			rowStart = current + 1
		}

		if inRuleRow &&
			!(p.items[current].IsDash() || p.items[current].IsPlus() || p.items[current].IsPipe()) {
			inRuleRow = false
		}

		current++
	}
}

func (p *parser) parseTableCells(node ast.Node, start, end int) {
	current := start + 1 // discard the firt pipe
	cellStart := start
	for current <= end {
		if p.items[current].IsPipe() {
			tc := ast.NewTableCellNode(node, p.items[cellStart+1:current])
			node.Append(tc)
			p.walkElements(tc, cellStart+1, current)
			cellStart = current
		}

		current++

	}
}

func (p *parser) matchesTable(current int) (found bool, end int) {
	itemsLength := len(p.items)
	token := p.items[current]

	if !token.IsPipe() {
		return false, -1
	}
	if current < itemsLength && current == 0 || p.items[current-1].IsNewline() {
		for current < itemsLength {
			token := p.items[current]
			if token.IsNewline() {
				if current < itemsLength && (p.items[current+1].Type() != lex.ItemPipe) {
					if p.items[current+1].IsEOF() {
						current++
						continue
					}
					break
				}
			}
			current++
		}
		return true, current
	}
	return false, -1
}
