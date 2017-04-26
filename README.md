# flux
[![Build Status](https://travis-ci.org/murlokswarm/flux.svg?branch=master)](https://travis-ci.org/murlokswarm/flux)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/flux)](https://goreportcard.com/report/github.com/murlokswarm/flux)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/flux/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/flux?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/flux?status.svg)](https://godoc.org/github.com/murlokswarm/flux)

Package flux is a Go implementation of the [Flux](https://facebook.github.io/flux/docs/overview.html) design pattern.

## Install
```
go get -u github.com/murlokswarm/flux
```

## How to use?

### Create and register a store
```go
// Implemtation.
type HelloStore struct {
	flux.Store
}

func (s *HelloStore) OnDispatch(a flux.Action) {
	if a.Name != "greet" {
		return
	}

	s.Emit(flux.Event{
		Name:    "greeted",
		Payload: fmt.Sprintf("Hello, %v", a.Payload),
	})
}

// Intialization and registration.
var (
	helloStore flux.Storer
)

func init() {
	helloStore := &HelloStore{}
	flux.Register(helloStore)

}

```

### Create and register a view
```go
type HelloView struct {
	Greeting string
}

func (v *HelloView) OnMount() {
    // Listen events from helloStore.
	helloStore.Register(v)
}

func (v *HelloView) OnDismount() {
    // Stop listening events from helloStore.
    // Avoid memory leak.
	helloStore.Unregister(v)
}

func (v *HelloView) Render() string {
	return `
<div>
    <h1>{{html .Greeting}}</h1>
    <input onchange="OnInputChange" />
</div>
    `
}

func (v *HelloView) OnInputChange(a app.ChangeArg) {
    // Dispatch an action.
	flux.Dispatch(flux.Action{
		Name:    "greet",
		Payload: a.Value,
	})
}

func (v *HelloView) OnStoreEvent(e flux.Event) {
	if e.Name != "greeted" {
		return
	}

	// Handling events from helloStore.
	if greet, ok := e.Payload.(string); ok{
		v.Gretting = greet
		app.Render(v)
	}
}
```


