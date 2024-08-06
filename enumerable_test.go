package main

import (
	"cmp"
	"testing"

	"aschoerk.de/go-ruby/ruby"
	"github.com/stretchr/testify/assert"
)

func TestEach(t *testing.T) {

	e := ruby.RStepped[int](1, 11, 1)
	var count int = 0
	e.Each(func(el int) {
		count++
		assert.Equal(t, el, count, "should be equal")
		lessOrEqual := func(a int, b int) bool { return a <= b }
		assert.True(t, e.Includes(el, lessOrEqual))
		assert.False(t, e.Includes(el+10, lessOrEqual))
	})

}

func TestEachWithIndex(t *testing.T) {

	e := ruby.RStepped[int](1, 11, 1)
	var count int = 0
	e.EachWithIndex(func(index int, el int) {
		count++
		if el != count {
			t.Errorf("Expected %d, but got %d", count, el)
		}
		if index != count-1 {
			t.Errorf("Expected index %d, but got %d", count-1, index)
		}
	})

}

func TestSlice(t *testing.T) {
	e := ruby.E(ruby.R[int](1, 3).Entries())
	if e.Count() != 2 {
		t.Errorf("Expected only 2 entries, but was %d", e.Count())
	}
	if e.Count(func(a int) bool { return a > 0 }) != 2 {
		t.Errorf("Expected only 2 entries, but was %d", e.Count())
	}
}

func LessEqualFunc[T cmp.Ordered](x T) func(T) bool {
	return func(y T) bool {
		return y <= x
	}
}

func TestString(t *testing.T) {
	e := ruby.E([]string{"a", "b", "b"})
	assert.True(t, e.One(LessEqualFunc("a")))
	assert.False(t, e.One(LessEqualFunc("b")))

}
