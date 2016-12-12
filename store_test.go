package flux

import (
	"testing"
	"time"
)

type StoreTest struct {
	Store
	OnDispatched bool
}

func (s *StoreTest) OnDispatch(a Action) {
	s.Emit(Event{
		Name:    "Success",
		Payload: 42,
	})
	s.OnDispatched = true
}

type BadStore struct {
	*Store
	OnDispatched bool
}

func (s BadStore) OnDispatch(a Action) {}

type ListenerTest struct {
	OnEventCalled bool
}

func (l *ListenerTest) OnStoreEvent(e Event) {
	l.OnEventCalled = true
}

type BadListener struct{}

func (l BadListener) OnStoreEvent(e Event) {}

func TestRegister(t *testing.T) {
	s := &StoreTest{}
	Register(s)

	if l := len(stores); l != 1 {
		t.Error("stores should have 1 element:", l)
	}

	Register(s)

	if l := len(stores); l != 1 {
		t.Error("stores should have 1 element:", l)
	}

	Unregister(s)

	if l := len(stores); l != 0 {
		t.Error("stores should be empty:", l)
	}
}

func TestRegisterPanic(t *testing.T) {
	defer func() { recover() }()

	s := BadStore{}
	Register(s)
	t.Error("should panic")
}

func TestStoreRegister(t *testing.T) {
	s := &Store{}
	l := &ListenerTest{}
	s.Register(l)

	if l := len(s.listeners); l != 1 {
		t.Error("s.listeners should have 1 element:", l)
	}

	s.Register(l)

	if l := len(s.listeners); l != 1 {
		t.Error("s.listeners should have 1 element:", l)
	}

	s.Unregister(l)

	if l := len(s.listeners); l != 0 {
		t.Error("s.listeners should be empty:", l)
	}
}

func TestStoreRegisterPanic(t *testing.T) {
	defer func() { recover() }()

	s := &Store{}
	l := BadListener{}
	s.Register(l)
	t.Error("should panic")
}

func TestStoreEmit(t *testing.T) {
	l := &ListenerTest{}
	s := &Store{}

	s.Register(l)
	s.Emit(Event{
		Name: "TestEmit",
	})
	time.Sleep(time.Millisecond * 1)

	if !l.OnEventCalled {
		t.Error("l.OnEventCalled should be true")
	}
}

func BenchmarkEmit(b *testing.B) {
	b.StopTimer()
	s := StoreTest{}

	for i := 0; i < 500; i++ {
		s.Register(&ListenerTest{})
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		s.Emit(Event{Name: "Benchmark", Payload: 42})
	}
}
