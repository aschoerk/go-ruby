package ruby

func (e *enumerableImpl[T]) Filter(p Predicate[T]) Enumerable[T] {
	return &enumerableImpl[T]{&filterEnum[T]{initCommonEnumerator(e.enumerator), p}}
}

type filterEnum[T any] struct {
	commonEnumerator[T]
	p Predicate[T]
}

func (e *filterEnum[T]) HasNext() bool {
	if e.didHasNext {
		return e.exists
	} else {
		e.didHasNext = true
		e.exists = false
		e.found = nil
		for e.enum.HasNext() {
			tmp, _ := e.enum.Next()
			if e.p(*tmp) {
				e.exists = true
				e.found = tmp
				break
			}
		}
		return e.exists
	}
}

func (e *filterEnum[T]) Next() (*T, bool) {
	if !e.didHasNext {
		e.HasNext()
	}
	e.didHasNext = false
	return e.found, e.exists
}

func (e *filterEnum[T]) Clone() Enumerator[T] {
	return &filterEnum[T]{e.clone(), e.p}
}
