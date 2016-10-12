package flux

import (
	"fmt"
	"testing"
)

var (
	SimpleListener = func(e Event) {
		fmt.Println("SimpleListener called")
	}
)

func TestNewStoreBase(t *testing.T) {
	NewStoreBase()
}

func TestStoreBaseID(t *testing.T) {
	s := NewStoreBase()
	s.SetID(42)

	if id := s.ID(); id != StoreID(42) {
		t.Error("id should be 42:", id)
	}
}

func TestStoreBaseAddListener(t *testing.T) {
	s := NewStoreBase()

	s.AddListener(SimpleListener)

	if l := len(s.listeners); l != 1 {
		t.Error("l should be 1:", l)
	}
}

func TestStoreBaseRemoveListener(t *testing.T) {
	s := NewStoreBase()

	lid := s.AddListener(SimpleListener)
	s.AddListener(SimpleListener)

	s.RemoveListener(lid)

	if l := len(s.listeners); l != 1 {
		t.Fatal("l should be 1:", l)
	}
}

func TestStoreBaseEmit(t *testing.T) {
	s := NewStoreBase()
	c := make(chan bool)

	l := func(e Event) {
		t.Log("l ->", e)
		c <- true
	}

	l2 := func(e Event) {
		t.Log("l2 ->", e)
		c <- true
	}

	s.AddListener(l)
	s.AddListener(l2)

	go s.Emit("test")
	go s.EmitWithPayload("test", "j'aime les filles")

	<-c
	<-c
	<-c
	<-c
}
