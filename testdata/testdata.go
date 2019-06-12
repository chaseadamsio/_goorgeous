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

	filepath := path.Join(path.Dir(callerFilename), filename)
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
	Headline1 = "in/headline/headline-1.org"

	OrderedListBasic                   = "in/ordered-list/basic.org"
	OrderedListNotAList                = "in/ordered-list/not-a-list.org"
	OrderedListWithStartingNewline     = "in/ordered-list/with-starting-newline.org"
	OrderedListFollowParagraph         = "in/ordered-list/with-follow-paragraph.org"
	OrderedListFollowNumberNotList     = "in/ordered-list/with-follow-number-not-list.org"
	OrderedListFollowAsteriskHeading   = "in/ordered-list/with-follow-asterisk-heading.org"
	OrderedListWithFollowUnOrderedList = "in/ordered-list/with-follow-unordered-list.org"
	OrderedListWithNestedOrderedList   = "in/ordered-list/with-nested-ordered-list.org"
	OrderedListWithNestedContent       = "in/ordered-list/with-nested-content.org"

	UnorderedListBasic                  = "in/unordered-list/basic.org"
	UnorderedListNotAList               = "in/unordered-list/not-a-list.org"
	UnorderedListWithStartingNewline    = "in/unordered-list/with-starting-newline.org"
	UnorderedListFollowParagraph        = "in/unordered-list/with-follow-paragraph.org"
	UnorderedListFollowDashNotList      = "in/unordered-list/with-follow-dash-not-list.org"
	UnorderedListFollowAsteriskHeading  = "in/unordered-list/with-follow-asterisk-heading.org"
	UnorderedListWithFollowOrderedList  = "in/unordered-list/with-follow-ordered-list.org"
	UnorderedListWithNestedOrderedList  = "in/unordered-list/with-nested-ordered-list.org"
	UnorderedListWithDeepNestedChildren = "in/unordered-list/with-2-deep-nested-children-list.org"
	UnorderedListWithNestedContent      = "in/unordered-list/with-nested-content.org"
)
