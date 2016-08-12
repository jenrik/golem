package golem

import (
	"net/url"
	"testing"
)

type ResolvData struct {
	link   url.URL
	rel    url.URL
	result string
}

func TestResolveLink(t *testing.T) {
	var ResolvTestData []ResolvData = []ResolvData{
		{ // 1
			link:   parseUrl("../bar"),
			rel:    parseUrl("http://example.com/foo"),
			result: "http://example.com/bar",
		},
		{ // 2
			link:   parseUrl("../bar"),
			rel:    parseUrl("http://example.com/foo/"),
			result: "http://example.com/bar",
		},
		{ // 3
			link:   parseUrl("../../bar"),
			rel:    parseUrl("http://example.com/foo"),
			result: "http://example.com/bar",
		},
		{ // 4
			link:   parseUrl("/bar"),
			rel:    parseUrl("http://example.com/foo"),
			result: "http://example.com/bar",
		},
		{ // 5
			link:   parseUrl("bar"),
			rel:    parseUrl("http://example.com/foo"),
			result: "http://example.com/foo/bar",
		},
		{ // 6
			link:   parseUrl("../bar"),
			rel:    parseUrl("http://example.com/foo/tar"),
			result: "http://example.com/foo/bar",
		},
		{ // 7
			link:   parseUrl("../"),
			rel:    parseUrl("http://example.com/foo/"),
			result: "http://example.com/",
		},
		{ // 8
			link:   parseUrl("/"),
			rel:    parseUrl("http://example.com/foo/"),
			result: "http://example.com/",
		},
		{ // 9
			link:   parseUrl("//"),
			rel:    parseUrl("http://example.com/foo/"),
			result: "http://example.com",
		},
		{ // 10
			link:   parseUrl("/../"),
			rel:    parseUrl("http://example.com/foo/"),
			result: "http://example.com/../", // Note: is this correct?
		},
		{ // 10
			link:   parseUrl("%2E%2E/"),
			rel:    parseUrl("http://example.com/foo/"),
			result: "http://example.com/foo/../",
		},
	}

	for i, test := range ResolvTestData {
		if result := ResolveLink(&test.link, &test.rel); result.String() != test.result {
			t.Errorf("test %v: Expected %v got %v", i+1, test.result, result.String())
		}
	}
}

func parseUrl(l string) url.URL {
	link, _ := url.Parse(l)
	return *link
}
