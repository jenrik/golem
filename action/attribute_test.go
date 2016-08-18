package action

import (
	"github.com/jenrik/golem"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestAttributeExtract(t *testing.T) {
	str := "<html><body><a href=\"example.com\">foo</div></body></html>"
	doc, _ := html.Parse(strings.NewReader(str))
	ctx := golem.Context{
		HTML: doc,
	}

	config := make(map[string]interface{})
	config["selector"] = "a"
	config["attribute"] = "href"

	action := new(ActionAttribute)
	conf, _ := action.ExpandConfig(nil, nil, config)
	output, err := action.Action(&ctx, conf)

	if err != nil {
		t.FailNow()
	}
	if val, ok := output.([]string); ok {
		if len(val) != 1 {
			t.FailNow()
		}
		if val[0] != "example.com" {
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

func TestAttributeValidateConfig(t *testing.T) {
	config := make(map[string]interface{})

	action := new(ActionAttribute)
	output := action.Validate(config)
	if output {
		t.FailNow()
	}

	config["selector"] = "a"
	config["attribute"] = "href"
	output = action.Validate(config)
	if !output {
		t.FailNow()
	}
}
