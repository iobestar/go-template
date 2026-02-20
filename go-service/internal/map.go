package internal

// entry represents a node in the doubly linked list
type entry[K comparable, V any] struct {
	key   K
	value V
	prev  *entry[K, V]
	next  *entry[K, V]
}

// LinkedMap maintains insertion order while providing O(1) access operations.
// It combines a hash map for fast lookups with a doubly linked list to preserve
// the order in which keys were inserted.
//
// This implementation is NOT thread-safe. If concurrent access is needed,
// external synchronization must be used.
type LinkedMap[K comparable, V any] struct {
	items map[K]*entry[K, V] // Hash map for O(1) lookups
	head  *entry[K, V]        // First entry in insertion order
	tail  *entry[K, V]        // Last entry in insertion order
	size  int                 // Number of entries
}

// NewLinkedMap creates a new empty LinkedMap.
func NewLinkedMap[K comparable, V any]() LinkedMap[K, V] {
	return LinkedMap[K, V]{
		items: make(map[K]*entry[K, V]),
	}
}

// NewLinkedMapWithCapacity creates a new LinkedMap with the specified initial capacity.
// This is useful when the approximate size is known in advance to avoid map resizing.
func NewLinkedMapWithCapacity[K comparable, V any](capacity int) LinkedMap[K, V] {
	return LinkedMap[K, V]{
		items: make(map[K]*entry[K, V], capacity),
	}
}

// Put inserts or updates a key-value pair in the map.
// If the key already exists, its value is updated without changing the insertion order.
// If the key is new, it is appended to the end of the insertion order.
//
// Time complexity: O(1) average case
func (m *LinkedMap[K, V]) Put(key K, value V) {
	if e, exists := m.items[key]; exists {
		// Key exists, just update the value
		e.value = value
		return
	}

	// Create new entry
	newEntry := &entry[K, V]{
		key:   key,
		value: value,
	}

	// Add to map
	m.items[key] = newEntry

	// Add to linked list (append to tail)
	if m.tail == nil {
		// First entry
		m.head = newEntry
		m.tail = newEntry
	} else {
		// Append to tail
		m.tail.next = newEntry
		newEntry.prev = m.tail
		m.tail = newEntry
	}

	m.size++
}

// Get retrieves the value associated with the given key.
// Returns the value and true if the key exists, or a zero value and false otherwise.
//
// Time complexity: O(1) average case
func (m *LinkedMap[K, V]) Get(key K) (V, bool) {
	if e, exists := m.items[key]; exists {
		return e.value, true
	}
	var zero V
	return zero, false
}

// Delete removes the key-value pair associated with the given key.
// Returns true if the key existed and was deleted, false otherwise.
//
// Time complexity: O(1) average case
func (m *LinkedMap[K, V]) Delete(key K) bool {
	e, exists := m.items[key]
	if !exists {
		return false
	}

	// Remove from map
	delete(m.items, key)

	// Remove from linked list
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		// Deleting head
		m.head = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	} else {
		// Deleting tail
		m.tail = e.prev
	}

	m.size--
	return true
}

// Contains checks whether the map contains the given key.
//
// Time complexity: O(1) average case
func (m *LinkedMap[K, V]) Contains(key K) bool {
	_, exists := m.items[key]
	return exists
}

// Size returns the number of key-value pairs in the map.
//
// Time complexity: O(1)
func (m *LinkedMap[K, V]) Size() int {
	return m.size
}

// Keys returns a slice of all keys in insertion order.
//
// Time complexity: O(n)
func (m *LinkedMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.size)
	for e := m.head; e != nil; e = e.next {
		keys = append(keys, e.key)
	}
	return keys
}

// Values returns a slice of all values in insertion order.
//
// Time complexity: O(n)
func (m *LinkedMap[K, V]) Values() []V {
	values := make([]V, 0, m.size)
	for e := m.head; e != nil; e = e.next {
		values = append(values, e.value)
	}
	return values
}

// Range iterates over all key-value pairs in insertion order.
// The iteration continues until all entries have been visited or the function returns false.
// If the function returns false, iteration stops early.
//
// Time complexity: O(n) for full iteration
func (m *LinkedMap[K, V]) Range(fn func(K, V) bool) {
	for e := m.head; e != nil; e = e.next {
		if !fn(e.key, e.value) {
			break
		}
	}
}
