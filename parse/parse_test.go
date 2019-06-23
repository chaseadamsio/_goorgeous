package parse

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/chaseadamsio/goorgeous/ast"
	"github.com/chaseadamsio/goorgeous/testdata"
	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update golden files")

func TestParse(t *testing.T) {
	for _, tc := range tests {
		// filter := "headline/headline-1.org"
		// if !strings.HasPrefix(tc.source, filter) {
		// 	continue
		// }

		t.Run(tc.source, func(t *testing.T) {

			value := testdata.GetOrgStr(tc.source)
			ast := Parse(value)

			if *update {
				out := fmt.Sprintf("%v", ast)
				err := os.MkdirAll(filepath.Dir(tc.golden), os.ModePerm)
				if err != nil {
					t.Errorf("failed to make directories for %s: %s", tc.golden, err)
				}
				if err := ioutil.WriteFile(tc.golden, []byte(out), os.ModePerm); err != nil {
					t.Errorf("failed to write %s file: %s", tc.golden, err)
				}
				return
			}

			gldn, err := ioutil.ReadFile(tc.golden)
			if err != nil {
				t.Fatalf("failed to read %s file: %s", tc.golden, err)
			}

			var expected map[string]interface{}
			err = json.Unmarshal([]byte(gldn), &expected)
			if err != nil {
				t.Errorf("failed to unmarshal expected string: %s", err)
				// use an empty string to get a view of the world and uncomment this next line...
				// fmt.Printf("\nexpected:\n\t%v\nactual:\n\t%v", expected, ast)
			}

			if cmp.Equal(ast, expected) {
				t.Errorf("expected %s AST shape to match expected.", tc.source)
			}
		})
	}

}

type testCase struct {
	source string
	golden string
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
		testdata.ElementPlain,
		fmt.Sprintf("testdata/%s.json", testdata.ElementPlain),
	},
	{
		testdata.ElementNested,
		fmt.Sprintf("testdata/%s.json", testdata.ElementNested),
	},
	{
		testdata.ElementBold,
		fmt.Sprintf("testdata/%s.json", testdata.ElementBold),
	},
	{
		testdata.ElementHorizontalRule,
		fmt.Sprintf("testdata/%s.json", testdata.ElementHorizontalRule),
	},
	{
		testdata.Headline1,
		fmt.Sprintf("testdata/%s.json", testdata.Headline1),
	},
	{
		testdata.Headline1And2,
		fmt.Sprintf("testdata/%s.json", testdata.Headline1And2),
	},
	{
		testdata.Headline1WithContent,
		fmt.Sprintf("testdata/%s.json", testdata.Headline1WithContent),
	},
	{
		testdata.HeadersBasic,
		fmt.Sprintf("testdata/%s.json", testdata.HeadersBasic),
	},
	{
		testdata.LinkStandard,
		fmt.Sprintf("testdata/%s.json", testdata.LinkStandard),
	},
	{
		testdata.LinkSelfDescriptive,
		fmt.Sprintf("testdata/%s.json", testdata.LinkSelfDescriptive),
	},
	{
		testdata.LinkBoth,
		fmt.Sprintf("testdata/%s.json", testdata.LinkBoth),
	},
	{
		testdata.OrderedListBasic,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListBasic),
	},
	{
		testdata.OrderedListNotAList,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListNotAList),
	},
	{
		testdata.OrderedListWithStartingNewline,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListWithStartingNewline),
	},
	{
		testdata.OrderedListFollowParagraph,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListFollowParagraph),
	},
	{
		testdata.OrderedListFollowNumberNotList,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListFollowNumberNotList),
	},
	{
		testdata.OrderedListFollowAsteriskHeading,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListFollowAsteriskHeading),
	},
	{
		testdata.OrderedListWithFollowUnOrderedList,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListWithFollowUnOrderedList),
	},
	{
		testdata.OrderedListWithNestedOrderedList,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListWithNestedOrderedList),
	},
	{
		testdata.OrderedListWithNestedUnorderedList,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListWithNestedUnorderedList),
	},
	{
		testdata.OrderedListWithDeepNestedChildren,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListWithDeepNestedChildren),
	},
	{
		testdata.OrderedListWithNestedContent,
		fmt.Sprintf("testdata/%s.json", testdata.OrderedListWithNestedContent),
	},
	{
		testdata.UnorderedListBasic,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListBasic),
	},
	{
		testdata.UnorderedListNotAList,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListNotAList),
	},
	{
		testdata.UnorderedListWithStartingNewline,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListWithStartingNewline),
	},
	{
		testdata.UnorderedListFollowParagraph,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListFollowParagraph),
	},
	{
		testdata.UnorderedListFollowDashNotList,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListFollowDashNotList),
	},
	{
		testdata.UnorderedListFollowAsteriskHeading,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListFollowAsteriskHeading),
	},
	{
		testdata.UnorderedListWithFollowOrderedList,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListWithFollowOrderedList),
	},
	{
		testdata.UnorderedListWithNestedOrderedList,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListWithNestedOrderedList),
	},
	{
		testdata.UnorderedListWithNestedUnorderedList,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListWithNestedUnorderedList),
	},
	{
		testdata.UnorderedListWithDeepNestedChildren,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListWithDeepNestedChildren),
	},
	{
		testdata.UnorderedListWithNestedContent,
		fmt.Sprintf("testdata/%s.json", testdata.UnorderedListWithNestedContent),
	},
	{
		testdata.TableBasic,
		fmt.Sprintf("testdata/%s.json", testdata.TableBasic),
	},
	// {
	// 	"headers",
	// 	"#+title: headers\n#+author: Chase Adams\n#+description: This is my description!",
	// 	[]testNode{{
	// 		"Root",
	// 		[]testNode{{
	//			"Section",
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
