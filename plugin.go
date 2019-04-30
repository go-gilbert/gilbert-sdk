package sdk

// Actions is actions map
//
// Key is an action name and value is action constructor
type Actions map[string]ActionFactory

// Plugin represents plugin structure
type Plugin struct {
	Name    string
	Actions Actions
}
