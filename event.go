package flux

// Event represents an event that would be sent to a listener.
type Event struct {
	ID      EventID
	Payload interface{}
}

// EventID represents an identifier for an event.
// Every event should have an unique identifier.
type EventID string
