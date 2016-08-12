package action

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestAttributeExtract(t *testing.T) {
	str := "<html><body><a href=\"example.com\">foo</div></body></html>"
	doc, _ := html.Parse(strings.NewReader(str))

	config := make(map[string]string)
	config["selector"] = "a"
	config["attribute"] = "href"

	action := new(ActionAttribute)
	output := action.Extract(doc, config)

	if output != "example.com" {
		t.FailNow()
	}
}

func TestAttributeValidateConfig(t *testing.T) {
	config := make(map[string]string)

	action := new(ActionAttribute)
	output := action.ValidateConfig(config)
	if output {
		t.FailNow()
	}

	config["selector"] = "a"
	config["attribute"] = "href"
	output = action.ValidateConfig(config)
	if !output {
		t.FailNow()
	}
}
