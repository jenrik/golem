package golem

import (
	"net/url"
	"testing"
)

type testSet struct {
	Match   bool
	Link    string
	Matcher PageMatcher
}

func TestCheck(t *testing.T) {
	tests := []testSet{
		testSet{
			true,
			"http://user@example.com/foo?tar=kar#frag",
			PageMatcher{
				[]string{"http"},
				[]string{"user"},
				[]string{"example.com"},
				[]map[string]string{{
					"min": "1",
					"max": "1",
					"1":   "foo",
				}},
				[]map[string]string{{
					"tar": "kar",
				}},
				[]string{"frag"},
			},
		},
		testSet{
			false,
			"not://user@example.com/foo?tar=kar#frag",
			PageMatcher{
				[]string{"http"},
				[]string{"user"},
				[]string{"example.com"},
				[]map[string]string{{
					"min": "1",
					"max": "1",
					"1":   "foo",
				}},
				[]map[string]string{{
					"tar": "kar",
				}},
				[]string{"frag"},
			},
		},
		testSet{
			true,
			"http://example.com/",
			PageMatcher{
				nil,
				nil,
				[]string{"example.com"},
				nil,
				nil,
				nil,
			},
		},
	}
	// ToDo write more test urls

	for i, test := range tests {
		if url, err := url.Parse(test.Link); err == nil {
			match, err := test.Matcher.Check(url)
			if err != nil {
				t.Error(err)
			}
			if match != test.Match {
				t.Logf("Check didn't match expected result for test %v", i)
			}
		} else {
			t.Error(err)
		}
	}
}
