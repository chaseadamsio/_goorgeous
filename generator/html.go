package generator

import (
	"strings"

	"github.com/chaseadamsio/goorgeous/transform"
)

func Generate(htmlNodes []transform.HTMLNode) string {
	var html []string
	for _, node := range htmlNodes {
		if node.Type == "List" {
			html = append(html, "<ul></ul>")
		}
	}
	return strings.Join(html, "\n")
}
