package flux

import (
	"fmt"
	"testing"
)

// ============================================================================
// StoreTest
// ============================================================================

const (
	TestAction ActionID = "test-action"
)

type StoreTest struct {
	*StoreBase
}

func (s *StoreTest) OnDispatch(a Action) {
	switch a.ID {
	case TestAction:
		fmt.Println("OnDispatch for", a)
	}
}

func newStoreTest() *StoreTest {
	return &StoreTest{
		StoreBase: NewStoreBase(),
	}
}

// ============================================================================
// Tests
// ============================================================================

func TestNewDispatcher(t *testing.T) {
	newDispatcher()
}

func TestDispatcherRegisterStore(t *testing.T) {
	d := newDispatcher()
	s := newStoreTest()
	d.RegisterStore(s)

	if l := len(d.stores); l != 1 {
		t.Error("l should be 1:", l)
	}
}

func TestDispatcherRegisterStoreRegistered(t *testing.T) {
	defer func() { recover() }()

	d := newDispatcher()
	s := newStoreTest()
	d.RegisterStore(s)
	d.RegisterStore(s)
	t.Error("should panic")
}

func TestDispatcherUnregisterStore(t *testing.T) {
	d := newDispatcher()
	s := newStoreTest()
	s2 := newStoreTest()
	d.RegisterStore(s)
	d.RegisterStore(s2)

	d.UnregisterStore(s)

	if l := len(d.stores); l != 1 {
		t.Fatal("l should be 1:", l)
	}

	if s3 := d.stores[s2.ID()]; s3 != s2 {
		t.Errorf("s3 should be s2: %v != %v", s3, s2)
	}

}

func TestDispatcherDispatch(t *testing.T) {
	d := newDispatcher()
	s := newStoreTest()
	d.RegisterStore(s)
	d.dispatch(Action{
		ID:      TestAction,
		Payload: "La vie en rose",
	})
}
