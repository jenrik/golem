package action

import (
	"github.com/andybalholm/cascadia"
	"github.com/jenrik/golem"
)

type ActionAttribute struct{}

func (action *ActionAttribute) Action(ctx *golem.Context, config interface{}) (interface{}, error) {
	conf := config.(map[string]interface{})
	selector, err := cascadia.Compile(conf["selector"].(string))
	if err != nil {
		return nil, err
	}

	nodes := selector.MatchAll(ctx.HTML)
	if nodes == nil {
		return nil, nil
	}
	values := make([]string, len(nodes))

	for i, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == conf["attribute"].(string) {
				values[i] = attr.Val
				break
			}
		}
	}
	return values, nil
}

func (action *ActionAttribute) Validate(config interface{}) bool {
	conf := config.(map[string]interface{})
	selectorVal, selector := conf["selector"]
	_, selectorType := selectorVal.(string)
	attributeVal, attribute := conf["attribute"]
	_, attributeType := attributeVal.(string)
	return selector && selectorType && attribute && attributeType
}

func (action *ActionAttribute) ExpandConfig(stages *[]map[string]interface{}, ctx *golem.Context, config interface{}) (interface{}, error) {
	return config, nil
}

func init() {
	golem.RegisterAction(new(ActionAttribute), "attribute")
}
