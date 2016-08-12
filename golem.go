package golem

type Pair struct {
	key   string
	value string
}

type Action interface {
	// Perform the action on the page
	Action(*Context, interface{}) (interface{}, error)

	// Expands any template string in the config and validate the expanded config
	ExpandConfig(*[]map[string]interface{}, *Context, interface{}) (interface{}, error)

	// Check if the supplied config is valid for the Action
	Validate(interface{}) bool
}

// Registered Actions
var actions map[string]Action = make(map[string]Action)

// Register a new action for use in definitions
// Return whether or not a existing action was overriden
func RegisterAction(action Action, name string) bool {
	_, exists := actions[name]
	actions[name] = action
	return exists
}
