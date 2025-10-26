package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"sync"
)

type App interface {
	Runner
	GetContext() AppContext
}

type HttpApp interface {
	App
	http.Handler
}

type AppContext interface {
	GetByName(name string) (any, bool)
	MustGetByName(name string) any
	GetAllByType(t reflect.Type) []any
	GetByType(t reflect.Type) (any, bool)
	MustGetByType(t reflect.Type) any
}

type AppContainer interface {
	Starter
	Runner
	AppContext
	Register(value ...any) error
	MustRegister(value ...any)
}

type AppComponent struct {
	Order int
	Name  string
	Value any
}

type defaultContainer struct {
	componentsByName map[string]AppComponent
	componentsByType map[reflect.Type][]AppComponent

	mux sync.RWMutex

	disposer Disposer
}

func NewAppContainer() AppContainer {
	return &defaultContainer{
		componentsByName: make(map[string]AppComponent),
		componentsByType: make(map[reflect.Type][]AppComponent),

		mux: sync.RWMutex{},

		disposer: NewDisposer(),
	}
}

func (c *defaultContainer) Register(value ...any) error {
	for _, v := range value {
		switch a := v.(type) {
		case AppComponent:
			if err := c.registerComponents(a); nil != err {
				return err
			}
		case DisposeFunc:
			c.disposer.Attach(a)
		default:
			return errors.New("value must be AppComponent or []AppComponent")
		}
	}
	return nil
}

func (c *defaultContainer) MustRegister(value ...any) {
	if err := c.Register(value...); err != nil {
		panic(err)
	}
}

func (c *defaultContainer) registerComponents(components ...AppComponent) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, component := range components {
		name := component.Name
		if _, exists := c.componentsByName[name]; exists {
			return errors.New(fmt.Sprintf("component already exists: %s", name))
		} else {
			c.componentsByName[name] = component
		}

		t := reflect.TypeOf(component.Value)

		if values, exists := c.componentsByType[t]; exists {
			c.componentsByType[t] = append(values, component)
		} else {
			c.componentsByType[t] = []AppComponent{component}
		}
	}
	return nil
}

// collectAndSortStarters returns all Starter components sorted by Order.
func (c *defaultContainer) collectAndSortStarters() []Starter {
	components := make([]AppComponent, 0, len(c.componentsByType))
	for _, component := range c.componentsByName {
		components = append(components, component)
	}

	return filterAndSortComponentsByType[Starter](components)
}

// collectAndSortRunners returns all Runner components sorted by Order.
func (c *defaultContainer) collectAndSortRunners() []Runner {
	components := make([]AppComponent, 0, len(c.componentsByType))
	for _, component := range c.componentsByName {
		components = append(components, component)
	}

	return filterAndSortComponentsByType[Runner](components)
}

func (c *defaultContainer) Start() error {
	c.mux.RLock()
	defer c.mux.RUnlock()

	starters := c.collectAndSortStarters()
	for _, starter := range starters {
		if err := starter.Start(); err != nil {
			return fmt.Errorf("failed to start component: %w", err)
		}
	}
	return nil
}

func (c *defaultContainer) Stop() error {
	c.mux.RLock()
	defer c.mux.RUnlock()

	starters := c.collectAndSortStarters()
	for i := len(starters) - 1; i >= 0; i-- {
		if err := starters[i].Stop(); err != nil {
			return fmt.Errorf("failed to stop component: %w", err)
		}
	}

	c.disposer.Dispose()
	return nil
}

func (c *defaultContainer) Run(ctx context.Context, wg *sync.WaitGroup, disposer Disposer) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, runner := range c.collectAndSortRunners() {
		if err := runner.Run(ctx, wg, disposer); nil != err {
			return err
		}
	}
	return nil
}

func (c *defaultContainer) GetByName(name string) (any, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if component, ok := c.componentsByName[name]; ok {
		return component.Value, true
	}
	return nil, false
}

func (c *defaultContainer) GetAllByType(t reflect.Type) []any {
	c.mux.RLock()
	defer c.mux.RUnlock()

	result := make([]any, 0)
	if components, ok := c.componentsByType[t]; ok {
		for _, component := range components {
			result = append(result, component.Value)
		}
	}

	for _, component := range c.componentsByName {
		if reflect.TypeOf(component.Value).AssignableTo(t) {
			result = append(result, component.Value)
		}
	}
	return result
}

func (c *defaultContainer) GetByType(t reflect.Type) (any, bool) {
	if components := c.GetAllByType(t); len(components) > 0 {
		return components[0], true
	}
	return nil, false
}

func (c *defaultContainer) MustGetByName(name string) any {
	if components, ok := c.GetByName(name); ok {
		return components
	} else {
		panic(fmt.Sprintf("component not found: %s", name))
	}
}

func (c *defaultContainer) MustGetByType(t reflect.Type) any {
	if component, ok := c.GetByType(t); ok {
		return component
	} else {
		panic(fmt.Sprintf("component not found for type: %s", t.String()))
	}
}

func filterAndSortComponentsByType[T any](components []AppComponent) []T {
	var filtered []AppComponent
	for _, c := range components {
		if _, ok := c.Value.(T); ok {
			filtered = append(filtered, c)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Order < filtered[j].Order
	})

	result := make([]T, len(filtered))
	for i, c := range filtered {
		result[i] = c.Value.(T)
	}
	return result
}

func GetAllByType[T any](ctx AppContext) []T {
	components := ctx.GetAllByType(reflect.TypeOf((*T)(nil)).Elem())
	result := make([]T, len(components))
	for i, c := range components {
		result[i] = c.(T)
	}
	return result
}

func MustGetByType[T any](ctx AppContext) T {
	if component, ok := ctx.GetByType(reflect.TypeOf((*T)(nil)).Elem()); ok {
		return component.(T)
	}
	panic(fmt.Sprintf("component not found for type: %s", reflect.TypeOf(ctx.GetAllByType(reflect.TypeOf((*T)(nil)).Elem()))))
}

func MustGetByName[T any](ctx AppContext, name string) T {
	if component, ok := ctx.GetByName(name); ok {
		return component.(T)
	}
	panic(fmt.Sprintf("component not found: %s", name))
}
