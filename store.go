package flux

import (
	"reflect"
	"sync"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
)

var (
	// Mutable vars related to store management, the mutex must be locked to
	// access them concurrently.
	storesMutex sync.Mutex
	stores      []Storer
)

// Storer defines the interface to implement a store.
// Stores are subject to concurrency.
// Implementations should take this in consideration and protect mutable fields.
type Storer interface {
	// Once a store is registered, OnDispatch will be called for every dispatched
	// actions.
	OnDispatch(a Action)

	Register(l Listener)

	Unregister(l Listener)

	Emit(e Event)
}

// Register registers s as dispatch target.
// Does nothing if s is already registered.
// Panic if s is not a pointer.
func Register(s Storer) {
	if v := reflect.ValueOf(s); v.Kind() != reflect.Ptr {
		log.Panicf("s is not a pointer: %T", s)
	}

	storesMutex.Lock()
	defer storesMutex.Unlock()

	for _, store := range stores {
		if store == s {
			return
		}
	}

	stores = append(stores, s)
}

// Unregister removes s from dispatch targets.
// Does nothing if s is not registered.
func Unregister(s Storer) {
	storesMutex.Lock()
	defer storesMutex.Unlock()

	for i, store := range stores {
		if store == s {
			copy(stores[i:], stores[i+1:])
			stores[len(stores)-1] = nil
			stores = stores[:len(stores)-1]
			return
		}
	}
}

// Listener describes a listener.
// Listeners should be implemented as pointers.
type Listener interface {
	OnStoreEvent(e Event)
}

// Event represents data to be passed to a listener.
type Event struct {
	Name    string
	Payload interface{}
	Error   error
}

// Store implements logic to register/unregister a listener and emit events.
// Should be embedded in Storer implementations.
type Store struct {
	mutex     sync.Mutex
	listeners []Listener
}

// Register registers l for event emissions.
// Does nothing if l is already registered.
// Panic if l is not a pointer.
func (s *Store) Register(l Listener) {
	if val := reflect.ValueOf(l); val.Kind() != reflect.Ptr {
		log.Panicf("l is not a pointer: %T", l)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, listener := range s.listeners {
		if l == listener {
			return
		}
	}

	s.listeners = append(s.listeners, l)
}

// Unregister removes l from event emissions.
// Does nothing if l is not registered.
func (s *Store) Unregister(l Listener) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, listener := range s.listeners {
		if l == listener {
			copy(s.listeners[i:], s.listeners[i+1:])
			s.listeners[len(s.listeners)-1] = nil
			s.listeners = s.listeners[:len(s.listeners)-1]
			return
		}
	}
}

// Emit emits an event.
// Calls OnEvent method from all registered listeners.
// All emissions are guaranteed to run the app UI goroutine.
func (s *Store) Emit(e Event) {
	s.mutex.Lock()
	listeners := make([]Listener, len(s.listeners))
	copy(listeners, s.listeners)
	s.mutex.Unlock()

	for _, l := range listeners {
		app.UIChan <- func() {
			l.OnStoreEvent(e)
		}
	}
}
