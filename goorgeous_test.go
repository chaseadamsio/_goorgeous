package goorgeous

import "testing"

func TestOrgHeaders(t *testing.T) {
	testCases := []struct {
		in       string
		expected map[string]interface{}
	}{
		{"#+TITLE: my org mode content\n#+author: Chase Adams\n#+DESCRIPTION: This is my description!",
			map[string]interface{}{
				"Title":       "my org mode content",
				"Author":      "Chase Adams",
				"Description": "This is my description!",
			}},
		{"#+TITLE: my org mode content\n#+author: Chase Adams\n#+DESCRIPTION: This is my description!\n* This shouldn't get captured!",
			map[string]interface{}{
				"Title":       "my org mode content",
				"Author":      "Chase Adams",
				"Description": "This is my description!",
			}},
	}

	for _, tc := range testCases {
		out, err := OrgHeaders([]byte(tc.in))
		if err != nil {
			t.Fatalf("OrgHeaders() failed: %s", err)
		}
		for k, v := range tc.expected {
			if out[k] != v {
				t.Errorf("OrgHeaders() = %v\n wants: %v\n", out[k], tc.expected[k])
			}
		}
	}
}

func TestOrgCommon(t *testing.T) {
	testCases := []struct {
		in       string
		expected string
	}{
		{"* A Headline Level 1!\n** A Headline Level 2!\n", "<h1 id=\"a-headline-level-1\">A Headline Level 1!</h1>\n\n<h2 id=\"a-headline-level-2\">A Headline Level 2!</h2>\n"},
		{"** A Headline Level 2!\n", "<h2 id=\"a-headline-level-2\">A Headline Level 2!</h2>\n"},
		{"**Not A Headline Level 2!\n", "<p>**Not A Headline Level 2!</p>\n"},
		{" - my definition list :: this is a definition\n", "<dl>\n<dt>my definition list</dt>\n\n<dd>this is a definition</dd>\n</dl>\n"},
		{"# this is a comment\n", "<!-- this is a comment -->\n"},
		{"this is a paragraph\n", "<p>this is a paragraph</p>\n"},
		{"      this is a paragraph      \n", "<p>this is a paragraph</p>\n"},
		{"contains =verbatim= code.\n", "<p>contains <code>verbatim</code> code.</p>\n"},
		{"contains = not a verbatim = code because spaces.\n", "<p>contains = not a verbatim = code because spaces.</p>\n"},
		{"contains just an = sign.\n", "<p>contains just an = sign.</p>\n"},
		{"this is a paragraph with an /emphasis/.\n", "<p>this is a paragraph with an <em>emphasis</em>.</p>\n"},
		{"this is a paragraph with an /emphasis/\n", "<p>this is a paragraph with an <em>emphasis</em></p>\n"},
		{"/begins with/ an /emphasis/\n", "<p><em>begins with</em> an <em>emphasis</em></p>\n"},
		{"/begins with/ an /emphasis/ and ends with a *bold*!\n", "<p><em>begins with</em> an <em>emphasis</em> and ends with a <strong>bold</strong>!</p>\n"},
		{"/begins with/ an /emphasis and contains a *bold*/!\n", "<p><em>begins with</em> an <em>emphasis and contains a <strong>bold</strong></em>!</p>\n"},
		{"/begins with/ a *bold that contains an /emphasis/*!\n", "<p><em>begins with</em> a <strong>bold that contains an <em>emphasis</em></strong>!</p>\n"},
		{"This +is a strikethrough+.\n", "<p>This <del>is a strikethrough</del>.</p>\n"},
		{"#+BEGIN_SRC sh\n echo 'hello!'\necho 'world!'\n#+END_SRC sh\n", "<pre><code>echo 'hello!'</code></pre>"},
	}

	for _, tc := range testCases {
		out := OrgCommon([]byte(tc.in))
		if string(out) != tc.expected {
			t.Errorf("OrgCommon(%s) = %s\nwants %s", tc.in, out, tc.expected)
		}
	}
}
