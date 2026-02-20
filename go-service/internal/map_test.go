package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLinkedMap(t *testing.T) {
	m := NewLinkedMap[string, int]()
	require.Equal(t, 0, m.Size())
	require.NotNil(t, m.items)
}

func TestNewLinkedMapWithCapacity(t *testing.T) {
	m := NewLinkedMapWithCapacity[string, int](10)
	require.Equal(t, 0, m.Size())
	require.NotNil(t, m.items)
}

func TestLinkedMap_Put_NewKey(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)

	require.Equal(t, 1, m.Size())
	val, exists := m.Get("key1")
	require.True(t, exists)
	require.Equal(t, 1, val)
}

func TestLinkedMap_Put_UpdateExisting(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key1", 2)

	require.Equal(t, 1, m.Size())
	val, exists := m.Get("key1")
	require.True(t, exists)
	require.Equal(t, 2, val)
}

func TestLinkedMap_Put_MultipleKeys(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key3", 3)

	require.Equal(t, 3, m.Size())

	// Verify all keys exist with correct values
	val1, exists1 := m.Get("key1")
	require.True(t, exists1)
	require.Equal(t, 1, val1)

	val2, exists2 := m.Get("key2")
	require.True(t, exists2)
	require.Equal(t, 2, val2)

	val3, exists3 := m.Get("key3")
	require.True(t, exists3)
	require.Equal(t, 3, val3)
}

func TestLinkedMap_Get_NonExistent(t *testing.T) {
	m := NewLinkedMap[string, int]()
	val, exists := m.Get("nonexistent")
	require.False(t, exists)
	require.Equal(t, 0, val) // Zero value
}

func TestLinkedMap_Delete_Existing(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)

	deleted := m.Delete("key1")
	require.True(t, deleted)
	require.Equal(t, 0, m.Size())

	_, exists := m.Get("key1")
	require.False(t, exists)
}

func TestLinkedMap_Delete_NonExistent(t *testing.T) {
	m := NewLinkedMap[string, int]()
	deleted := m.Delete("nonexistent")
	require.False(t, deleted)
	require.Equal(t, 0, m.Size())
}

func TestLinkedMap_Delete_Head(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key3", 3)

	deleted := m.Delete("key1")
	require.True(t, deleted)
	require.Equal(t, 2, m.Size())

	// Verify order is maintained
	keys := m.Keys()
	require.Equal(t, []string{"key2", "key3"}, keys)
}

func TestLinkedMap_Delete_Tail(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key3", 3)

	deleted := m.Delete("key3")
	require.True(t, deleted)
	require.Equal(t, 2, m.Size())

	// Verify order is maintained
	keys := m.Keys()
	require.Equal(t, []string{"key1", "key2"}, keys)
}

func TestLinkedMap_Delete_Middle(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key3", 3)

	deleted := m.Delete("key2")
	require.True(t, deleted)
	require.Equal(t, 2, m.Size())

	// Verify order is maintained
	keys := m.Keys()
	require.Equal(t, []string{"key1", "key3"}, keys)
}

func TestLinkedMap_Delete_SingleElement(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)

	deleted := m.Delete("key1")
	require.True(t, deleted)
	require.Equal(t, 0, m.Size())
	require.Nil(t, m.head)
	require.Nil(t, m.tail)
}

func TestLinkedMap_Contains(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)

	require.True(t, m.Contains("key1"))
	require.False(t, m.Contains("nonexistent"))
}

func TestLinkedMap_Size(t *testing.T) {
	m := NewLinkedMap[string, int]()
	require.Equal(t, 0, m.Size())

	m.Put("key1", 1)
	require.Equal(t, 1, m.Size())

	m.Put("key2", 2)
	require.Equal(t, 2, m.Size())

	m.Delete("key1")
	require.Equal(t, 1, m.Size())

	m.Delete("key2")
	require.Equal(t, 0, m.Size())
}

func TestLinkedMap_Keys_InsertionOrder(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	keys := m.Keys()
	require.Equal(t, []string{"c", "a", "b"}, keys)
}

func TestLinkedMap_Keys_Empty(t *testing.T) {
	m := NewLinkedMap[string, int]()
	keys := m.Keys()
	require.Equal(t, []string{}, keys)
}

func TestLinkedMap_Keys_AfterUpdate(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key1", 10) // Update existing

	// Order should not change when updating
	keys := m.Keys()
	require.Equal(t, []string{"key1", "key2"}, keys)
}

func TestLinkedMap_Values_InsertionOrder(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	values := m.Values()
	require.Equal(t, []int{3, 1, 2}, values)
}

func TestLinkedMap_Values_Empty(t *testing.T) {
	m := NewLinkedMap[string, int]()
	values := m.Values()
	require.Equal(t, []int{}, values)
}

func TestLinkedMap_Range_FullIteration(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key3", 3)

	var keys []string
	var values []int

	m.Range(func(k string, v int) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})

	require.Equal(t, []string{"key1", "key2", "key3"}, keys)
	require.Equal(t, []int{1, 2, 3}, values)
}

func TestLinkedMap_Range_EarlyExit(t *testing.T) {
	m := NewLinkedMap[string, int]()
	m.Put("key1", 1)
	m.Put("key2", 2)
	m.Put("key3", 3)

	var keys []string

	m.Range(func(k string, v int) bool {
		keys = append(keys, k)
		return v < 2 // Stop when value is 2
	})

	require.Equal(t, []string{"key1", "key2"}, keys)
}

func TestLinkedMap_Range_Empty(t *testing.T) {
	m := NewLinkedMap[string, int]()

	called := false
	m.Range(func(k string, v int) bool {
		called = true
		return true
	})

	require.False(t, called)
}

func TestLinkedMap_ComplexScenario(t *testing.T) {
	m := NewLinkedMap[string, string]()

	// Insert
	m.Put("user1", "Alice")
	m.Put("user2", "Bob")
	m.Put("user3", "Charlie")

	// Update
	m.Put("user2", "Bobby")

	// Delete
	m.Delete("user3")

	// Add more
	m.Put("user4", "David")

	// Verify size
	require.Equal(t, 3, m.Size())

	// Verify order and values
	keys := m.Keys()
	require.Equal(t, []string{"user1", "user2", "user4"}, keys)

	values := m.Values()
	require.Equal(t, []string{"Alice", "Bobby", "David"}, values)

	// Verify specific gets
	val, exists := m.Get("user2")
	require.True(t, exists)
	require.Equal(t, "Bobby", val)

	// Verify deleted key doesn't exist
	_, exists = m.Get("user3")
	require.False(t, exists)
}

func TestLinkedMap_IntKeys(t *testing.T) {
	m := NewLinkedMap[int, string]()
	m.Put(3, "three")
	m.Put(1, "one")
	m.Put(2, "two")

	keys := m.Keys()
	require.Equal(t, []int{3, 1, 2}, keys)

	val, exists := m.Get(1)
	require.True(t, exists)
	require.Equal(t, "one", val)
}

func TestLinkedMap_StructValues(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	m := NewLinkedMap[string, Person]()
	m.Put("p1", Person{Name: "Alice", Age: 30})
	m.Put("p2", Person{Name: "Bob", Age: 25})

	person, exists := m.Get("p1")
	require.True(t, exists)
	require.Equal(t, "Alice", person.Name)
	require.Equal(t, 30, person.Age)

	values := m.Values()
	require.Len(t, values, 2)
	require.Equal(t, "Alice", values[0].Name)
	require.Equal(t, "Bob", values[1].Name)
}
