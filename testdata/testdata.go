package testdata

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

// cache the files
var files = map[string]string{}

// GetOrgStr takes a filepath relative to testdata/in directory
// (not the test file that's using it)
// as an example, in parser/list_test.go:
// GetOrg(ordered-list/basic.org)
func GetOrgStr(filename string) (content string) {
	_, callerFilename, _, ok := runtime.Caller(0) // get this file's filepath
	if !ok {
		panic("Unable to determine runtime caller.")
	}

	filepath := path.Join(path.Dir(callerFilename), "in", filename)
	if found, ok := files[filepath]; ok {
		content = found
	} else {
		contentByt, err := ioutil.ReadFile(filepath)
		if err != nil {
			panic(fmt.Errorf("Could not open file %s: %s", filepath, err))
		}
		content = string(contentByt)
		files[filepath] = content

	}
	return content
}

const (
	OrderedListBasic                   = "ordered-list/basic.org"
	OrderedListNotAList                = "ordered-list/not-a-list.org"
	OrderedListWithStartingNewline     = "ordered-list/with-starting-newline.org"
	OrderedListFollowParagraph         = "ordered-list/with-follow-paragraph.org"
	OrderedListFollowNumberNotList     = "ordered-list/with-follow-number-not-list.org"
	OrderedListFollowAsteriskHeading   = "ordered-list/with-follow-asterisk-heading.org"
	OrderedListWithFollowUnOrderedList = "ordered-list/with-follow-unordered-list.org"
	OrderedListWithNestedOrderedList   = "ordered-list/with-nested-ordered-list.org"
	OrderedListWithNestedContent       = "ordered-list/with-nested-content.org"

	UnorderedListBasic                  = "unordered-list/basic.org"
	UnorderedListNotAList               = "unordered-list/not-a-list.org"
	UnorderedListWithStartingNewline    = "unordered-list/with-starting-newline.org"
	UnorderedListFollowParagraph        = "unordered-list/with-follow-paragraph.org"
	UnorderedListFollowDashNotList      = "unordered-list/with-follow-dash-not-list.org"
	UnorderedListFollowAsteriskHeading  = "unordered-list/with-follow-asterisk-heading.org"
	UnorderedListWithFollowOrderedList  = "unordered-list/with-follow-ordered-list.org"
	UnorderedListWithNestedOrderedList  = "unordered-list/with-nested-ordered-list.org"
	UnorderedListWithDeepNestedChildren = "unordered-list/with-2-deep-nested-children-list.org"
	UnorderedListWithNestedContent      = "unordered-list/with-nested-content.org"
)
