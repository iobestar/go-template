package internal

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestInterface interface {
	Runner
	Starter
}

type TestComponent struct {
	startCount atomic.Int64
	stopCount  atomic.Int64
	runCount   atomic.Int64
}

func (t *TestComponent) Run(ctx context.Context, wg *sync.WaitGroup, disposer Disposer) error {
	t.runCount.Add(1)
	return nil
}

func (t *TestComponent) Start() error {
	t.startCount.Add(1)
	return nil
}

func (t *TestComponent) Stop() error {
	t.stopCount.Add(1)
	return nil
}

func TestAppContainer_GetByName(t *testing.T) {
	appCnt := NewAppContainer()

	c := AppComponent{Order: 1, Name: "foo", Value: "bar"}

	err := appCnt.Register(c)
	require.NoError(t, err)

	r, ok := appCnt.GetByName("foo")
	require.True(t, ok)
	require.Equal(t, "bar", r.(string))
}

func TestAppContainer_GetAllByType(t *testing.T) {
	appCnt := NewAppContainer()

	c1 := AppComponent{Order: 1, Name: "foo", Value: "bar"}
	c2 := AppComponent{Order: 2, Name: "baz", Value: 123}
	c3 := AppComponent{Order: 3, Name: "qux", Value: "quux"}
	err := appCnt.Register(c1, c2, c3)
	require.NoError(t, err)

	results := appCnt.GetAllByType(reflect.TypeFor[string]())
	require.Len(t, results, 2)
	values := []string{results[0].(string), results[1].(string)}
	sort.Strings(values)
	require.Equal(t, []string{"bar", "quux"}, values)
}

func TestAppContainer_Start(t *testing.T) {
	appCnt := NewAppContainer()

	c1 := AppComponent{Order: 2, Name: "comp1", Value: &TestComponent{}}
	appCnt.MustRegister(c1)

	require.NoError(t, appCnt.Start())
	require.NoError(t, appCnt.Run(context.Background(), &sync.WaitGroup{}, NewDisposer()))

	assert.Equal(t, 1, int(c1.Value.(*TestComponent).startCount.Load()))
	assert.Equal(t, 1, int(c1.Value.(*TestComponent).runCount.Load()))
	require.NoError(t, appCnt.Stop())
	assert.Equal(t, 1, int(c1.Value.(*TestComponent).stopCount.Load()))
}

func TestAppContainer_RunOrder(t *testing.T) {
	appCnt := NewAppContainer()

	var components []*TestComponent

	for i := 0; i < 5; i++ {
		tc := &TestComponent{}
		components = append(components, tc)
		c := AppComponent{Order: i, Name: "comp" + string(rune('A'+i)), Value: tc}
		appCnt.MustRegister(c)
	}
}

func TestNewAppContainer_GetAllByTypeInterface(t *testing.T) {
	appCnt := NewAppContainer()

	c1 := AppComponent{Order: 2, Name: "comp1", Value: &TestComponent{}}
	appCnt.MustRegister(c1)

	results := appCnt.GetAllByType(reflect.TypeFor[TestInterface]())
	require.Len(t, results, 1)
	_, ok := results[0].(TestInterface)
	require.True(t, ok)
}
