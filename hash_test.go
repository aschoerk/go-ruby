package main

import (
	"testing"

	"aschoerk.de/go-ruby/ruby"
	"github.com/stretchr/testify/assert"
)

// maps, functions as keys not implemented yet
// CompareGenerally should be replacable using reflect.DeepEqual

func testWithSingleValues[K any, V any](t *testing.T, k K, v V) {
	h := ruby.NewHash[K, V]()
	assert.Equal(t, 0, h.Len())
	h.Set(k, v)
	assert.Equal(t, 1, h.Len())
	val, exists := h.Get(k)
	assert.Equal(t, 1, h.Len())
	assert.True(t, exists)
	assert.Equal(t, v, val)
	assert.Equal(t, 1, h.Len())
	h.Set(k, v)
	val, exists = h.Get(k)
	assert.True(t, exists)
	assert.Equal(t, v, val)
	assert.Equal(t, 1, h.Len())
	h.Delete(k)
	assert.Equal(t, 0, h.Len())
	_, exists = h.Get(k)
	assert.False(t, exists)

}

func TestWithSingleValues(t *testing.T) {
	testWithSingleValues(t, 1, "string")
	testWithSingleValues(t, "string", 1)
	testWithSingleValues(t, 1, 1)
	testWithSingleValues(t, 1.0, 1.0)
	testWithSingleValues(t, 1.0000000000000000001, 1.0)
	testWithSingleValues(t, []int{1}, 1)
	testWithSingleValues(t, 1, []int{1})
	testWithSingleValues(t, [][]int{{1}}, 1)
	testWithSingleValues(t, 1, [][]int{{1}})
	testWithSingleValues(t, [][]*[]int{{&[]int{1}}}, 1)
	testWithSingleValues(t, 1, [][]*[]int{{&[]int{1}}})
}

func testWithTwoPairs[K any, V any](t *testing.T, k1 K, v1 V, k2 K, v2 V, num int) {
	h := ruby.NewHash[K, V]()
	h.Set(k1, v1)
	h.Set(k2, v2)
	assert.Equal(t, num, h.Len())
	if num == 2 {
		val, exists := h.Get(k1)
		assert.Equal(t, num, h.Len())
		assert.True(t, exists)
		assert.Equal(t, v1, val)
	} else {
		// k1 was overwritten
		val, exists := h.Get(k1)
		assert.Equal(t, num, h.Len())
		assert.True(t, exists)
		assert.Equal(t, v2, val)
	}
	val, exists := h.Get(k2)
	assert.True(t, exists)
	assert.Equal(t, v2, val)
	h.Delete(k1)
	assert.Equal(t, num-1, h.Len())
	h.Delete(k1)
	assert.Equal(t, num-1, h.Len())
	_, exists = h.Get(k1)
	assert.False(t, exists)
	if num == 2 {
		h.Delete(k2)
		assert.Equal(t, 0, h.Len())
		h.Delete(k2)
		assert.Equal(t, 0, h.Len())
		_, exists = h.Get(k2)
		assert.False(t, exists)
	}
}

func TestWithTwoPairs(t *testing.T) {
	testWithTwoPairs(t, 1, "string", 1, "string", 1)
	testWithTwoPairs(t, "string", 1, "string", 1, 1)
	testWithTwoPairs(t, 1, 1, 1, 1, 1)
	testWithTwoPairs(t, 1.0, 1.0, 1.0, 1.0, 1)
	testWithTwoPairs(t, 1.0000000000000000001, 1.0, 1.0000000000000000001, 1.0, 1)
	testWithTwoPairs(t, []int{1}, 1, []int{1}, 1, 1)
	testWithTwoPairs(t, 1, []int{1}, 1, []int{1}, 1)
	testWithTwoPairs(t, [][]int{{1}}, 1, [][]int{{1}}, 1, 1)
	testWithTwoPairs(t, 1, [][]int{{1}}, 1, [][]int{{1}}, 1)
	testWithTwoPairs(t, [][]*[]int{{&[]int{1}}}, 1, [][]*[]int{{&[]int{1}}}, 1, 1)
	testWithTwoPairs(t, 1, [][]*[]int{{&[]int{1}}}, 1, [][]*[]int{{&[]int{1}}}, 1)

	testWithTwoPairs(t, 1, "string", 1, "strng", 1)
	testWithTwoPairs(t, "string", 1, "string", 2, 1)
	testWithTwoPairs(t, 1, 1, 1, 2, 1)
	testWithTwoPairs(t, 1.0, 1.0, 1.0, 1.1, 1)
	testWithTwoPairs(t, 1.0000000000000000001, 1.0, 1.0000000000000000002, 1.0, 1)
	testWithTwoPairs(t, []int{1}, 1, []int{1}, 2, 1)
	testWithTwoPairs(t, 1, []int{1}, 1, []int{2}, 1)
	testWithTwoPairs(t, [][]int{{1}}, 1, [][]int{{1}}, 2, 1)
	testWithTwoPairs(t, 1, [][]int{{1}}, 1, [][]int{{2}}, 1)
	testWithTwoPairs(t, [][]*[]int{{&[]int{1}}}, 1, [][]*[]int{{&[]int{1}}}, 2, 1)
	testWithTwoPairs(t, 1, [][]*[]int{{&[]int{1}}}, 1, [][]*[]int{{&[]int{2}}}, 1)

	testWithTwoPairs(t, 1, "string", 2, "strng", 2)
	testWithTwoPairs(t, "string", 1, "strng", 2, 2)
	testWithTwoPairs(t, 1, 1, 3, 2, 2)
	testWithTwoPairs(t, 1.0, 1.0, 3.0, 1.1, 2)
	// testWithTwoPairs(t, 1.0000000000000000001, 1.0, 1.0000000000000000002, 1.0, 2)
	testWithTwoPairs(t, []int{1}, 1, []int{2}, 2, 2)
	testWithTwoPairs(t, 1, []int{1}, 2, []int{2}, 2)
	testWithTwoPairs(t, [][]int{{1}}, 1, [][]int{{2}}, 2, 2)
	testWithTwoPairs(t, 1, [][]int{{1}}, 2, [][]int{{2}}, 2)
	testWithTwoPairs(t, [][]*[]int{{&[]int{1}}}, 1, [][]*[]int{{&[]int{2}}}, 2, 2)
	testWithTwoPairs(t, 1, [][]*[]int{{&[]int{1}}}, 2, [][]*[]int{{&[]int{2}}}, 2)
}
