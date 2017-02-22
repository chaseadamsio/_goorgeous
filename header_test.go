package goorgeous

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestExtractOrgHeaders(t *testing.T) {
	source := "./testdata/test.org"
	golden := "./testdata/test.orgheaders.golden"
	content, err := os.Open(source)
	if err != nil {
		t.Fatalf("Could not open file %s: %s", source, err)
	}
	r := bufio.NewReader(content)
	fm, err := ExtractOrgHeaders(r)
	if err != nil {
		t.Fatalf("Could not extract org headers: %s", err)
	}

	if *update {
		if err := ioutil.WriteFile(golden, fm, 0644); err != nil {
			t.Errorf("failed to write %s file: %s", golden, err)
		}
		return
	}

	gld, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatalf("failed to read %s file: %s", golden, err)
	}

	if !bytes.Equal(fm, gld) {
		t.Errorf("ExtractOrgHeaders() from %s = %s\nwants: %s", source, string(fm), gld)
	}
}

func TestOrgHeaders(t *testing.T) {
	testCases := map[string]struct {
		in       string
		expected map[string]interface{}
	}{
		"no-content-headers": {"#title:\n",
			map[string]interface{}{},
		},
		"one-content-header": {"#title: my org content\n#+author:",
			map[string]interface{}{},
		},
		"basic-happy-path": {"#+title: my org mode content\n#+author: Chase Adams\n#+description: This is my description!",
			map[string]interface{}{
				"title":       "my org mode content",
				"author":      "Chase Adams",
				"description": "This is my description!",
			}},
		"basic-happy-path-new-content-after": {"#+title: my org mode content\n#+author: Chase Adams\n#+description: This is my description!\n* This shouldn't get captured!",
			map[string]interface{}{
				"title":       "my org mode content",
				"author":      "Chase Adams",
				"description": "This is my description!",
			}},
		"basic-happy-path-with-tags": {"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+tags: org-content org-mode hugo\n",
			map[string]interface{}{
				"title":       "my org mode tags content",
				"author":      "Chase Adams",
				"description": "This is my description!",
				"tags":        []string{"org-content", "org-mode", "hugo"},
			}},
		"basic-happy-path-with-categories": {"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+categories: org-content org-mode hugo\n",
			map[string]interface{}{
				"title":       "my org mode tags content",
				"author":      "Chase Adams",
				"description": "This is my description!",
				"categories":  []string{"org-content", "org-mode", "hugo"},
			}},
		"basic-happy-path-with-aliases": {"#+title: my org mode tags content\n#+author: Chase Adams\n#+description: This is my description!\n#+aliases: /org/content /org/mode /hugo\n",
			map[string]interface{}{
				"title":       "my org mode tags content",
				"author":      "Chase Adams",
				"description": "This is my description!",
				"aliases":     []string{"/org/content", "/org/mode", "/hugo"},
			}},
	}

	for caseName, tc := range testCases {
		out, err := OrgHeaders([]byte(tc.in))
		if err != nil {
			t.Fatalf("OrgHeaders() failed: %s", err)
		}
		for k, v := range tc.expected {
			switch out[k].(type) {
			case []string:
				outSlice := out[k].([]string)
				vSlice := v.([]string)
				for idx, val := range outSlice {
					if val != vSlice[idx] {
						t.Errorf("%s OrgHeaders() %v = %v\n wants: %v\n", caseName, k, out[k], tc.expected[k])
					}
				}
			case string:
				if out[k] != v {
					t.Errorf("%s OrgHeaders() %v = %v\n wants: %v\n", caseName, k, out[k], tc.expected[k])
				}
			case nil:
				t.Errorf("%s OrgHeaders() %v is nil", caseName, k)
			default:
			}
		}

	}
}
