package internal

import (
	"context"
	"sync"
)

type DisposeFunc func()

type Disposer interface {
	Attach(DisposeFunc)
	Dispose()
}

type disposer struct {
	lock      sync.Mutex
	disposers []DisposeFunc
}

func NewDisposer() Disposer {
	return &disposer{
		lock:      sync.Mutex{},
		disposers: []DisposeFunc{},
	}
}

func (d *disposer) Dispose() {
	d.lock.Lock()
	defer d.lock.Unlock()
	for i := len(d.disposers) - 1; i >= 0; i-- {
		d.disposers[i]()
	}
}

func (d *disposer) Attach(disposer DisposeFunc) {
	if nil != disposer {
		d.lock.Lock()
		defer d.lock.Unlock()
		d.disposers = append(d.disposers, disposer)
	}
}

type Runner interface {
	Run(ctx context.Context, wg *sync.WaitGroup, disposer Disposer) error
}

type Starter interface {
	Start() error
	Stop() error
}
