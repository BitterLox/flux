package flux

import (
	"sync"

	"github.com/murlokswarm/log"
)

type dispatcher struct {
	DispChan    chan Action
	stores      map[StoreID]Store
	lastStoreID StoreID
	mtx         sync.Mutex
}

func (d *dispatcher) RegisterStore(s Store) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.stores[s.ID()]; ok {
		log.Panicf("store with ID = %v is already registered", s.ID())
	}

	d.lastStoreID++
	id := d.lastStoreID
	s.SetID(id)
	d.stores[id] = s
}

func (d *dispatcher) UnregisterStore(s Store) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	delete(d.stores, s.ID())
}

func (d *dispatcher) Run() {
	defer close(d.DispChan)

	for {
		select {
		case action := <-d.DispChan:
			d.dispatch(action)
		}

	}
}

func (d *dispatcher) dispatch(a Action) {
	for _, store := range d.stores {
		store.OnDispatch(a)
	}
}

func newDispatcher() *dispatcher {
	return &dispatcher{
		DispChan: make(chan Action, 42),
		stores:   map[StoreID]Store{},
	}
}
