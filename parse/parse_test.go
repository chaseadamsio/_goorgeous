package parse

import (
	"strings"
	"testing"

	"github.com/chaseadamsio/goorgeous/ast"
)

func TestParse(t *testing.T) {
	for _, tc := range tests {
		if !strings.HasPrefix(tc.name, "link") {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			ast := Parse(tc.input)
			if true {
				t.Errorf("\nname: %s\n\tinput: %s\n\t%v", tc.name, tc.input, ast)
			}
		})
	}
}

type testCase struct {
	name     string
	input    string
	expected []testNode
}

type testNode struct {
	ast.NodeType
	// value    string
	children childrenTestNodes
}

type childrenTestNodes []testNode

func (n *testNode) Type() ast.NodeType {
	return n.NodeType
}

// func (n *testNode) Parent() ast.Node {
// 	return n.parent
// }

// func (n *testNode) String() string {
// 	return n.value
// }

// func (n *testNode) Children() []ast.Node {
// 	return n.children
// }

// func (n *testNode) Append(child ast.Node) {
// 	n.children = append(n.children, child)
// }

var tests = []testCase{
	{
		"headers",
		"#+title: headers\n#+author: Chase Adams\n#+description: This is my description!",
		[]testNode{{
			"Root",
			[]testNode{{
				"Section",
				[]testNode{{
					"Keyword",
					nil,
				}, {
					"Keyword",
					nil,
				}, {
					"Keyword",
					nil,
				}},
			}},
		}},
	},
	{
		"basic-happy-path-new-content-after",
		"#+title: basic-happy-path-new-content-after\n#+author: Chase Adams\n#+description: This is my description!\n* This starts the content!",
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
		"#+title: basic-happy-path-with-tags\n#+author: Chase Adams\n#+description: This is my description!\n#+tags: org-content org-mode hugo\n",
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
		"basic-text",
		"this is a line.\nthis is a newline.",
		[]testNode{{
			"Root",
			[]testNode{{
				"Section",
				[]testNode{{
					"Text",
					nil,
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
		"headline with paragraph children - deep",
		"* headline1\nthis is a line.\nthis is another line.\n** headline2\n*** headline3\n* headline1-2\n",
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
		"[[https://github.com][this is a link with *some bold text*.]]",
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
		"[[https://github.com][this is a link]]\n",
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
		"link-self-describing",
		"[[https://github.com]]\n",
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
	{
		"footnote-number",
		"The Org homepage[fn:1] now looks a lot better than it used to.\n[fn:1] The link is: https://orgmode.org",
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
		"footnote-anonymous-inline-definition",
		"The Org homepage[fn::This is the inline definition of this footnote] now looks a lot better than it used to.\n",
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
		"footnote-inline-description",
		"The Org homepage[fn:name:a definition]	now looks a lot better than it used to.\n",
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
		"bold-with-italic-child",
		" *this is some /italic text/ in a bold element.*\n",
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
		"bold-with-italic-child-with-verbatim-child",
		" *this is some /italic text with =a verbatim child=/ in a bold element.*\n",
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
