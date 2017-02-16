package goorgeous

import (
	"bytes"
	"flag"
	"testing"

	"github.com/russross/blackfriday"
)

var update = flag.Bool("update", false, "update golden files")

func TestIsTable(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"|some table", true},
		{"| some table", true},
		{" | not a table", false},
		{"not a table", false},
		{"*not a table", false},
		{"-not a table", false},
		{"+not a table", false},
	}

	for _, tc := range testCases {
		isTable := isTable([]byte(tc.in))
		if isTable != tc.expected {
			t.Errorf("isTable(%s) = %T\nwants: %T", tc.in, isTable, tc.expected)
		}
	}
}

func TestSkipChar(t *testing.T) {
	testCases := []struct {
		in       string
		start    int
		char     byte
		expected int
	}{
		{"   check for spaces", 0, ' ', 3},
		{" -  check for spaces", 0, ' ', 1},
		{"check     for spaces", 1, ' ', 1},
		{"check     for spaces", 5, ' ', 10},
		{"check-----for dashes", 5, '-', 10},
	}

	for _, tc := range testCases {
		skipped := skipChar([]byte(tc.in), tc.start, tc.char)
		if skipped != tc.expected {
			t.Errorf("skipChar(%s, %d, %q) = %d\nwants: %d", tc.in, tc.start, tc.char, skipped, tc.expected)
		}
	}
}

func TestIsSpace(t *testing.T) {
	testCases := []struct {
		char     byte
		expected bool
	}{

		{' ', true},
		{'+', false},
	}

	for _, tc := range testCases {
		isSpace := isSpace(tc.char)
		if isSpace != tc.expected {
			t.Errorf("isSpace(%q) = %T/nwants: %T", tc.char, isSpace, tc.expected)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"\n", true},
		{"\t", true},
		{"\t\n\t", true},
		{"\tfoo\n\t", false},
	}

	for _, tc := range testCases {
		isEmpty := isEmpty([]byte(tc.in))
		if isEmpty != tc.expected {
			t.Errorf("isEmpty(%s) = %T\nwants: %T", tc.in, isEmpty, tc.expected)
		}
	}
}

func TestCharMatch(t *testing.T) {
	testCases := []struct {
		a        byte
		b        byte
		expected bool
	}{
		{' ', ' ', true},
		{'+', '+', true},
		{'#', '#', true},
		{'-', '-', true},
		{'+', '-', false},
		{' ', '-', false},
	}

	for _, tc := range testCases {
		matches := charMatches(tc.a, tc.b)
		if matches != tc.expected {
			t.Errorf("charMatches(%q, %q) = %T\nwants: %T", tc.a, tc.b, matches, tc.expected)
		}
	}
}

func TestIsHeader(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"#+TITLE", true},
		{"#+ TITLE", false},
		{"# TITLE", false},
		{"TITLE", false},
	}

	for _, tc := range testCases {
		isHeader := IsKeyword([]byte(tc.in))
		if isHeader != tc.expected {
			t.Errorf("isHeader(%s) = %T\nwants: %T", tc.in, isHeader, tc.expected)
		}
	}
}

func TestIsComment(t *testing.T) {
	testCases := []struct {
		in       string
		expected bool
	}{
		{"# this is a comment", true},
		{"#-this is not a comment", false},
		{"#+TITLE", false},
		{"This is not a comment", false},
	}

	for _, tc := range testCases {
		isComment := isComment([]byte(tc.in))
		if isComment != tc.expected {
			t.Errorf("isComment(%s) = %T\nwants: %T", tc.in, isComment, tc.expected)
		}
	}
}

func TestGenerateComment(t *testing.T) {
	p := NewParser(blackfriday.HtmlRenderer(blackfriday.HTML_USE_XHTML, "", ""))
	var out bytes.Buffer
	text := "This is a comment and we expect it to look a certain way."
	orgComment := []byte("# " + text)
	expected := "<!-- " + text + " -->\n"
	p.generateComment(&out, orgComment)
	if out.String() != expected {
		t.Errorf("generateComment(%s) = %s\nwants: %s", text, out.String(), expected)
	}
}

type testCase struct {
	in       string
	expected string
}

func TestRenderingHeadings(t *testing.T) {
	testCases := map[string]testCase{
		"h1-basic": {
			"* a h1 heading\n",
			"<h1 id=\"a-h1-heading\">a h1 heading</h1>\n",
		},
		"h2-basic": {
			"** a h2 heading\n",
			"<h2 id=\"a-h2-heading\">a h2 heading</h2>\n",
		},
		"h3-basic": {
			"*** a h3 heading\n",
			"<h3 id=\"a-h3-heading\">a h3 heading</h3>\n",
		},
		"h4-basic": {
			"**** a h4 heading\n",
			"<h4 id=\"a-h4-heading\">a h4 heading</h4>\n",
		},
		"h5-basic": {
			"***** a h5 heading\n",
			"<h5 id=\"a-h5-heading\">a h5 heading</h5>\n",
		},
		"h6-basic": {
			"****** a h6 heading\n",
			"<h6 id=\"a-h6-heading\">a h6 heading</h6>\n",
		},

		"h1-link": {
			"* [[https://github.com/chaseadamsio/goorgeous][a heading]]\n",
			"<h1 id=\"https-github-com-chaseadamsio-goorgeous-a-heading\"><a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"a heading\">a heading</a></h1>\n",
		},
		"h2-link": {
			"** [[https://github.com/chaseadamsio/goorgeous][a heading]]\n",
			"<h2 id=\"https-github-com-chaseadamsio-goorgeous-a-heading\"><a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"a heading\">a heading</a></h2>\n",
		},
		"h3-link": {
			"*** [[https://github.com/chaseadamsio/goorgeous][a heading]]\n",
			"<h3 id=\"https-github-com-chaseadamsio-goorgeous-a-heading\"><a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"a heading\">a heading</a></h3>\n",
		},

		"h3-emphasis": {
			"*** /a h3/ heading\n",
			"<h3 id=\"a-h3-heading\"><em>a h3</em> heading</h3>\n",
		},
		"h3-strong": {
			"*** *a h3* heading\n",
			"<h3 id=\"a-h3-heading\"><strong>a h3</strong> heading</h3>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func TestRenderingInline(t *testing.T) {
	testCases := map[string]testCase{
		"no-inline": {"this string should have no inline changes.\n",
			"<p>this string should have no inline changes.</p>\n",
		},
		"emphasis": {
			"this string /has emphasis text/.\n",
			"<p>this string <em>has emphasis text</em>.</p>\n",
		},
		"emphasis-not": {
			"this string does not /have emphasis text/p.\n",
			"<p>this string does not /have emphasis text/p.</p>\n",
		},
		"emphasis-not-no-spaces": {
			"this string does not/have emphasis textp/.\n",
			"<p>this string does not/have emphasis textp/.</p>\n",
		},
		"emphasis-not-single-slash": {
			"this string does not /have emphasis text.\n",
			"<p>this string does not /have emphasis text.</p>\n",
		},
		"emphasis-not-double-slash-no-spaces": {
			"this string does not/have emphasis text. feel/me?\n",
			"<p>this string does not/have emphasis text. feel/me?</p>\n",
		},
		"emphasis-not-slash-with-link": {
			"this string does not/have emphasis text [[https://somelinkshouldntrenderaccidentalemphasis.com]].\n",
			"<p>this string does not/have emphasis text <a href=\"https://somelinkshouldntrenderaccidentalemphasis.com\" title=\"https://somelinkshouldntrenderaccidentalemphasis.com\">https://somelinkshouldntrenderaccidentalemphasis.com</a>.</p>\n",
		},
		"bold": {
			"this string *has bold text*.\n",
			"<p>this string <strong>has bold text</strong>.</p>\n",
		},
		"bold-not-no-spaces": {
			"this string*doesn't have bold text*.\n",
			"<p>this string*doesn't have bold text*.</p>\n",
		},
		"bold-not-no-spaces-split": {
			"this string*doesn't have bold text.*\n",
			"<p>this string*doesn't have bold text.*</p>\n",
		},
		"underline": {
			"this is _underlined text_.\n",
			"<p>this is <span style=\"text-decoration: underline;\">underlined text</span>.</p>\n",
		},
		"verbatim": {
			"this is =inline code=.\n",
			"<p>this is <code>inline code</code>.</p>\n",
		},
		//"verbatim-with-equal-in-code": {
		//"this is =inline=code=.\n",
		//"<p>this is <code>inline=code</code>.</p>\n",
		//},
		//"verbatim-with-multiple-equals-in-code": {
		//"this is ==inline code==.\n",
		//"<p>this is <code>=inline code=</code>.</p>\n",
		//},
	}

	testOrgCommon(testCases, t)
}

func TestRenderingLinksAndImages(t *testing.T) {

	testCases := map[string]testCase{
		"anchor-basic": {"this has [[https://github.com/chaseadamsio/goorgeous]] as a link.\n",
			"<p>this has <a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"https://github.com/chaseadamsio/goorgeous\">https://github.com/chaseadamsio/goorgeous</a> as a link.</p>\n",
		},
		"anchor-text": {
			"this has [[https://github.com/chaseadamsio/goorgeous][goorgeous by chaseadamsio]] as a link.\n",
			"<p>this has <a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"goorgeous by chaseadamsio\">goorgeous by chaseadamsio</a> as a link.</p>\n",
		},
		"image-basic": {
			"this has [[file:https://github.com/chaseadamsio/goorgeous/img.png]] as an image.\n",
			"<p>this has <img src=\"https://github.com/chaseadamsio/goorgeous/img.png\" alt=\"https://github.com/chaseadamsio/goorgeous/img.png\" title=\"https://github.com/chaseadamsio/goorgeous/img.png\" /> as an image.</p>\n",
		},
		"image-alt": {
			"this has [[file:../gopher.gif][a uni-gopher]] as an image.",
			"<p>this has <img src=\"../gopher.gif\" alt=\"a uni-gopher\" title=\"a uni-gopher\" /> as an image.</p>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func TestRenderingBlock(t *testing.T) {

	testCases := map[string]testCase{
		"SRC": {"#+BEGIN_SRC sh\n echo \"foo\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\"> echo &quot;foo&quot;\n</code></pre>\n",
		},
		"EXAMPLE": {
			"#+BEGIN_EXAMPLE sh\n echo \"foo\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\"> echo &quot;foo&quot;\n</code></pre>\n",
		},
		"QUOTE": {
			"#+BEGIN_QUOTE\nthis is a quote.\n#+END_QUOTE\n",
			"<blockquote>\n<p>\nthis is a quote.\n</p>\n</blockquote>\n",
		},
		"QUOTE_MULTILINE": {
			"#+BEGIN_QUOTE\nthis is a quote\nwith multiple lines.\n#+END_QUOTE\n",
			"<blockquote>\n<p>\nthis is a quote\n</p>\n<p>\nwith multiple lines.\n</p>\n</blockquote>\n",
		},
		"CENTER": {
			"#+BEGIN_CENTER\nthis is a centered block.\n#+END_CENTER\n",
			"<center>\n<p>\nthis is a centered block.\n</p>\n</center>\n",
		},
		"CENTER_MULTILINE": {
			"#+BEGIN_CENTER\nthis is a\nmulti-lined centered block.\n#+END_CENTER\n",
			"<center>\n<p>\nthis is a\n</p>\n<p>\nmulti-lined centered block.\n</p>\n</center>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func TestRenderingTables(t *testing.T) {
	testCases := map[string]testCase{
		"no-table-heading-no-horizontal-splits": {
			"|foo|bar|baz|\n| d | e | f |\n| g | h | i |\n",
			"\n<table>\n<tbody>\n<tr>\n<td>foo</td>\n<td>bar</td>\n<td>baz</td>\n</tr>\n\n<tr>\n<td>d</td>\n<td>e</td>\n<td>f</td>\n</tr>\n\n<tr>\n<td>g</td>\n<td>h</td>\n<td>i</td>\n</tr>\n</tbody>\n</table>\n",
		},
		"table-heading": {
			"|foo|bar|baz|\n|---+---+---|\n| d | e | f |\n| g | h | i |\n",
			"\n<table>\n<thead>\n<tr>\n<th>foo</th>\n<th>bar</th>\n<th>baz</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>d</td>\n<td>e</td>\n<td>f</td>\n</tr>\n\n<tr>\n<td>g</td>\n<td>h</td>\n<td>i</td>\n</tr>\n</tbody>\n</table>\n",
		},
		"no-table-heading-horizontal-splits": {
			"|---+---+---|\n| d | e | f |\n|---+---+---|\n| g | h | i |\n|---+---+---|\n",
			"\n<table>\n<tbody>\n<tr>\n<td>d</td>\n<td>e</td>\n<td>f</td>\n</tr>\n\n<tr>\n<td>g</td>\n<td>h</td>\n<td>i</td>\n</tr>\n</tbody>\n</table>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func TestRenderingPropertiesDrawer(t *testing.T) {
	testCases := map[string]testCase{
		"basic": {
			"* Heading\n:PROPERTIES:\n:header-args: :tangle ~/.filename\n:END:\n next block.",
			"<h1 id=\"heading\">Heading</h1>\n\n<p>next block.</p>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func testOrgCommon(testCases map[string]testCase, t *testing.T) {
	for caseName, tc := range testCases {

		out := OrgCommon([]byte(tc.in))

		if !bytes.Equal(out, []byte(tc.expected)) {
			t.Errorf("case %s for OrgCommon() from %s = %s\nwants: %s", caseName, tc.in, out, tc.expected)
		}
	}
}
