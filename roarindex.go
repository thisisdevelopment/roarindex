package roarindex

import (
	"errors"
	"sync"

	roaring "github.com/RoaringBitmap/roaring"
)

// ErrKeyNotFound is returned when a key is not found in the RoarIndex.
var ErrKeyNotFound = errors.New("key not found")

// RoarIndex is a mapping from keys of type K to sets of values of type V.
// It uses RoaringBitmap internally for efficient storage and operations.
type RoarIndex[K comparable, V comparable] struct {
	mtx sync.RWMutex

	// Internal counters to assign unique IDs
	nextKeyID   uint32
	nextValueID uint32

	// Maps to assign unique IDs to keys and values
	keyToID   map[K]uint32
	idToKey   map[uint32]K
	valueToID map[V]uint32
	idToValue map[uint32]V

	// Map from key IDs to RoaringBitmap of value IDs
	data map[uint32]*roaring.Bitmap
}

// NewRoarIndex creates a new RoarIndex.
func NewRoarIndex[K comparable, V comparable]() *RoarIndex[K, V] {
	return &RoarIndex[K, V]{
		keyToID:   make(map[K]uint32),
		idToKey:   make(map[uint32]K),
		valueToID: make(map[V]uint32),
		idToValue: make(map[uint32]V),
		data:      make(map[uint32]*roaring.Bitmap),
	}
}

// PushMap associates a value with a key.
func (om *RoarIndex[K, V]) PushMap(key K, value V) {
	om.mtx.Lock()
	defer om.mtx.Unlock()

	// Get or assign key ID
	keyID, keyExists := om.keyToID[key]
	if !keyExists {
		keyID = om.nextKeyID
		om.nextKeyID++
		om.keyToID[key] = keyID
		om.idToKey[keyID] = key
	}

	// Get or assign value ID
	valueID, valueExists := om.valueToID[value]
	if !valueExists {
		valueID = om.nextValueID
		om.nextValueID++
		om.valueToID[value] = valueID
		om.idToValue[valueID] = value

	}

	// Get or create bitmap for the key
	bm, exists := om.data[keyID]
	if !exists {
		bm = roaring.NewBitmap()
		om.data[keyID] = bm
	}
	// Add the value ID to the bitmap
	bm.Add(valueID)
}

// GetMap retrieves the set of values associated with a key.
func (om *RoarIndex[K, V]) GetMap(key K) ([]V, error) {
	om.mtx.RLock()
	defer om.mtx.RUnlock()

	keyID, keyExists := om.keyToID[key]
	if !keyExists {
		return nil, ErrKeyNotFound
	}

	bm, exists := om.data[keyID]
	if !exists {
		return nil, nil // No values associated
	}

	values := make([]V, 0, bm.GetCardinality())
	it := bm.Iterator()
	for it.HasNext() {
		valueID := it.Next()
		value, valueExists := om.idToValue[valueID]
		if valueExists {
			values = append(values, value)
		}
	}

	return values, nil
}

// HasValue checks if a value is associated with a key.
func (om *RoarIndex[K, V]) HasValue(key K, value V) bool {
	om.mtx.RLock()
	defer om.mtx.RUnlock()

	keyID, keyExists := om.keyToID[key]
	if !keyExists {
		return false
	}

	valueID, valueExists := om.valueToID[value]
	if !valueExists {
		return false
	}

	bm, exists := om.data[keyID]
	if !exists {
		return false
	}

	return bm.Contains(valueID)
}

// DeleteMap removes a key and its associated values from the RoarIndex.
func (om *RoarIndex[K, V]) DeleteMap(key K) {
	om.mtx.Lock()
	defer om.mtx.Unlock()

	keyID, keyExists := om.keyToID[key]
	if !keyExists {
		return
	}

	// Remove the bitmap and key mappings
	delete(om.data, keyID)
	delete(om.keyToID, key)
	delete(om.idToKey, keyID)
}

// Keys returns a slice of all keys in the RoarIndex.
func (om *RoarIndex[K, V]) Keys() []K {
	om.mtx.RLock()
	defer om.mtx.RUnlock()

	keys := make([]K, 0, len(om.keyToID))
	for key := range om.keyToID {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a slice of all values in the RoarIndex.
func (om *RoarIndex[K, V]) Values() []V {
	om.mtx.RLock()
	defer om.mtx.RUnlock()

	values := make([]V, 0, len(om.valueToID))
	for value := range om.valueToID {
		values = append(values, value)
	}
	return values
}

// Count returns the number of keys in the RoarIndex.
func (om *RoarIndex[K, V]) Count() int {
	om.mtx.RLock()
	defer om.mtx.RUnlock()

	return len(om.keyToID)
}
