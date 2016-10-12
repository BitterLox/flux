package flux

import "sync"

// Store represents a store in the flux design pattern.
type Store interface {
	OnDispatch(a Action)

	ID() StoreID

	SetID(ID StoreID)
}

// StoreID represents a store identifier.
type StoreID int

// StoreBase is the base struct which should be embedded in every store implementation.
type StoreBase struct {
	id             StoreID
	listeners      map[ListenerID]Listener
	lastListenerID ListenerID
	mtx            sync.Mutex
}

// ID returns the store ID.
func (s *StoreBase) ID() StoreID {
	return s.id
}

// SetID sets the ID of the store.
// Used internally by the dispatcher.
// Should not be called.
func (s *StoreBase) SetID(ID StoreID) {
	s.id = ID
}

// AddListener adds a listener to the store.
func (s *StoreBase) AddListener(l Listener) ListenerID {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.lastListenerID++
	id := s.lastListenerID

	s.listeners[id] = l
	return id
}

// RemoveListener removes a listener from the store.
func (s *StoreBase) RemoveListener(ID ListenerID) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.listeners, ID)
}

// Emit sends an event by an event identifier to all registered listeners.
func (s *StoreBase) Emit(ID EventID) {
	s.EmitEvent(Event{
		ID: ID,
	})
}

// EmitWithPayload sends an event with a payload by
// an event identifier to all registered listeners.
func (s *StoreBase) EmitWithPayload(ID EventID, payload interface{}) {
	s.EmitEvent(Event{
		ID:      ID,
		Payload: payload,
	})
}

// EmitEvent sends an event to all registered listeners.
func (s *StoreBase) EmitEvent(e Event) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, l := range s.listeners {
		l(e)
	}
}

// NewStoreBase creates a new instance of StoreBase.
func NewStoreBase() *StoreBase {
	return &StoreBase{
		listeners: map[ListenerID]Listener{},
	}
}
