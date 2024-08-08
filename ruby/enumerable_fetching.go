package ruby

func (e *enumerableImpl[T]) First() (T, bool) {
	res, exists := e.enumerator.Clone().Next()
	if !exists {
		var zeroValue T
		return zeroValue, false
	} else {
		return *res, true
	}
}

func (e *enumerableImpl[T]) FirstN(max int) Enumerable[T] {
	return &enumerableImpl[T]{&firstTakeEnumerator[T]{e.enumerator.Clone(), 0, max}}
}

func (e *enumerableImpl[T]) TakeN(max int) Enumerable[T] {
	return &enumerableImpl[T]{&firstTakeEnumerator[T]{e.enumerator.Clone(), 0, max}}
}

func (e *enumerableImpl[T]) Drop(max int) Enumerable[T] {
	res := e.enumerator.Clone()
	for i := 0; i < max; i++ {
		res.Next()
	}
	return &enumerableImpl[T]{res}
}

func (e *enumerableImpl[T]) TakeWhile(p Predicate[T]) Enumerable[T] {
	return &enumerableImpl[T]{&takeWhileEnum[T]{initCommonEnumerator(e.enumerator.Clone()), p}}
}

func (e *enumerableImpl[T]) DropWhile(p Predicate[T]) Enumerable[T] {
	return &enumerableImpl[T]{&dropWhileEnum[T]{initCommonEnumerator(e.enumerator.Clone()), p, false}}
}

type takeWhileEnum[T any] struct {
	commonEnumerator[T]
	p Predicate[T]
}

func (e *takeWhileEnum[T]) HasNext() bool {
	if e.didHasNext {
		return e.exists
	} else {
		e.didHasNext = true
		e.exists = false
		e.found = nil
		if e.enum.HasNext() {
			tmp, _ := e.enum.Next()
			if e.p(*tmp) {
				e.exists = true
				e.found = tmp
			} else {
				e.exists = false
				e.found = nil
			}
		}
		return e.exists
	}
}

func (e *takeWhileEnum[T]) Next() (*T, bool) {
	if !e.didHasNext {
		e.HasNext()
	}
	if e.exists {
		// first false ends iterating, hasNext is not necessary anymore
		e.didHasNext = false
	}
	return e.found, e.exists
}

func (e *takeWhileEnum[T]) Clone() Enumerator[T] {
	return &takeWhileEnum[T]{e.clone(), e.p}
}

type dropWhileEnum[T any] struct {
	commonEnumerator[T]
	p          Predicate[T]
	droppedAll bool
}

func (e *dropWhileEnum[T]) HasNext() bool {
	if e.didHasNext {
		return e.exists
	}
	e.didHasNext = true
	if e.droppedAll {
		hasNext := e.enum.HasNext()
		if hasNext {
			e.found, e.exists = e.enum.Next()
		} else {
			e.found = nil
			e.exists = false
		}
	} else {
		for e.enum.HasNext() {
			tmp, _ := e.enum.Next()
			if !e.p(*tmp) {
				e.exists = true
				e.found = tmp
				e.droppedAll = true
				break
			}
		}
		if !e.droppedAll { // all had to be dropped, stream ends here
			e.droppedAll = true
			e.exists = false
			e.found = nil
		}
	}
	return e.exists
}

func (e *dropWhileEnum[T]) Next() (*T, bool) {
	if !e.didHasNext {
		e.HasNext()
	}
	e.didHasNext = false
	return e.found, e.exists
}

func (e *dropWhileEnum[T]) Clone() Enumerator[T] {
	return &dropWhileEnum[T]{e.clone(), e.p, false}
}

type firstTakeEnumerator[T any] struct {
	enumerator Enumerator[T]
	delivered  int
	max        int
}

func (e *firstTakeEnumerator[T]) HasNext() bool {
	if e.delivered < e.max {
		return e.enumerator.HasNext()
	} else {
		return false
	}
}

func (e *firstTakeEnumerator[T]) Next() (*T, bool) {
	if e.delivered < e.max {
		e.delivered++
		res, exists := e.enumerator.Next()
		return res, exists
	} else {
		return nil, false
	}
}

func (e *firstTakeEnumerator[T]) Clone() Enumerator[T] {
	return &firstTakeEnumerator[T]{e.enumerator.Clone(), 0, e.max}
}
