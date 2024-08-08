package ruby

import (
	"github.com/mitchellh/hashstructure/v2"
)

type hashImpl[K any, V any] struct {
	data map[any]V
	// map of hash-values of keys to pointers of keys which have the same hash but different value according to CompareGenerally
	uncomparable map[uint64][]any
}

// NewCustomMap creates a new CustomMap
func NewHash[K any, V any]() Hash[K, V] {
	return &hashImpl[K, V]{
		data:         make(map[any]V),
		uncomparable: make(map[uint64][]any),
	}
}

// wrapKey wraps non-comparable keys in a pointer
func checkKey[K any](key K) (bool, uint64) {
	if !isComparable(key) {
		hash, _ := hashstructure.Hash(key, hashstructure.FormatV2, nil)
		return false, hash
	}
	return true, 0
}

// keys can be equal, but can have different addresses in memory
// so the first key used for adding an entry to the map is representative. To do later changes for that key-value
// using the correct pointer, the original is returned as 3rd parameter.
// the hash is returned to be able to do updates without needing to recalculate the hash.
func (m *hashImpl[K, V]) getIncludingOriginalKey(key K) (V, bool, *K, uint64) {
	isComparable, hash := checkKey(key)
	if isComparable {
		value, exists := m.data[key]
		return value, exists, nil, hash
	} else {
		value, exists := m.data[&key]
		if exists {
			return value, exists, &key, hash
		}
		keys, keysExist := m.uncomparable[hash]
		if keysExist {
			for _, k := range keys {
				if CompareGenerally(any(&key), k) && CompareGenerally(k, any(&key)) {
					value2, exists2 := m.data[k]
					return value2, exists2, k.(*K), hash
				}
			}
		}
		return value, false, nil, hash
	}
}

// Set adds or updates a key-value pair in the map
func (m *hashImpl[K, V]) Set(key K, value V) {
	_, _, original, hash := m.getIncludingOriginalKey(key)
	if hash == 0 {
		m.data[key] = value
	} else {
		// handle uncomparable
		if original != nil {
			// equal key was already found, use original to set new value
			m.data[original] = value
		} else {
			// equal key was not found yet
			m.data[&key] = value
			// maintain map of hash of the key to be able to find equal candidates later
			keyslice, exists := m.uncomparable[hash]
			if !exists {
				keyslice = make([]any, 0)
			}
			keyslice = append(keyslice, &key)
			m.uncomparable[hash] = keyslice
		}
	}
}

// Get retrieves a value by key from the map
func (m *hashImpl[K, V]) Get(key K) (V, bool) {
	v, exists, _, _ := m.getIncludingOriginalKey(key)
	return v, exists
}

// Delete removes a key-value pair from the map
func (m *hashImpl[K, V]) Delete(key K) {
	_, _, original, hash := m.getIncludingOriginalKey(key)
	if hash == 0 {
		delete(m.data, key)
	} else {
		if original != nil {
			delete(m.data, original)
		}
		keys, exists := m.uncomparable[hash]
		if exists {
			if len(keys) == 0 {
				delete(m.uncomparable, hash)
			} else {
				keyslice := make([]any, 0)
				for _, el := range keys {
					if el != original {
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

func (m *hashImpl[K, V]) Len() int {
	return len(m.data)
}
