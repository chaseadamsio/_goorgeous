package goorgeous

import "testing"

type lexTest struct {
	name  string
	input string
	items []item
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

var (
	tEOF     = mkItem(itemEOF, "")
	tNewLine = mkItem(itemNewLine, "\n")

	tH1 = mkItem(itemHeadline, delimH1)
	tH2 = mkItem(itemHeadline, delimH2)
	tH3 = mkItem(itemHeadline, delimH3)
	tH4 = mkItem(itemHeadline, delimH4)
	tH5 = mkItem(itemHeadline, delimH5)
	tH6 = mkItem(itemHeadline, delimH6)

	tEmphasis      = mkItem(itemEmphasis, "/")
	tBold          = mkItem(itemBold, "*")
	tStrikethrough = mkItem(itemStrikethrough, "+")
	tVerbatim      = mkItem(itemVerbatim, "=")
	tCode          = mkItem(itemCode, "~")
	tUnderline     = mkItem(itemUnderline, "_")

	tStatusTODO = mkItem(itemStatus, "TODO")
	tStatusDONE = mkItem(itemStatus, "DONE")
	tPriorityA  = mkItem(itemPriority, "[A]")
	tTag        = mkItem(itemTags, ":")

	tTable = mkItem(itemTable, "|")

	tImgOrLinkOpen        = mkItem(itemImgOrLinkOpen, "[[")
	tImgOrLinkOpenSingle  = mkItem(itemImgOrLinkOpenSingle, "[")
	tImgPre               = mkItem(itemImgPre, "file:")
	tImgOrLinkClose       = mkItem(itemImgOrLinkClose, "]]")
	tImgOrLinkCloseSingle = mkItem(itemImgOrLinkCloseSingle, "]")
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"spaces", " \t\n", []item{mkItem(itemText, " \t"), tNewLine, tEOF}},
	{"text", "now is the time", []item{mkItem(itemText, "now is the time"), tEOF}},
	{"text", "now is the time\n", []item{mkItem(itemText, "now is the time"), tNewLine, tEOF}},
	{"text", "now is the time\n\n", []item{mkItem(itemText, "now is the time"), tNewLine, tNewLine, tEOF}},

	// BASIC HEADLINES
	{"h1", "* A h1 Headline", []item{tH1, mkItem(itemText, "A h1 Headline"), tEOF}},
	{"not-h1", "not an * h1 Headline", []item{mkItem(itemText, "not an * h1 Headline"), tEOF}},
	{"alt-not-h1", " * not an h1 Headline", []item{mkItem(itemText, " * not an h1 Headline"), tEOF}},
	{"alt-not-h1-2", "*not an h1 Headline", []item{mkItem(itemText, "*not an h1 Headline"), tEOF}},
	{"h1-with-status-todo", "* TODO a h1 headline", []item{tH1, tStatusTODO, mkItem(itemText, "a h1 headline"), tEOF}},
	{"h1-with-status-done", "* DONE a h1 headline", []item{tH1, tStatusDONE, mkItem(itemText, "a h1 headline"), tEOF}},
	{"h1-with-priority-a", "* [A] a h1 headline", []item{tH1, tPriorityA, mkItem(itemText, "a h1 headline"), tEOF}},
	{"h1-not-priority", "* [Z] a h1 headline", []item{tH1, mkItem(itemText, "[Z] a h1 headline"), tEOF}},
	{"h1-with-single-tag", "* A h1 Headline :singletag:", []item{tH1, mkItem(itemText, "A h1 Headline "), tTag, mkItem(itemText, "singletag"), tTag, tEOF}},
	{"h1-with-status-todo-priority-a", "* TODO [A] a h1 headline", []item{tH1, tStatusTODO, tPriorityA, mkItem(itemText, "a h1 headline"), tEOF}},
	{"h1-with-status-todo-not-priority", "* TODO [Z] a h1 headline", []item{tH1, tStatusTODO, mkItem(itemText, "[Z] a h1 headline"), tEOF}},
	{"h2", "** A h2 Headline", []item{tH2, mkItem(itemText, "A h2 Headline"), tEOF}},
	{"not-h2", "not an ** h2 Headline", []item{mkItem(itemText, "not an ** h2 Headline"), tEOF}},
	{"alt-not-h2", " ** not an h2 Headline", []item{mkItem(itemText, " ** not an h2 Headline"), tEOF}},
	{"alt-not-h2-2", "**not an h2 Headline", []item{mkItem(itemText, "**not an h2 Headline"), tEOF}},
	{"h3", "*** A h3 Headline", []item{tH3, mkItem(itemText, "A h3 Headline"), tEOF}},
	{"not-h3", "not an *** h3 Headline", []item{mkItem(itemText, "not an *** h3 Headline"), tEOF}},
	{"alt-not-h3", " *** not an h3 Headline", []item{mkItem(itemText, " *** not an h3 Headline"), tEOF}},
	{"alt-not-h3-2", "***not an h3 Headline", []item{mkItem(itemText, "***not an h3 Headline"), tEOF}},
	{"h4", "**** A h4 Headline", []item{tH4, mkItem(itemText, "A h4 Headline"), tEOF}},
	{"not-h4", "not an **** h4 Headline", []item{mkItem(itemText, "not an **** h4 Headline"), tEOF}},
	{"alt-not-h4", " **** not an h4 Headline", []item{mkItem(itemText, " **** not an h4 Headline"), tEOF}},
	{"alt-not-h4-2", "****not an h4 Headline", []item{mkItem(itemText, "****not an h4 Headline"), tEOF}},
	{"h5", "***** A h5 Headline", []item{tH5, mkItem(itemText, "A h5 Headline"), tEOF}},
	{"not-h5", "not an ***** h5 Headline", []item{mkItem(itemText, "not an ***** h5 Headline"), tEOF}},
	{"alt-not-h5", " ***** not an h5 Headline", []item{mkItem(itemText, " ***** not an h5 Headline"), tEOF}},
	{"alt-not-h5-2", "*****not an h5 Headline", []item{mkItem(itemText, "*****not an h5 Headline"), tEOF}},
	{"h6", "****** A h6 Headline", []item{tH6, mkItem(itemText, "A h6 Headline"), tEOF}},
	{"not-h6", "not an ****** h6 Headline", []item{mkItem(itemText, "not an ****** h6 Headline"), tEOF}},
	{"alt-not-h6", " ****** not an h6 Headline", []item{mkItem(itemText, " ****** not an h6 Headline"), tEOF}},
	{"alt-not-h6-2", "******not an h6 Headline", []item{mkItem(itemText, "******not an h6 Headline"), tEOF}},

	// COMPLEX
	{"complex-h1-emphasis-h6", "* A h1 Headline\nsome /emphasis text/.\n****** A h6 Headline",
		[]item{
			tH1,
			mkItem(itemText, "A h1 Headline"),
			tNewLine,
			mkItem(itemText, "some "),
			tEmphasis,
			mkItem(itemText, "emphasis text"),
			tEmphasis,
			mkItem(itemText, "."),
			tNewLine,
			tH6,
			mkItem(itemText, "A h6 Headline"),
			tEOF}},
	{"complex-h1-emphasis-h6", "* A h1 Headline\nsome /*emphasis* text/.\n****** A h6 Headline",
		[]item{
			tH1,
			mkItem(itemText, "A h1 Headline"),
			tNewLine,
			mkItem(itemText, "some "),
			tEmphasis,
			tBold,
			mkItem(itemText, "emphasis"),
			tBold,
			mkItem(itemText, " text"),
			tEmphasis,
			mkItem(itemText, "."),
			tNewLine,
			tH6,
			mkItem(itemText, "A h6 Headline"),
			tEOF}},

	// PROPERTY DRAWER
	{"property-drawer", ":PROPERTIES:\n:Title:     Goldberg Variations\n:END:",
		[]item{
			mkItem(itemPropertyDrawer, ":PROPERTIES:"),
			tNewLine,
			mkItem(itemText, ":Title:     Goldberg Variations"),
			tNewLine,
			mkItem(itemPropertyDrawer, ":END:"),
			tEOF,
		},
	},

	// BLOCKS
	{"src", "#+BEGIN_SRC sh\necho \"foo\"\n#+END_SRC",
		[]item{
			mkItem(itemBlock, "#+BEGIN_SRC"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			mkItem(itemBlock, "#+END_SRC"),
			tEOF,
		},
	},
	{"src-multiline", "#+BEGIN_SRC sh\necho \"foo\"\necho \"bar\"\n#+END_SRC",
		[]item{
			mkItem(itemBlock, "#+BEGIN_SRC"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			mkItem(itemBlock, "#+END_SRC"),
			tEOF,
		},
	},
	{"src-multiline-multi-newline", "#+BEGIN_SRC sh\necho \"foo\"\n\necho \"bar\"\n#+END_SRC",
		[]item{
			mkItem(itemBlock, "#+BEGIN_SRC"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			mkItem(itemBlock, "#+END_SRC"),
			tEOF,
		},
	},
	{"src-multiline-many-multi-newline", "#+BEGIN_SRC sh\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n#+END_SRC",
		[]item{
			mkItem(itemBlock, "#+BEGIN_SRC"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			mkItem(itemBlock, "#+END_SRC"),
			tEOF,
		},
	},
	{"src-multiline-many-multi-newline-text",
		"#+BEGIN_SRC text\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n#+END_SRC",
		[]item{
			mkItem(itemBlock, "#+BEGIN_SRC"),
			mkItem(itemText, " text"),
			tNewLine,
			mkItem(itemText, "/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo"),
			tNewLine,
			mkItem(itemText, "ligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque"),
			tNewLine,
			mkItem(itemText, "eu, sem. Nulla consequat massa quis enim./"),
			tNewLine,
			tNewLine,
			mkItem(itemText, "/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam"),
			tNewLine,
			mkItem(itemText, "dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus"),
			tNewLine,
			mkItem(itemText, "elementum semper nisi./"),
			tNewLine,
			mkItem(itemBlock, "#+END_SRC"),
			tEOF,
		},
	},
	{"example", "#+BEGIN_EXAMPLE sh\necho \"foo\"\n#+END_EXAMPLE",
		[]item{
			mkItem(itemBlock, "#+BEGIN_EXAMPLE"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			mkItem(itemBlock, "#+END_EXAMPLE"),
			tEOF,
		},
	},
	{"example-multiline", "#+BEGIN_EXAMPLE sh\necho \"foo\"\necho \"bar\"\n#+END_EXAMPLE",
		[]item{
			mkItem(itemBlock, "#+BEGIN_EXAMPLE"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			mkItem(itemBlock, "#+END_EXAMPLE"),
			tEOF,
		},
	},
	{"examplec-multiline-multi-newline", "#+BEGIN_EXAMPLE sh\necho \"foo\"\n\necho \"bar\"\n#+END_EXAMPLE",
		[]item{
			mkItem(itemBlock, "#+BEGIN_EXAMPLE"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			mkItem(itemBlock, "#+END_EXAMPLE"),
			tEOF,
		},
	},
	{"example-multiline-many-multi-newline", "#+BEGIN_EXAMPLE sh\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n#+END_EXAMPLE",
		[]item{
			mkItem(itemBlock, "#+BEGIN_EXAMPLE"),
			mkItem(itemText, " sh"),
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"foo\""),
			tNewLine,
			tNewLine,
			mkItem(itemText, "echo \"bar\""),
			tNewLine,
			mkItem(itemBlock, "#+END_EXAMPLE"),
			tEOF,
		},
	},

	{"quote", "#+BEGIN_QUOTE\nthis is a quote.\n#+END_QUOTE",
		[]item{
			mkItem(itemBlock, "#+BEGIN_QUOTE"),
			tNewLine,
			mkItem(itemText, "this is a quote."),
			tNewLine,
			mkItem(itemBlock, "#+END_QUOTE"),
			tEOF,
		},
	},
	{"quote-multiline", "#+BEGIN_QUOTE\nthis is a quote\nwith multiple lines.\n#+END_QUOTE",
		[]item{
			mkItem(itemBlock, "#+BEGIN_QUOTE"),
			tNewLine,
			mkItem(itemText, "this is a quote"),
			tNewLine,
			mkItem(itemText, "with multiple lines."),
			tNewLine,
			mkItem(itemBlock, "#+END_QUOTE"),
			tEOF,
		},
	},

	{"center", "#+BEGIN_CENTER\nthis is a centered block.\n#+END_CENTER",
		[]item{
			mkItem(itemBlock, "#+BEGIN_CENTER"),
			tNewLine,
			mkItem(itemText, "this is a centered block."),
			tNewLine,
			mkItem(itemBlock, "#+END_CENTER"),
			tEOF,
		},
	},
	{"center-multiline", "#+BEGIN_CENTER\nthis is a centered block\nwith multiple lines.\n#+END_CENTER",
		[]item{
			mkItem(itemBlock, "#+BEGIN_CENTER"),
			tNewLine,
			mkItem(itemText, "this is a centered block"),
			tNewLine,
			mkItem(itemText, "with multiple lines."),
			tNewLine,
			mkItem(itemBlock, "#+END_CENTER"),
			tEOF,
		},
	},

	// COMMENT
	{"comment-basic", "# this is a comment",
		[]item{
			mkItem(itemComment, "#"),
			mkItem(itemText, " this is a comment"),
			tEOF,
		},
	},
	{"comment-multiline", "# this is a comment\n# with multiple lines",
		[]item{
			mkItem(itemComment, "#"),
			mkItem(itemText, " this is a comment"),
			tNewLine,
			mkItem(itemComment, "#"),
			mkItem(itemText, " with multiple lines"),
			tEOF,
		},
	},
	{"comment-with-whitespace", "# this is a comment\n       # this is also a comment",
		[]item{
			mkItem(itemComment, "#"),
			mkItem(itemText, " this is a comment"),
			tNewLine,
			mkItem(itemComment, "       #"),
			mkItem(itemText, " this is also a comment"),
			tEOF,
		},
	},

	// HEADLINES VARIANTS
	{"multi-headlines", "* A h1 Headline\n****** A h6 Headline",
		[]item{
			tH1,
			mkItem(itemText, "A h1 Headline"),
			tNewLine,
			tH6,
			mkItem(itemText, "A h6 Headline"),
			tEOF}},
	{"h1-with-text", "* A h1 Headline\nThis is a new line.\n",
		[]item{
			tH1,
			mkItem(itemText, "A h1 Headline"),
			tNewLine,
			mkItem(itemText, "This is a new line."),
			tNewLine,
			tEOF}},

	// tables
	{"table-basic", "| Peter |  1234 |  17 |\n| Anna  |  4321 |  25 |",
		[]item{
			tTable, mkItem(itemText, " Peter "),
			tTable, mkItem(itemText, "  1234 "),
			tTable, mkItem(itemText, "  17 "),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, " Anna  "),
			tTable, mkItem(itemText, "  4321 "),
			tTable, mkItem(itemText, "  25 "),
			tTable,
			tEOF}},

	{"table-header", "| Name  | Phone | Age |\n|-------+-------+-----|\n| Peter |  1234 |  17 |\n| Anna  |  4321 |  25 |",
		[]item{
			tTable, mkItem(itemText, " Name  "),
			tTable, mkItem(itemText, " Phone "),
			tTable, mkItem(itemText, " Age "),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, "-------+-------+-----"),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, " Peter "),
			tTable, mkItem(itemText, "  1234 "),
			tTable, mkItem(itemText, "  17 "),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, " Anna  "),
			tTable, mkItem(itemText, "  4321 "),
			tTable, mkItem(itemText, "  25 "),
			tTable,
			tEOF}},

	{"table-header-horizontal-splits", "|---+---+---|\n| d | e | f |\n|---+---+---|\n| g | h | i |\n|---+---+---|",
		[]item{
			tTable, mkItem(itemText, "---+---+---"),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, " d "),
			tTable, mkItem(itemText, " e "),
			tTable, mkItem(itemText, " f "),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, "---+---+---"),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, " g "),
			tTable, mkItem(itemText, " h "),
			tTable, mkItem(itemText, " i "),
			tTable,
			tNewLine,
			tTable, mkItem(itemText, "---+---+---"),
			tTable,
			tEOF}},

	// list
	{"ordered-list-periods", "1. this\n2. is\n3. a list",
		[]item{
			mkItem(itemOrderedList, "1."),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemOrderedList, "2."),
			mkItem(itemText, " is"),
			tNewLine,
			mkItem(itemOrderedList, "3."),
			mkItem(itemText, " a list"),
			tEOF,
		},
	},
	{"not-full-ordered-list-periods", "1. this\nfoo 2. is\n3. a list",
		[]item{
			mkItem(itemOrderedList, "1."),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemText, "foo 2. is"),
			tNewLine,
			mkItem(itemOrderedList, "3."),
			mkItem(itemText, " a list"),
			tEOF,
		},
	},
	{"ordered-list-periods-change-number", "1. this\n2. is\n3. [@10] a list\n4. that jumps",
		[]item{
			mkItem(itemOrderedList, "1."),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemOrderedList, "2."),
			mkItem(itemText, " is"),
			tNewLine,
			mkItem(itemOrderedList, "3."),
			mkItem(itemOrderedListNumber, " [@10]"),
			mkItem(itemText, " a list"),
			tNewLine,
			mkItem(itemOrderedList, "4."),
			mkItem(itemText, " that jumps"),
			tEOF,
		},
	},
	{"not-ordered-list-periods-change-number", "1. this\n2. is\n3. [@10 ] a list\n4. that jumps",
		[]item{
			mkItem(itemOrderedList, "1."),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemOrderedList, "2."),
			mkItem(itemText, " is"),
			tNewLine,
			mkItem(itemOrderedList, "3."),
			mkItem(itemText, " [@10 ] a list"),
			tNewLine,
			mkItem(itemOrderedList, "4."),
			mkItem(itemText, " that jumps"),
			tEOF,
		},
	},

	{"ordered-list-parens", "1) this\n2) is\n3) a list",
		[]item{
			mkItem(itemOrderedList, "1)"),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemOrderedList, "2)"),
			mkItem(itemText, " is"),
			tNewLine,
			mkItem(itemOrderedList, "3)"),
			mkItem(itemText, " a list"),
			tEOF,
		},
	},

	{"unordered-list-dashes", "- this\n- is\n- a list",
		[]item{
			mkItem(itemUnorderedList, "-"),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemUnorderedList, "-"),
			mkItem(itemText, " is"),
			tNewLine,
			mkItem(itemUnorderedList, "-"),
			mkItem(itemText, " a list"),
			tEOF,
		},
	},
	{"not-full-ordered-list-dashes", "- this\nfoo - is\n- a list",
		[]item{
			mkItem(itemUnorderedList, "-"),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemText, "foo - is"),
			tNewLine,
			mkItem(itemUnorderedList, "-"),
			mkItem(itemText, " a list"),
			tEOF}},

	{"definition-list", "- definition lists :: these are useful sometimes\n- item 2 :: M-RET again gives another item, and long lines wrap in a tidy way underneath the definition",
		[]item{
			mkItem(itemDefinitionTerm, "-"),
			mkItem(itemText, " definition lists "),
			mkItem(itemDefinitionDescription, "::"),
			mkItem(itemText, " these are useful sometimes"),
			tNewLine,
			mkItem(itemDefinitionTerm, "-"),
			mkItem(itemText, " item 2 "),
			mkItem(itemDefinitionDescription, "::"),
			mkItem(itemText, " M-RET again gives another item, and long lines wrap in a tidy way underneath the definition"),
			tEOF,
		},
	},

	{"unordered-list-pluses", "+ this\n+ is\n+ a list",
		[]item{
			mkItem(itemUnorderedList, "+"),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemUnorderedList, "+"),
			mkItem(itemText, " is"),
			tNewLine,
			mkItem(itemUnorderedList, "+"),
			mkItem(itemText, " a list"),
			tEOF,
		},
	},
	{"not-full-ordered-list-pluses", "+ this\nfoo + is\n+ a list",
		[]item{
			mkItem(itemUnorderedList, "+"),
			mkItem(itemText, " this"),
			tNewLine,
			mkItem(itemText, "foo + is"),
			tNewLine,
			mkItem(itemUnorderedList, "+"),
			mkItem(itemText, " a list"),
			tEOF,
		},
	},
	// links & images
	{"link-basic", "this has [[https://github.com/chaseadamsio/goorgeous]] as a link.",
		[]item{
			mkItem(itemText, "this has "),
			tImgOrLinkOpen,
			mkItem(itemImgOrLinkURL, "https://github.com/chaseadamsio/goorgeous"),
			tImgOrLinkClose,
			mkItem(itemText, " as a link."),
			tEOF,
		},
	},

	{"link-with-alt", "this has [[https://github.com/chaseadamsio/goorgeous][goorgeous by chaseadamsio]] as a link.",
		[]item{
			mkItem(itemText, "this has "),
			tImgOrLinkOpen,
			mkItem(itemImgOrLinkURL, "https://github.com/chaseadamsio/goorgeous"),
			tImgOrLinkCloseSingle,
			tImgOrLinkOpenSingle,
			mkItem(itemImgOrLinkText, "goorgeous by chaseadamsio"),
			tImgOrLinkClose,
			mkItem(itemText, " as a link."),
			tEOF,
		},
	},

	{"not-link-with-alt", "this has [[https://github.com/chaseadamsio/goorgeous]foo[goorgeous by chaseadamsio]] as a link.",
		[]item{
			mkItem(itemText, "this has [[https://github.com/chaseadamsio/goorgeous]foo[goorgeous by chaseadamsio]] as a link."),
			tEOF,
		},
	},

	{"not-link", "this has [[https://github.com/chaseadamsio/goorgeous] as a link.",
		[]item{
			mkItem(itemText, "this has [[https://github.com/chaseadamsio/goorgeous] as a link."),
			tEOF,
		},
	},

	{"image", "this has [[file:https://github.com/chaseadamsio/goorgeous/img.png]] as an image.",
		[]item{
			mkItem(itemText, "this has "),
			tImgOrLinkOpen,
			tImgPre,
			mkItem(itemImgOrLinkURL, "https://github.com/chaseadamsio/goorgeous/img.png"),
			tImgOrLinkClose,
			mkItem(itemText, " as an image."),
			tEOF,
		},
	},

	{"image-with-alt", "this has [[file:https://github.com/chaseadamsio/goorgeous/img.png][a uni-gopher in the wild]] as an image.",
		[]item{
			mkItem(itemText, "this has "),
			tImgOrLinkOpen,
			tImgPre,
			mkItem(itemImgOrLinkURL, "https://github.com/chaseadamsio/goorgeous/img.png"),
			tImgOrLinkCloseSingle,
			tImgOrLinkOpenSingle,
			mkItem(itemImgOrLinkText, "a uni-gopher in the wild"),
			tImgOrLinkClose,
			mkItem(itemText, " as an image."),
			tEOF,
		},
	},

	{"not-image", "this has [[file:https://github.com/chaseadamsio/goorgeous/img.png] as an image.",
		[]item{
			mkItem(itemText, "this has [[file:https://github.com/chaseadamsio/goorgeous/img.png] as an image."),
			tEOF,
		},
	},

	{"not-image-with-alt", "this has [[file:https://github.com/chaseadamsio/goorgeous/img.png]foo[a uni-gopher in the wild]] as an image.",
		[]item{
			mkItem(itemText, "this has [[file:https://github.com/chaseadamsio/goorgeous/img.png]foo[a uni-gopher in the wild]] as an image."),
			tEOF,
		},
	},

	// emphasis
	{"emphasis", "/now is the time/", []item{tEmphasis, mkItem(itemText, "now is the time"), tEmphasis, tEOF}},
	{"emphasis-surrounded", "now is /the/ time", []item{
		mkItem(itemText, "now is "),
		tEmphasis,
		mkItem(itemText, "the"),
		tEmphasis,
		mkItem(itemText, " time"),
		tEOF}},
	{"emphasis-surrounded", "now is /the/ time\nthis is some more text!", []item{
		mkItem(itemText, "now is "),
		tEmphasis,
		mkItem(itemText, "the"),
		tEmphasis,
		mkItem(itemText, " time"),
		tNewLine,
		mkItem(itemText, "this is some more text!"),
		tEOF}},
	{"not-emphasis", "no/w is the time/", []item{mkItem(itemText, "no/w is the time/"), tEOF}},

	// bold
	{"bold", "*now is the time*", []item{tBold, mkItem(itemText, "now is the time"), tBold, tEOF}},
	{"bold-inside", "they say *now is the time*", []item{mkItem(itemText, "they say "), tBold, mkItem(itemText, "now is the time"), tBold, tEOF}},

	// strikethrough
	{"strikethrough", "+now is the time+", []item{tStrikethrough, mkItem(itemText, "now is the time"), tStrikethrough, tEOF}},
	{"strikethrough-surrounded", "now is +the+ time", []item{
		mkItem(itemText, "now is "),
		tStrikethrough,
		mkItem(itemText, "the"),
		tStrikethrough,
		mkItem(itemText, " time"),
		tEOF}},
	{"strikethrough-surrounded", "now is +the+ time\nthis is some more text!", []item{
		mkItem(itemText, "now is "),
		tStrikethrough,
		mkItem(itemText, "the"),
		tStrikethrough,
		mkItem(itemText, " time"),
		tNewLine,
		mkItem(itemText, "this is some more text!"),
		tEOF}},
	{"not-strikethrough", "no+w is the time+", []item{mkItem(itemText, "no+w is the time+"), tEOF}},

	// verbatim
	{"verbatim", "=simple verbatim=", []item{tVerbatim, mkItem(itemText, "simple verbatim"), tVerbatim, tEOF}},
	{"verbatim", "=simple=verbatim=", []item{tVerbatim, mkItem(itemText, "simple=verbatim"), tVerbatim, tEOF}},
	{"verbatim", "==simple=verbatim==", []item{tVerbatim, mkItem(itemText, "=simple=verbatim="), tVerbatim, tEOF}},
	{"verbatim-surrounded", "now is =the= time", []item{
		mkItem(itemText, "now is "),
		tVerbatim,
		mkItem(itemText, "the"),
		tVerbatim,
		mkItem(itemText, " time"),
		tEOF}},
	{"verbatim-surrounded", "now is =the= time\nthis is some more text!", []item{
		mkItem(itemText, "now is "),
		tVerbatim,
		mkItem(itemText, "the"),
		tVerbatim,
		mkItem(itemText, " time"),
		tNewLine,
		mkItem(itemText, "this is some more text!"),
		tEOF}},
	{"not-verbatim", "no=w is the time=", []item{mkItem(itemText, "no=w is the time="), tEOF}},

	// code
	{"code", "~now is the time~", []item{tCode, mkItem(itemText, "now is the time"), tCode, tEOF}},
	{"code-surrounded", "now is ~the~ time", []item{
		mkItem(itemText, "now is "),
		tCode,
		mkItem(itemText, "the"),
		tCode,
		mkItem(itemText, " time"),
		tEOF}},
	{"code-surrounded", "now is ~the~ time\nthis is some more text!", []item{
		mkItem(itemText, "now is "),
		tCode,
		mkItem(itemText, "the"),
		tCode,
		mkItem(itemText, " time"),
		tNewLine,
		mkItem(itemText, "this is some more text!"),
		tEOF}},
	{"not-code", "no~w is the time~", []item{mkItem(itemText, "no~w is the time~"), tEOF}},

	// underline
	{"underline", "_now is the time_", []item{tUnderline, mkItem(itemText, "now is the time"), tUnderline, tEOF}},
	{"underline-surrounded", "now is _the_ time", []item{
		mkItem(itemText, "now is "),
		tUnderline,
		mkItem(itemText, "the"),
		tUnderline,
		mkItem(itemText, " time"),
		tEOF}},
	{"underline-surrounded", "now is _the_ time\nthis is some more text!", []item{
		mkItem(itemText, "now is "),
		tUnderline,
		mkItem(itemText, "the"),
		tUnderline,
		mkItem(itemText, " time"),
		tNewLine,
		mkItem(itemText, "this is some more text!"),
		tEOF}},
	{"not-underline", "no_w is the time_", []item{mkItem(itemText, "no_w is the time_"), tEOF}},
}

func equal(i1, i2 []item, checkPos bool, t *testing.T) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			t.Logf("types not equal: %s: %T, %s: %T", i1[k].val, i1[k].typ, i2[k].val, i2[k].typ)
			return false
		}
		if i1[k].val != i2[k].val {
			t.Logf("vals not equal: %s, %s", i1[k].val, i2[k].val)
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			t.Logf("pos not equal: %d, %d", i1[k].pos, i2[k].pos)
			return false
		}
	}
	return true
}

func collect(t *lexTest) (items []item) {
	l := lex(t.name, t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF {
			break
		}
	}
	return
}

func TestLex(t *testing.T) {
	for _, tc := range lexTests {
		t.Run(tc.name, func(t *testing.T) {
			items := collect(&tc)
			if !equal(items, tc.items, false, t) {
				t.Errorf("got\n\t%v\nexpected\n\t%v", items, tc.items)
			}
		})
	}
}
