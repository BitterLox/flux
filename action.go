package flux

// Action represents an action to be dispatched.
type Action struct {
	ID      ActionID
	Payload interface{}
}

// ActionID represents an identifier for an action.
// Every action should have an unique identifier.
type ActionID string
