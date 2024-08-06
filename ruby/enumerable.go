package ruby

import "reflect"

type enumerableImpl[T any] struct {
	EnumeratorGenerator[T]
}

type comparableEnumerableImpl[T comparable] struct {
	enumerableImpl[T]
}

func (e *enumerableImpl[T]) Each(f func(T)) {
	enumerator := e.EnumeratorGenerator.create()
	for enumerator.hasNext() {
		f(enumerator.next())
	}
}

func (e *enumerableImpl[T]) EachWithIndex(f func(int, T)) {
	enumerator := e.EnumeratorGenerator.create()
	i := 0
	for enumerator.hasNext() {
		f(i, enumerator.next())
		i++
	}
}

func getComparer[T any](funcs []func(T, T) bool) func(T, T) bool {
	if len(funcs) > 1 {
		panic("invalid call")
	}
	var comparer func(T, T) bool
	if len(funcs) == 0 {
		comparer = CompareGenerally
	} else {
		comparer = funcs[0]
	}
	return comparer
}

func (e *enumerableImpl[T]) Includes(t T, lessOrEqual ...func(T, T) bool) bool {
	comparer := getComparer(lessOrEqual)
	enumerator := e.EnumeratorGenerator.create()
	for enumerator.hasNext() {
		el := enumerator.next()
		if comparer(t, el) && comparer(el, t) {
			return true
		}
	}
	return false
}

func (e *enumerableImpl[T]) Entries() []T {
	a := make([]T, 0)
	e.Each(func(el T) {
		a = append(a, el)
	})
	return a
}

func isNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func notIsNil[T any](t T) bool {
	return !isNil(reflect.ValueOf(t))
}

func (e *enumerableImpl[T]) All(f ...Predicate[T]) bool {
	if len(f) > 1 {
		panic("Invalid usage of All")
	}
	if len(f) == 0 {
		return e.All(notIsNil)
	} else {
		enumerator := e.EnumeratorGenerator.create()
		for enumerator.hasNext() {
			if !f[0](enumerator.next()) {
				return false
			}
		}
		return true
	}
}

func (e *enumerableImpl[T]) Any(f ...Predicate[T]) bool {
	if len(f) > 1 {
		panic("Invalid usage of All")
	}
	if len(f) == 0 {
		return e.Any(notIsNil)
	} else {
		enumerator := e.EnumeratorGenerator.create()
		for enumerator.hasNext() {
			if !f[0](enumerator.next()) {
				return true
			}
		}
		return false
	}
}

func (e *enumerableImpl[T]) None(f ...Predicate[T]) bool {
	if len(f) > 1 {
		panic("Invalid usage of None")
	}
	if len(f) == 0 {
		return !e.Any()
	} else {
		return !e.Any(f[0])
	}
}

func (e *enumerableImpl[T]) One(f ...Predicate[T]) bool {
	if len(f) > 1 {
		panic("Invalid usage of One")
	}
	if len(f) == 0 {
		return e.One(notIsNil)
	} else {
		enumerator := e.EnumeratorGenerator.create()
		found := false
		for enumerator.hasNext() {
			if f[0](enumerator.next()) {
				if found {
					return false
				} else {
					found = true
				}
			}
		}
		return found
	}
}

func (e *enumerableImpl[T]) Count(f ...Predicate[T]) int {
	if len(f) > 1 {
		panic("Invalid usage of Count")
	}
	if len(f) == 0 {
		return e.Count(notIsNil)
	} else {
		res := 0
		enumerator := e.EnumeratorGenerator.create()
		for enumerator.hasNext() {
			if f[0](enumerator.next()) {
				res++
			}
		}
		return res
	}
}

func (e *comparableEnumerableImpl[T]) Tally(hash ...map[T]int) map[T]int {
	if len(hash) > 1 {
		panic("invalid usage of Tally")
	}
	var res map[T]int
	if len(hash) == 1 {
		res = hash[0]
	} else {
		res = make(map[T]int)
	}
	e.Each(func(el T) {
		count, exists := res[el]
		if exists {
			res[el] = count + 1
		} else {
			res[el] = 1
		}
	})
	return res
}
