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
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"spaces", " \t\n", []item{mkItem(itemText, " \t"), tNewLine, tEOF}},
	{"text", "now is the time", []item{mkItem(itemText, "now is the time"), tEOF}},
	{"text", "now is the time\n", []item{mkItem(itemText, "now is the time"), tNewLine, tEOF}},
	{"text", "now is the time\n\n", []item{mkItem(itemText, "now is the time"), tNewLine, tNewLine, tEOF}},

	// BASIC HEADLINES
	{"h1", "* A h1 Headline", []item{tH1, mkItem(itemText, " A h1 Headline"), tEOF}},
	{"not-h1", "not an * h1 Headline", []item{mkItem(itemText, "not an * h1 Headline"), tEOF}},
	{"alt-not-h1", " * not an h1 Headline", []item{mkItem(itemText, " * not an h1 Headline"), tEOF}},
	{"alt-not-h1-2", "*not an h1 Headline", []item{mkItem(itemText, "*not an h1 Headline"), tEOF}},
	{"h2", "** A h2 Headline", []item{tH2, mkItem(itemText, " A h2 Headline"), tEOF}},
	{"not-h2", "not an ** h2 Headline", []item{mkItem(itemText, "not an ** h2 Headline"), tEOF}},
	{"alt-not-h2", " ** not an h2 Headline", []item{mkItem(itemText, " ** not an h2 Headline"), tEOF}},
	{"alt-not-h2-2", "**not an h2 Headline", []item{mkItem(itemText, "**not an h2 Headline"), tEOF}},
	{"h3", "*** A h3 Headline", []item{tH3, mkItem(itemText, " A h3 Headline"), tEOF}},
	{"not-h3", "not an *** h3 Headline", []item{mkItem(itemText, "not an *** h3 Headline"), tEOF}},
	{"alt-not-h3", " *** not an h3 Headline", []item{mkItem(itemText, " *** not an h3 Headline"), tEOF}},
	{"alt-not-h3-2", "***not an h3 Headline", []item{mkItem(itemText, "***not an h3 Headline"), tEOF}},
	{"h4", "**** A h4 Headline", []item{tH4, mkItem(itemText, " A h4 Headline"), tEOF}},
	{"not-h4", "not an **** h4 Headline", []item{mkItem(itemText, "not an **** h4 Headline"), tEOF}},
	{"alt-not-h4", " **** not an h4 Headline", []item{mkItem(itemText, " **** not an h4 Headline"), tEOF}},
	{"alt-not-h4-2", "****not an h4 Headline", []item{mkItem(itemText, "****not an h4 Headline"), tEOF}},
	{"h5", "***** A h5 Headline", []item{tH5, mkItem(itemText, " A h5 Headline"), tEOF}},
	{"not-h5", "not an ***** h5 Headline", []item{mkItem(itemText, "not an ***** h5 Headline"), tEOF}},
	{"alt-not-h5", " ***** not an h5 Headline", []item{mkItem(itemText, " ***** not an h5 Headline"), tEOF}},
	{"alt-not-h5-2", "*****not an h5 Headline", []item{mkItem(itemText, "*****not an h5 Headline"), tEOF}},
	{"h6", "****** A h6 Headline", []item{tH6, mkItem(itemText, " A h6 Headline"), tEOF}},
	{"not-h6", "not an ****** h6 Headline", []item{mkItem(itemText, "not an ****** h6 Headline"), tEOF}},
	{"alt-not-h6", " ****** not an h6 Headline", []item{mkItem(itemText, " ****** not an h6 Headline"), tEOF}},
	{"alt-not-h6-2", "******not an h6 Headline", []item{mkItem(itemText, "******not an h6 Headline"), tEOF}},

	// HEADLINES VARIANTS
	{"multi-headlines", "* A h1 Headline\n****** A h6 Headline",
		[]item{
			tH1,
			mkItem(itemText, " A h1 Headline"),
			tNewLine,
			tH6,
			mkItem(itemText, " A h6 Headline"),
			tEOF}},
	{"h1-with-text", "* A h1 Headline\nThis is a new line.\n",
		[]item{
			tH1,
			mkItem(itemText, " A h1 Headline"),
			tNewLine,
			mkItem(itemText, "This is a new line."),
			tNewLine,
			tEOF}},

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

	// COMPLEX
	{"complex-h1-emphasis-h6", "* A h1 Headline\nsome /emphasis text/.\n****** A h6 Headline",
		[]item{
			tH1,
			mkItem(itemText, " A h1 Headline"),
			tNewLine,
			mkItem(itemText, "some "),
			tEmphasis,
			mkItem(itemText, "emphasis text"),
			tEmphasis,
			mkItem(itemText, "."),
			tNewLine,
			tH6,
			mkItem(itemText, " A h6 Headline"),
			tEOF}},
	{"complex-h1-emphasis-h6", "* A h1 Headline\nsome /*emphasis* text/.\n****** A h6 Headline",
		[]item{
			tH1,
			mkItem(itemText, " A h1 Headline"),
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
			mkItem(itemText, " A h6 Headline"),
			tEOF}},
}

func equal(i1, i2 []item, checkPos bool, t *testing.T) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
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
