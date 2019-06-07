package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chaseadamsio/goorgeous/ast"
)

func TestParse(t *testing.T) {
	for _, tc := range tests {
		if !strings.HasPrefix(tc.name, "headline - deep") {
			continue
		}
		fmt.Println(tc.name)
		t.Run(tc.name, func(t *testing.T) {
			ast := Parse(tc.input)
			if true {
				t.Error(ast)
			}
		})
	}
}

func getPathForNode(n ast.Node, path []ast.NodeType) []ast.NodeType {
	path = append(path, n.Type())
	if n.Parent().Type() != "Root" {
		path = getPathForNode(n.Parent(), path)
	} else {
		var reversedPath []ast.NodeType
		path = append(path, "Root")
		for idx := len(path) - 1; idx >= 0; idx-- {
			reversedPath = append(reversedPath, path[idx])
		}
		return reversedPath
	}
	return path
}

type testCase struct {
	name     string
	input    string
	expected []testNode
}

type testNode struct {
	NodeType string
	children childrenTestNodes
}

type childrenTestNodes []testNode

var tests = []testCase{
	{
		"headers",
		"#+title: my org mode content\n#+author: Chase Adams\n#+description: This is my description!",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"basic-happy-path-new-content-after",
		"#+title: my org mode content\n#+author: Chase Adams\n#+description: This is my description!\n* This starts the content!",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"basic-happy-path-with-tags",
		"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+tags: org-content org-mode hugo\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},

	{
		"basic-happy-path-with-categories",
		"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+categories: org-content org-mode hugo\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"basic-happy-path-with-aliases",
		"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+aliases: /org/content /org/mode /hugo\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"basic - text",
		"this is a line.\nthis is a newline.",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"headline - level 1",
		"* this is a headline",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"headline - level 1 w/ newline",
		"* this is a headline\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"headline - deep",
		"* headline1\n** headline2\n*** headline3\n* headline1-2\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"link",
		"[this is a link](https://github.com)",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"link w/ newline",
		"[this is a link](https://github.com)\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"complex",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"complex w/ newline",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"complex w/ trailing text",
		"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\nhello",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"list",
		"- apples\n- oranges\n- bananas\nsomething else",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
	{
		"table",
		"| Name  | Phone | Age |\n|-------+-------+-----|\n| Peter |  1234 |  17 |\n| Anna  |  4321 |  25 |\n",
		[]testNode{{
			"Root",
			[]testNode{{
				"Headline",
				[]testNode{{
					"Headline",
					[]testNode{{
						"Headline",
						nil,
					}},
				}},
			}},
		}},
	},
}
