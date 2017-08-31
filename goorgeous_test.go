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
		"newline": {"this string should have\nan inline change.\n",
			"<p>this string should have\nan inline change.</p>\n",
		},
		"double-newline": {"this string should have\nan inline change.\n\nAnd a new paragraph.\n",
			"<p>this string should have\nan inline change.</p>\n\n<p>And a new paragraph.</p>\n",
		},
		"emphasis": {
			"this string /has emphasis text/.\n",
			"<p>this string <em>has emphasis text</em>.</p>\n",
		},
		"emphasis-with-dot-at-front": {
			"this string /.has emphasis text/.\n",
			"<p>this string <em>.has emphasis text</em>.</p>\n",
		},
		"emphasis-with-slash-inside": {
			"this string /has a slash/inside and emphasis text/.\n",
			"<p>this string <em>has a slash/inside and emphasis text</em>.</p>\n",
		},
		"emphasis-with-slash-and-space-inside": {
			"this string /has a slash/ inside and emphasis text/.\n",
			"<p>this string <em>has a slash</em> inside and emphasis text/.</p>\n",
		},
		"emphasis-with-slash-inside-and-another-emphasis": {
			"this string /has a slash/inside and emphasis text/ and another /set of emphasis/.\n",
			"<p>this string <em>has a slash/inside and emphasis text</em> and another <em>set of emphasis</em>.</p>\n",
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
		"emphasis-inside-simple-ol": {
			"1. this\n2. is\n3. an\n4. ordered\n5. /list/\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>ordered</li>\n<li><em>list</em></li>\n</ol>\n",
		},
		"emphasis-inside-simple-ul": {
			"- this\n- is\n- an\n- unordered\n- /list/\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li><em>list</em></li>\n</ul>\n",
		},
		"bold": {
			"this string *has bold text*.\n",
			"<p>this string <strong>has bold text</strong>.</p>\n",
		},
		"bold-with-dot-at-front": {
			"this string *.has bold text*.\n",
			"<p>this string <strong>.has bold text</strong>.</p>\n",
		},
		"bold-with-asterisk-inside": {
			"this string *has *asterisk and bold text*.\n",
			"<p>this string <strong>has *asterisk and bold text</strong>.</p>\n",
		},
		"bold-not-no-spaces": {
			"this string*doesn't have bold text*.\n",
			"<p>this string*doesn't have bold text*.</p>\n",
		},
		"bold-not-no-spaces-split": {
			"this string*doesn't have bold text.*\n",
			"<p>this string*doesn't have bold text.*</p>\n",
		},
		"bold-inside-simple-ol": {
			"1. this\n2. is\n3. an\n4. ordered\n5. *list*\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>ordered</li>\n<li><strong>list</strong></li>\n</ol>\n",
		},
		"bold-inside-simple-ul": {
			"- this\n- is\n- an\n- unordered\n- *list*\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li><strong>list</strong></li>\n</ul>\n",
		},
		"underline": {
			"this is _underlined text_.\n",
			"<p>this is <span style=\"text-decoration: underline;\">underlined text</span>.</p>\n",
		},
		"underline-with-dot-at-front": {
			"this is _.underlined text_.\n",
			"<p>this is <span style=\"text-decoration: underline;\">.underlined text</span>.</p>\n",
		},
		"verbatim": {
			"this is =inline code=.\n",
			"<p>this is <code>inline code</code>.</p>\n",
		},
		"verbatim-with-dot-at-front": {
			"this is =.inline code=.\n",
			"<p>this is <code>.inline code</code>.</p>\n",
		},
		"verbatim-with-equal-in-code": {
			"this is =inline=code=.\n",
			"<p>this is <code>inline=code</code>.</p>\n",
		},
		"verbatim-with-multiple-equals-in-code": {
			"this is ==inline code==.\n",
			"<p>this is <code>=inline code=</code>.</p>\n",
		},
		"verbatim-no-surrounding-text": {
			"==Verbatim==\n",
			"<p><code>=Verbatim=</code></p>\n",
		},
		"verbatim-inside-simple-ol": {
			"1. this\n2. is\n3. an\n4. ordered\n5. =list=\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>ordered</li>\n<li><code>list</code></li>\n</ol>\n",
		},
		"verbatim-inside-simple-ul": {
			"- this\n- is\n- an\n- unordered\n- =list=\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li><code>list</code></li>\n</ul>\n",
		},
		"code": {
			"this has ~code~.\n",
			"<p>this has <code>code</code>.</p>\n",
		},
		"code-with-dot-at-front": {
			"this has ~.code~.\n",
			"<p>this has <code>.code</code>.</p>\n",
		},
		"code-not": {
			"this has not~code~.\n",
			"<p>this has not~code~.</p>\n",
		},
		"code-with-tilde": {
			"this has ~~code~.\n",
			"<p>this has <code>~code</code>.</p>\n",
		},
		"code-inside-simple-ol": {
			"1. this\n2. is\n3. an\n4. ordered\n5. ~list~\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>ordered</li>\n<li><code>list</code></li>\n</ol>\n",
		},
		"code-inside-simple-ul": {
			"- this\n- is\n- an\n- unordered\n- ~list~\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li><code>list</code></li>\n</ul>\n",
		},
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
			"this has [[file:../gopher.gif][a uni-gopher]] as an image.\n",
			"<p>this has <img src=\"../gopher.gif\" alt=\"a uni-gopher\" title=\"a uni-gopher\" /> as an image.</p>\n",
		},
		"link-inside-simple-ol": {
			"1. this\n2. is\n3. an\n4. ordered\n5. list with [[https://github.com/chaseadamsio/goorgeous][goorgeous by chaseadamsio]] as a link\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>ordered</li>\n<li>list with <a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"goorgeous by chaseadamsio\">goorgeous by chaseadamsio</a> as a link</li>\n</ol>\n",
		},
		"link-inside-simple-ul": {
			"- this\n- is\n- an\n- unordered\n- list with [[https://github.com/chaseadamsio/goorgeous][goorgeous by chaseadamsio]] as a link\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li>list with <a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"goorgeous by chaseadamsio\">goorgeous by chaseadamsio</a> as a link</li>\n</ul>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func TestRenderingFootnotes(t *testing.T) {
	testCases := testCase{"Test 1[fn:1] and Test 2[fn: 2] and Test 3[fn:3] also test lettres[fn:let] then final test[fn:4]\n\n[fn:let] what?\n\n\n[fn:6] And test it[fn:7].\n\n* Footnotes\n\n[fn:1] Test 1\n\n[fn:3] Test 3\n\n[fn:2] Test2\n\n[fn:5] missing?\n\n[fn:6] Six?\n\n[fn:7] Seven??",
		"<p>Test 1<sup class=\"footnote-ref\" id=\"fnref:1\"><a rel=\"footnote\" href=\"#fn:1\">1</a></sup> and Test 2[fn: 2] and Test 3<sup class=\"footnote-ref\" id=\"fnref:3\"><a rel=\"footnote\" href=\"#fn:3\">2</a></sup> also test lettres<sup class=\"footnote-ref\" id=\"fnref:let\"><a rel=\"footnote\" href=\"#fn:let\">3</a></sup> then final test<sup class=\"footnote-ref\" id=\"fnref:4\"><a rel=\"footnote\" href=\"#fn:4\">4</a></sup></p>\n\n<h1 id=\"footnotes\">Footnotes</h1>\n<div class=\"footnotes\">\n\n<hr />\n\n<ol>\n<li id=\"fn:1\">Test 1 <a class=\"footnote-return\" href=\"#fnref:1\"><sup>↩</sup></a></li>\n\n<li id=\"fn:3\">Test 3 <a class=\"footnote-return\" href=\"#fnref:3\"><sup>↩</sup></a></li>\n\n<li id=\"fn:let\">what? <a class=\"footnote-return\" href=\"#fnref:let\"><sup>↩</sup></a></li>\n\n<li id=\"fn:4\">DEFINITION NOT FOUND <a class=\"footnote-return\" href=\"#fnref:4\"><sup>↩</sup></a></li>\n</ol>\n</div>\n"}

	flags := blackfriday.HTML_USE_XHTML
	flags |= blackfriday.LIST_ITEM_BEGINNING_OF_LIST
	flags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	var parameters blackfriday.HtmlRendererParameters
	parameters.FootnoteReturnLinkContents = "<sup>↩</sup>"

	renderer := blackfriday.HtmlRendererWithParameters(flags, "", "", parameters)
	out := Org([]byte(testCases.in), renderer)

	if !bytes.Equal(out, []byte(testCases.expected)) {
		t.Errorf("Footnote for Org() from: \n %s \n result: \n  %s\n wants:\n  %s", testCases.in, string(out), testCases.expected)
	}
}

func TestRenderingBlock(t *testing.T) {

	testCases := map[string]testCase{
		"SRC": {
			"#+BEGIN_SRC sh\necho \"foo\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\n</code></pre>\n",
		},
		"SRC_MULTILINE": {
			"#+BEGIN_SRC sh\necho \"foo\"\necho \"bar\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\necho &quot;bar&quot;\n</code></pre>\n",
		},
		"SRC_MULTILINE_MULTI_NEWLINE": {
			"#+BEGIN_SRC sh\necho \"foo\"\n\necho \"bar\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\n\necho &quot;bar&quot;\n</code></pre>\n",
		},
		"SRC_MULTILINE_MANY_MULTI_NEWLINE": {
			"#+BEGIN_SRC sh\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\n\necho &quot;bar&quot;\n\necho &quot;foo&quot;\n\necho &quot;bar&quot;\n</code></pre>\n",
		},
		"SRC_MULTILINE_MANY_MULTI_NEWLINE_TEXT": {
			"#+BEGIN_SRC text\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n#+END_SRC",
			"<pre><code class=\"language-text\">\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n</code></pre>\n",
		},
		"EXAMPLE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\n</code></pre>\n",
		},
		"EXAMPLE_MULTILINE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\necho \"bar\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\necho &quot;bar&quot;\n</code></pre>\n",
		},
		"EXAMPLE_MULTILINE_MULTI_NEWLINE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\n\necho \"bar\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\n\necho &quot;bar&quot;\n</code></pre>\n",
		},
		"EXAMPLE_MULTILINE_MANY_MULTI_NEWLINE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho &quot;foo&quot;\n\necho &quot;bar&quot;\n\necho &quot;foo&quot;\n\necho &quot;bar&quot;\n</code></pre>\n",
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
		"table-with-inlined-elements": {
			"| Format           | Org mode markup syntax |\n| *Bold*           | =*Bold*=               |\n| /Italics/        | =/Italics/=            |\n| _Underline_      | =_Underline_=          |\n| =Verbatim=       | ==Verbatim== |\n| +Strike-through+ | =+Strike-through+=     |\n",
			"\n<table>\n<tbody>\n<tr>\n<td>Format</td>\n<td>Org mode markup syntax</td>\n</tr>\n\n<tr>\n<td><strong>Bold</strong></td>\n<td><code>*Bold*</code></td>\n</tr>\n\n<tr>\n<td><em>Italics</em></td>\n<td><code>/Italics/</code></td>\n</tr>\n\n<tr>\n<td><span style=\"text-decoration: underline;\">Underline</span></td>\n<td><code>_Underline_</code></td>\n</tr>\n\n<tr>\n<td><code>Verbatim</code></td>\n<td><code>=Verbatim=</code></td>\n</tr>\n\n<tr>\n<td><del>Strike-through</del></td>\n<td><code>+Strike-through+</code></td>\n</tr>\n</tbody>\n</table>\n",
		},
		"table-single-cell": {
			"| r |\n",
			"\n<table>\n<tbody>\n<tr>\n<td>r</td>\n</tr>\n</tbody>\n</table>\n",
		},
	}

	testOrgCommon(testCases, t)
}

func TestLists(t *testing.T) {
	testCases := map[string]testCase{
		"simple-definition": {
			"- definition lists :: these are useful sometimes\n- item 2 :: M-RET again gives another item, and long lines wrap in a tidy way underneath the definition\n",
			"<dl>\n<dt>definition lists</dt>\n<dd>these are useful sometimes</dd>\n<dt>item 2</dt>\n<dd>M-RET again gives another item, and long lines wrap in a tidy way underneath the definition</dd>\n</dl>\n",
		},
		"simple-ol": {
			"1. this\n2. is\n3. an\n4. ordered\n5. list\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>ordered</li>\n<li>list</li>\n</ol>\n",
		},
		"ol-change-number": {
			"1. this\n2. is\n3. [@10] changed to 10\n4. ordered\n5. list\n",
			"<ol>\n<li>this</li>\n<li>is</li>\n<li value=\"10\">changed to 10</li>\n<li>ordered</li>\n<li>list</li>\n</ol>\n",
		},
		"simple-ul-plus-sign": {
			"+ this\n+ is\n+ an\n+ unordered\n+ list\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li>list</li>\n</ul>\n",
		},
		"simple-ul-dash": {
			"- this\n- is\n- an\n- unordered\n- list\n",
			"<ul>\n<li>this</li>\n<li>is</li>\n<li>an</li>\n<li>unordered</li>\n<li>list</li>\n</ul>\n",
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

func TestRenderingComplexTexts(t *testing.T) {
	testCases := map[string]testCase{
		"newline": {
			"** Start a new paragraph\nAn empty line starts a new paragraph.\n#+BEGIN_SRC text\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n#+END_SRC\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi./\n",
			"<h2 id=\"start-a-new-paragraph\">Start a new paragraph</h2>\n\n<p>An empty line starts a new paragraph.</p>\n\n<pre><code class=\"language-text\">\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n</code></pre>\n\n<p><em>Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, sem. Nulla consequat massa quis enim.</em></p>\n\n<p><em>In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi.</em></p>\n",
		},
	}
	testOrgCommon(testCases, t)
}

func TestFixedWidthAreas(t *testing.T) {
	testCases := map[string]testCase{
		"single-line-fixed-width": {
			": Fixed width area\n",
			"<pre class=\"example\">\nFixed width area\n</pre>\n",
		},
		"double-line-fixed-width": {
			": Line 1\n:    Line 2\n",
			"<pre class=\"example\">\nLine 1\n   Line 2\n</pre>\n",
		},
		"paragraph-fixed-width-transition": {
			"     Here is an example\n        : Some example from a text file.",
			"<p>Here is an example</p>\n<pre class=\"example\">\nSome example from a text file.\n</pre>\n",
		},
		"babel-output-fixed-width": {
			"#+RESULTS[26164fcdb01a7c6a9329d20d4754e15f7739ad20]:\n: Hello, World!\n: The output of this code can be seen in the post.",
			"<pre class=\"example\">\nHello, World!\nThe output of this code can be seen in the post.\n</pre>\n",
		},
		"fixed-width-to-paragraph-transition": {
			": example\nEnd\n",
			"<pre class=\"example\">\nexample\n</pre>\n<p>End</p>\n",
		},
		"fixed-width-paragraph-comment": {
			"p\n: e\n# comment\n",
			"<p>p</p>\n<pre class=\"example\">\ne\n</pre>\n<!-- comment -->\n",
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
