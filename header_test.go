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
