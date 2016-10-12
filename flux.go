// Package flux is a Go implementation of the flux design pattern.
// https://facebook.github.io/flux/docs/overview.html
package flux

var (
	currentDispatcher = newDispatcher()
)

func init() {
	go currentDispatcher.Run()
}

// RegisterStore registers a store to receive dispatched actions.
func RegisterStore(s Store) {
	currentDispatcher.RegisterStore(s)
}

// UnregisterStore unregisters a store from the dispatcher.
// The store will no longer receive dispatched actions.
func UnregisterStore(s Store) {
	currentDispatcher.UnregisterStore(s)
}

// Dispatch dispatches an action by an action identifier
// to all the registered stores.
func Dispatch(a ActionID) {
	DispatchWithPayload(a, nil)
}

// DispatchWithPayload dispatches an action with a payload by an action identifier
// to all the registered stores.
func DispatchWithPayload(a ActionID, payload interface{}) {
	DispatchAction(Action{
		ID:      a,
		Payload: payload,
	})
}

// DispatchAction dispatches an action to all the registered stores.
func DispatchAction(a Action) {
	currentDispatcher.DispChan <- a
}
