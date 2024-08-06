package ruby

import (
	"reflect"

	"github.com/mitchellh/hashstructure/v2"
)

type hashImpl[K any, V any] struct {
	data         map[any]V
	uncomparable map[uint64][]any
}

// NewCustomMap creates a new CustomMap
func newHash[K any, V any]() *hashImpl[K, V] {
	return &hashImpl[K, V]{
		data: make(map[any]V),
	}
}

// wrapKey wraps non-comparable keys in a pointer
func checkKey[K any](key K) (bool, uint64) {
	if !isComparable(reflect.TypeOf(key)) {
		hash, _ := hashstructure.Hash(key, hashstructure.FormatV2, nil)
		return false, hash
	}
	return true, 0
}

// Set adds or updates a key-value pair in the map
func (m *hashImpl[K, V]) Set(key K, value V) {
	isComparable, hash := checkKey(key)
	if isComparable {
		m.data[key] = value
	} else {
		m.data[&key] = value
		keyslice, exists := m.uncomparable[hash]
		if !exists {
			keyslice = make([]any, 0)
		}
		keyslice = append(keyslice, &key)
		m.uncomparable[hash] = keyslice
	}
}

// Get retrieves a value by key from the map
func (m *hashImpl[K, V]) Get(key K) (V, bool) {
	isComparable, hash := checkKey(key)
	if isComparable {
		value, exists := m.data[key]
		return value, exists
	} else {
		value, exists := m.data[&key]
		if exists {
			return value, exists
		}
		keys, keysExist := m.uncomparable[hash]
		if keysExist {
			for _, k := range keys {
				if CompareGenerally(any(&key), k) && CompareGenerally(k, any(&key)) {
					value2, exists2 := m.data[k]
					return value2, exists2
				}
			}
		}
		return value, false
	}
}

// Delete removes a key-value pair from the map
func (m *hashImpl[K, V]) Delete(key K) {
	isComparable, hash := checkKey(key)
	if isComparable {
		delete(m.data, key)
	} else {
		delete(m.data, &key)
		keys, exists := m.uncomparable[hash]
		if exists {
			if len(keys) == 0 {
				delete(m.uncomparable, hash)
			} else {
				keyslice := make([]any, 0)
				for _, el := range keys {
					if !CompareGenerally(any(&key), el) || !CompareGenerally(el, any(&key)) {
						keyslice = append(keyslice, el)
					}
				}
				if len(keyslice) == 0 {
					delete(m.uncomparable, hash)
				} else {
					m.uncomparable[hash] = keyslice
				}
			}
		}
	}
}
