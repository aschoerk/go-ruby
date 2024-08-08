package ruby

func (e *enumerableImpl[T]) Each(f func(T)) {
	enum := e.enumerator.Clone()
	for enum.HasNext() {
		f(*firstParam(enum.Next))
	}
}

func (e *enumerableImpl[T]) EachWithIndex(f func(int, T)) {
	i := 0
	enum := e.enumerator.Clone()
	for enum.HasNext() {
		f(i, *firstParam(enum.Next))
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
	enum := e.enumerator.Clone()
	for enum.HasNext() {
		el := *firstParam(enum.Next)
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

func (e *enumerableImpl[T]) All(f Predicate[T]) bool {

	enum := e.enumerator.Clone()
	for enum.HasNext() {
		if !f(*firstParam(enum.Next)) {
			return false
		}
	}
	return true

}

func (e *enumerableImpl[T]) Any(f Predicate[T]) bool {

	enum := e.enumerator.Clone()
	for enum.HasNext() {
		if !f(*firstParam(enum.Next)) {
			return true
		}
	}
	return false

}

func (e *enumerableImpl[T]) None(f Predicate[T]) bool {
	return !e.Any(f)
}

func (e *enumerableImpl[T]) One(f Predicate[T]) bool {

	found := false
	enum := e.enumerator.Clone()
	for enum.HasNext() {
		if f(*firstParam(enum.Next)) {
			if found {
				return false
			} else {
				found = true
			}
		}
	}
	return found

}

func (e *enumerableImpl[T]) Count(f Predicate[T]) int {
	res := 0
	enum := e.enumerator.Clone()
	for enum.HasNext() {
		if f(*firstParam(enum.Next)) {
			res++
		}
	}
	return res
}

func (e *enumerableImpl[T]) CountAll() int {
	res := 0
	enum := e.enumerator.Clone()
	for existsParam(enum.Next) {
		res++
	}
	return res
}

func (e *enumerableImpl[T]) TallyTo(hash Hash[T, int]) Hash[T, int] {
	e.Each(func(el T) {
		count, exists := hash.Get(el)
		if exists {
			hash.Set(el, count+1)
		} else {
			hash.Set(el, 1)
		}
	})
	return hash
}

func (e *enumerableImpl[T]) Tally() Hash[T, int] {
	res := NewHash[T, int]()
	return e.TallyTo(res)
}
