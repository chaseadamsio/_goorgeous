package parse

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/testdata"
)

func TestParse(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// if !strings.HasPrefix(tc.name, "unordered-list-with-child-ordered-list-with-unordered-list-child") {
			// 	return
			// }
			value := testdata.GetOrgStr(tc.source)
			// _ = Parse(value)
			ast := Parse(value)
			var expected []interface{}
			err := json.Unmarshal([]byte(tc.expected), &expected)
			if err != nil {
				t.Errorf("failed to unmarshal expected string: %s", err)
				// use an empty string to get a view of the world and uncomment this next line...
				// fmt.Printf("\nexpected:\n\t%v\nactual:\n\t%v", expected, ast)
			}

			if reflect.DeepEqual(ast, expected) {
				t.Errorf("expected %s AST shape to match expected.", tc.name)
			}
		})
	}

}

type testCase struct {
	name     string
	source   string
	expected string
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
	// {
	// 	"headers",
	// 	"#+title: headers\n#+author: Chase Adams\n#+description: This is my description!",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Section",
	// 			[]testNode{{
	// 				"Keyword",
	// 				nil,
	// 			}, {
	// 				"Keyword",
	// 				nil,
	// 			}, {
	// 				"Keyword",
	// 				nil,
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"basic-happy-path-new-content-after",
	// 	"#+title: basic-happy-path-new-content-after\n#+author: Chase Adams\n#+description: This is my description!\n* This starts the content!",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"basic-happy-path-with-tags",
	// 	"#+title: basic-happy-path-with-tags\n#+author: Chase Adams\n#+description: This is my description!\n#+tags: org-content org-mode hugo\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },

	// {
	// 	"basic-happy-path-with-categories",
	// 	"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+categories: org-content org-mode hugo\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"basic-happy-path-with-aliases",
	// 	"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+aliases: /org/content /org/mode /hugo\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"basic-text",
	// 	"this is a line.\nthis is a newline.",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Section",
	// 			[]testNode{{
	// 				"Text",
	// 				nil,
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"headline - level 1",
	// 	"* this is a headline",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"headline - level 1 w/ newline",
	// 	"* this is a headline\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"headline - deep",
	// 	"* headline1\n** headline2\n*** headline3\n* headline1-2\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"headline with paragraph children - deep",
	// 	"* headline1\nthis is a line.\nthis is another line.\n** headline2\n*** headline3\n* headline1-2\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"link",
	// 	"[[https://github.com][this is a link with *some bold text*.]]",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"link w/ newline",
	// 	"[[https://github.com][this is a link]]\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"link-self-describing",
	// 	"[[https://github.com]]\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"complex",
	// 	"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"complex w/ newline",
	// 	"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"complex w/ trailing text",
	// 	"** hello\nthis is some text\n#+BEGIN_SRC javascript\nconsole.log(\"hello world\");\n#+END_SRC\nhello",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	{
		"unordered-list",
		testdata.UnorderedListBasic,
		` [
			{
			  "NodeType": "Root",
			  "Children": [
				{
				  "NodeType": "List",
				  "Children": [
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- apples"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- bananas"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- oranges"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- pears"
					}
				  ],
				  "Value": "- apples\n- bananas\n- oranges\n- pears\n"
				}
			  ],
			  "Value": ""
			}
		  ]`,
	},
	{
		"unordered-list-with-child-ordered-list",
		testdata.UnorderedListWithNestedOrderedList,
		`[
			{
			  "NodeType": "Root",
			  "Children": [
				{
				  "NodeType": "List",
				  "Children": [
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- apples"
					},
					{
					  "NodeType": "ListItem",
					  "Children": [
						{
						  "NodeType": "List",
						  "Children": [
							{
							  "NodeType": "ListItem",
							  "Children": null,
							  "Value": "1. more apples"
							},
							{
							  "NodeType": "ListItem",
							  "Children": null,
							  "Value": "2. more bananas"
							},
							{
							  "NodeType": "ListItem",
							  "Children": null,
							  "Value": "3. more oranges"
							},
							{
							  "NodeType": "ListItem",
							  "Children": null,
							  "Value": "4. more pears"
							}
						  ],
						  "Value": "    1. more apples\n    2. more bananas\n    3. more oranges\n    4. more pears\n"
						}
					  ],
					  "Value": "- bananas"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- oranges"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- pears"
					}
				  ],
				  "Value": "- apples\n- bananas\n    1. more apples\n    2. more bananas\n    3. more oranges\n    4. more pears\n- oranges\n- pears\n"
				},
				{
				  "NodeType": "Section",
				  "Children": [
					{
					  "NodeType": "Paragraph",
					  "Children": [
						{
						  "NodeType": "Text",
						  "Children": null,
						  "Value": "with an extra paragraph\n"
						}
					  ],
					  "Value": "with an extra paragraph\n"
					}
				  ],
				  "Value": ""
				}
			  ],
			  "Value": ""
			}
		  ]`,
	},
	// {
	// 	"unordered-list-with-child-ordered-list-with-ordered-list-child",
	// 	"- apples\n\t1. in apples 1\n\t2. in apples 2\n\t\t1. in apples 1\n\t\t2. in apples 2\n\t\t3. in apples 3\n\t3. in apples 3\n- oranges\n- bananas\nsomething else",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	{
		"unordered-list-with-child-ordered-list-with-unordered-list-child",
		testdata.UnorderedListWithDeepNestedChildren,
		`[
			{
			  "NodeType": "Root",
			  "Children": [
				{
				  "NodeType": "List",
				  "Children": [
					{
					  "NodeType": "ListItem",
					  "Children": [
						{
						  "NodeType": "List",
						  "Children": [
							{
							  "NodeType": "ListItem",
							  "Children": null,
							  "Value": "1. in apples 1"
							},
							{
							  "NodeType": "ListItem",
							  "Children": [
								{
								  "NodeType": "List",
								  "Children": [
									{
									  "NodeType": "ListItem",
									  "Children": null,
									  "Value": "- in apples 1"
									},
									{
									  "NodeType": "ListItem",
									  "Children": null,
									  "Value": "- in apples 2"
									},
									{
									  "NodeType": "ListItem",
									  "Children": null,
									  "Value": "- in apples 3"
									}
								  ],
								  "Value": "    - in apples 1\n    - in apples 2\n    - in apples 3\n"
								}
							  ],
							  "Value": "2. in apples 2"
							},
							{
							  "NodeType": "ListItem",
							  "Children": null,
							  "Value": "3. in apples 3"
							}
						  ],
						  "Value": "  1. in apples 1\n  2. in apples 2\n    - in apples 1\n    - in apples 2\n    - in apples 3\n  3. in apples 3\n"
						}
					  ],
					  "Value": "- apples"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- oranges"
					},
					{
					  "NodeType": "ListItem",
					  "Children": null,
					  "Value": "- bananas"
					}
				  ],
				  "Value": "- apples\n  1. in apples 1\n  2. in apples 2\n    - in apples 1\n    - in apples 2\n    - in apples 3\n  3. in apples 3\n- oranges\n- bananas\n"
				},
				{
				  "NodeType": "Section",
				  "Children": [
					{
					  "NodeType": "Paragraph",
					  "Children": [
						{
						  "NodeType": "Text",
						  "Children": null,
						  "Value": "something else\n"
						}
					  ],
					  "Value": "something else\n"
					}
				  ],
				  "Value": ""
				}
			  ],
			  "Value": ""
			}
		  ]`,
	},
	// {
	// 	"ordered-list-with-unordered-list-child",
	// 	"1. apples\n2. oranges\n\t- apples\n\t- oranges\n\t- bananas\n\t\tsomething else\n3. bananas\nsomething else",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"table",
	// 	"| Name  | Phone | Age |\n|-------+-------+-----|\n| Peter |  1234 |  17 |\n| Anna  |  4321 |  25 |\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"footnote-number",
	// 	"The Org homepage[fn:1] now looks a lot better than it used to.\n[fn:1] The link is: https://orgmode.org",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"footnote-anonymous-inline-definition",
	// 	"The Org homepage[fn::This is the inline definition of this footnote] now looks a lot better than it used to.\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"footnote-inline-description",
	// 	"The Org homepage[fn:name:a definition]	now looks a lot better than it used to.\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"bold-with-italic-child",
	// 	" *this is some /italic text/ in a bold element.*\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
	// {
	// 	"bold-with-italic-child-with-verbatim-child",
	// 	" *this is some /italic text with =a verbatim child=/ in a bold element.*\n",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	// 			"Headline",
	// 			[]testNode{{
	// 				"Headline",
	// 				[]testNode{{
	// 					"Headline",
	// 					nil,
	// 				}},
	// 			}},
	// 		}},
	// 	}},
	// },
}
