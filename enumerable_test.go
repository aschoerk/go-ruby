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
	if e.CountAll() != 2 {
		t.Errorf("Expected only 2 entries, but was %d", e.CountAll())
	}
	if e.Count(func(a int) bool { return a > 0 }) != 2 {
		t.Errorf("Expected only 2 entries, but was %d", e.CountAll())
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

func TestTally(t *testing.T) {
	e := ruby.E([]string{"a", "b", "b"})
	tally := e.Tally()
	numA, exists := tally.Get("a")
	assert.True(t, exists)
	assert.Equal(t, numA, 1)
	numA, exists = tally.Get("b")
	assert.True(t, exists)
	assert.Equal(t, numA, 2)
}

func TestTallyWithNotComparables(t *testing.T) {
	e := ruby.E([][]string{[]string{"a"}, []string{"b"}, []string{"b"}})
	tally := e.Tally()
	numA, exists := tally.Get([]string{"a"})
	assert.True(t, exists)
	assert.Equal(t, 1, numA)
	numA, exists = tally.Get([]string{"b"})
	assert.True(t, exists)
	assert.Equal(t, 2, numA)
}

func TestFilter(t *testing.T) {
	e := ruby.E([]string{"a", "b", "b"})
	assert.Equal(t, 1, e.Filter(ruby.GenEquals("a")).CountAll())
	assert.Equal(t, 2, e.Filter(ruby.GenEquals("b")).CountAll())
	assert.Equal(t, 0, e.Filter(ruby.GenEquals("c")).CountAll())
	assert.Equal(t, 3, e.Filter(ruby.GenTrue[string]()).CountAll())
	assert.Equal(t, 0, e.Filter(ruby.GenFalse[string]()).CountAll())
	assert.Equal(t, 1,
		e.Filter(ruby.GenEquals("a")).
			Filter(ruby.GenEquals("a")).CountAll())
	assert.Equal(t, 2,
		e.Filter(ruby.GenEquals("b")).
			Filter(ruby.GenEquals("b")).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenEquals("c")).
			Filter(ruby.GenEquals("c")).CountAll())
	assert.Equal(t, 3,
		e.Filter(ruby.GenTrue[string]()).
			Filter(ruby.GenTrue[string]()).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenFalse[string]()).
			Filter(ruby.GenFalse[string]()).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenEquals("a")).
			Filter(ruby.GenEquals("b")).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenEquals("b")).
			Filter(ruby.GenEquals("a")).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenEquals("c")).
			Filter(ruby.GenEquals("a")).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenTrue[string]()).
			Filter(ruby.GenFalse[string]()).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenFalse[string]()).
			Filter(ruby.GenTrue[string]()).CountAll())
	assert.Equal(t, 1,
		e.Filter(ruby.GenTrue[string]()).
			Filter(ruby.GenEquals("a")).CountAll())
	assert.Equal(t, 2,
		e.Filter(ruby.GenTrue[string]()).
			Filter(ruby.GenEquals("b")).CountAll())
	assert.Equal(t, 0,
		e.Filter(ruby.GenTrue[string]()).
			Filter(ruby.GenEquals("c")).CountAll())
}

func firstParam[T any](f func() (T, bool)) T {
	res, _ := f()
	return res
}
func existsParam[T any](f func() (T, bool)) bool {
	_, exists := f()
	return exists
}

func TestFetching(t *testing.T) {
	e := ruby.E([]string{"a", "b", "b"})
	assert.Equal(t, firstParam(e.First), "a")
	assert.Equal(t, firstParam(e.FirstN(1).First), "a")
	assert.False(t, existsParam(e.FirstN(0).First))
	assert.Equal(t, e.FirstN(0).CountAll(), 0)
	assert.Equal(t, firstParam(e.FirstN(2).First), "a")
	assert.Equal(t, e.FirstN(2).CountAll(), 2)
	assert.Equal(t, firstParam(e.FirstN(10).First), "a")
	assert.Equal(t, e.FirstN(10).CountAll(), 3)
	assert.Equal(t, e.TakeWhile(ruby.GenEquals("a")).CountAll(), 1)
	assert.Equal(t, e.TakeWhile(ruby.GenEquals("b")).CountAll(), 0)
	assert.Equal(t, ruby.E([]string{}).TakeWhile(ruby.GenEquals("b")).CountAll(), 0)
	assert.Equal(t, e.DropWhile(ruby.GenEquals("a")).CountAll(), 2)
	assert.Equal(t, e.DropWhile(ruby.GenEquals("b")).CountAll(), 3)
	assert.Equal(t, e.DropWhile(ruby.GenEquals("c")).CountAll(), 3)
	assert.Equal(t, ruby.E([]string{}).DropWhile(ruby.GenEquals("b")).CountAll(), 0)

}
