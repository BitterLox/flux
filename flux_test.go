package flux

import "testing"

func TestRegisterUnregisterStore(t *testing.T) {
	s := newStoreTest()

	RegisterStore(s)
	defer UnregisterStore(s)
}

func TestDispatch(t *testing.T) {
	Dispatch("Test flux dispatch")
}

func TestDispatchWithPayload(t *testing.T) {
	DispatchWithPayload("Test flux dispatch with payload", "Bancs publics")
}

func TestDispatchAction(t *testing.T) {
	DispatchAction(Action{
		ID:      "Test dlux dispatch action",
		Payload: 42,
	})
}
