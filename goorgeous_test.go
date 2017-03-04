package goorgeous

import (
	"flag"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestOrgCommon(t *testing.T) {
	testCases := map[string]struct {
		in       string
		expected string
	}{
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
			"<h1 id=\"a-heading\"><a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"a heading\">a heading</a></h1>\n",
		},
		"h2-link": {
			"** [[https://github.com/chaseadamsio/goorgeous][a heading]]\n",
			"<h2 id=\"a-heading\"><a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"a heading\">a heading</a></h2>\n",
		},
		"h3-link": {
			"*** [[https://github.com/chaseadamsio/goorgeous][a heading]]\n",
			"<h3 id=\"a-heading\"><a href=\"https://github.com/chaseadamsio/goorgeous\" title=\"a heading\">a heading</a></h3>\n",
		},

		"h3-emphasis": {
			"*** /a h3/ heading\n",
			"<h3 id=\"a-h3-heading\"><em>a h3</em> heading</h3>\n",
		},
		"h3-strong": {
			"*** *a h3* heading\n",
			"<h3 id=\"a-h3-heading\"><strong>a h3</strong> heading</h3>\n",
		},

		"no-inline": {"this string should have no inline changes.\n",
			"<p>this string should have no inline changes.</p>\n",
		},
		"emphasis": {
			"this string /has emphasis text/.\n",
			"<p>this string <em>has emphasis text</em>.</p>\n",
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
			"this string does not /have emphasis text/bk.\n",
			"<p>this string does not /have emphasis text/bk.</p>\n",
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
		"underline": {
			"this is _underlined text_.\n",
			"<p>this is <span style=\"text-decoration: underline;\">underlined text</span>.</p>\n",
		},
		"verbatim": {
			"this is =inline code=.\n",
			"<p>this is <code>inline code</code>.</p>\n",
		},
		"verbatim-with-equal-in-code": {
			"this is =inline=code=.\n",
			"<p>this is <code>inline=code</code>.</p>\n",
		},
		"verbatim-with-multiple-equals-in-code": {
			"this is ==inline code==.\n",
			"<p>this is <code>=inline code=</code>.</p>\n",
		},
		//"verbatim-no-surrounding-text": {
		//"==Verbatim==\n",
		//"<p><code>=Verbatim=</code></p>\n",
		//},
		"code": {
			"this has ~code~.\n",
			"<p>this has <code>code</code>.</p>\n",
		},
		"code-not": {
			"this has not~code~.\n",
			"<p>this has not~code~.</p>\n",
		},
		"code-with-tilde": {
			"this has ~~code~.\n",
			"<p>this has <code>~code</code>.</p>\n",
		},

		"src": {
			"#+BEGIN_SRC sh\necho \"foo\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\n</code></pre>\n",
		},
		"src_multiline": {
			"#+BEGIN_SRC sh\necho \"foo\"\necho \"bar\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\necho \"bar\"\n</code></pre>\n",
		},
		"src_multiline_multi_newline": {
			"#+BEGIN_SRC sh\necho \"foo\"\n\necho \"bar\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\n\necho \"bar\"\n</code></pre>\n",
		},
		"SRC_MULTILINE_MANY_MULTI_NEWLINE": {
			"#+BEGIN_SRC sh\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n#+END_SRC\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n</code></pre>\n",
		},
		"SRC_MULTILINE_MANY_MULTI_NEWLINE_TEXT": {
			"#+BEGIN_SRC text\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n#+END_SRC",
			"<pre><code class=\"language-text\">\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n\n/In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam\ndictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus\nelementum semper nisi./\n</code></pre>",
		},
		"EXAMPLE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\n</code></pre>\n",
		},
		"EXAMPLE_MULTILINE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\necho \"bar\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\necho \"bar\"\n</code></pre>\n",
		},
		"EXAMPLE_MULTILINE_MULTI_NEWLINE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\n\necho \"bar\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\n\necho \"bar\"\n</code></pre>\n",
		},
		"EXAMPLE_MULTILINE_MANY_MULTI_NEWLINE": {
			"#+BEGIN_EXAMPLE sh\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n#+END_EXAMPLE\n",
			"<pre><code class=\"language-sh\">\necho \"foo\"\n\necho \"bar\"\n\necho \"foo\"\n\necho \"bar\"\n</code></pre>\n",
		},
		"QUOTE": {
			"#+BEGIN_QUOTE\nthis is a quote.\n#+END_QUOTE\n",
			"<blockquote>\nthis is a quote.\n</blockquote>\n",
		},
		"QUOTE_MULTILINE": {
			"#+BEGIN_QUOTE\nthis is a quote\nwith multiple lines.\n#+END_QUOTE\n",
			"<blockquote>\nthis is a quote\nwith multiple lines.\n</blockquote>\n",
		},
		"CENTER": {
			"#+BEGIN_CENTER\nthis is a centered block.\n#+END_CENTER\n",
			"<center>\nthis is a centered block.\n</center>\n",
		},
		"CENTER_MULTILINE": {
			"#+BEGIN_CENTER\nthis is a\nmulti-lined centered block.\n#+END_CENTER\n",
			"<center>\nthis is a\nmulti-lined centered block.\n</center>\n",
		},
		"center_multiline_mutli_newline": {
			"#+BEGIN_CENTER\nthis is a\n\nmulti-lined centered block.\n#+END_CENTER\n",
			"<center>\nthis is a\n\nmulti-lined centered block.\n</center>\n",
		},
	}

	for caseName, tc := range testCases {
		t.Run(caseName, func(t *testing.T) {
			out := string(OrgCommon([]byte(tc.in)))
			if out != tc.expected {
				t.Errorf("%s failed. \n\t got: %s<EOF>\n\twant: %s", caseName, out, tc.expected)
			}
		})
	}
}

func BenchmarkOrgCommon(b *testing.B) {
	in := "* Title 1\nThis exports.\n* Title 2                                                          :noexport:\nThis should not export.\n* Org table\n|---+---+---|\n| a | b | c |\n|---+---+---|\n| d | e | f |\n|---+---+---|\n\n* Formatting\n** Fonts\n| Format           | Org mode markup syntax |\n| *Bold*           | =*Bold*=               |\n| /Italics/        | =/Italics/=            |\n| _Underline_      | =_Underline_=          |\n| =Verbatim=       | ==Verbatim==           |\n| +Strike-through+ | =+Strike-through+=     |\n==Verbatim==\n=Verbatim=\nthis is ==inline code==.\nthis has ~~code~.\n** Start a new paragraph\nAn empty line starts a new paragraph.\n#+BEGIN_SRC text\n/Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo\nligula nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque\neu, sem. Nulla consequat massa quis enim./\n"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		OrgCommon([]byte(in))
	}
}
