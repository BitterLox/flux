package flux

// Listener represents a function signature to listen events.
type Listener func(e Event)

// ListenerID represents a listener identifier.
type ListenerID int
