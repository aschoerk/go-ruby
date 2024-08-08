package ruby

import "golang.org/x/exp/constraints"

func RStepped[T constraints.Integer](start, end, step T) Enumerable[T] {
	return &enumerableImpl[T]{enumerator: &rangeEnumerator[T]{start, end, step, start}}
}

func R[T constraints.Integer](start, end T) Enumerable[T] {
	return RStepped(start, end, 1)
}

type rangeEnumerator[T constraints.Integer] struct {
	start, end, step T
	pos              T
}

func (e *rangeEnumerator[T]) HasNext() bool {
	return e.pos < e.end
}

func (e *rangeEnumerator[T]) Next() (*T, bool) {
	if e.HasNext() {
		res := e.pos
		e.pos += e.step
		return &res, true
	} else {
		return nil, false
	}
}

func (e *rangeEnumerator[T]) Clone() Enumerator[T] {
	return &rangeEnumerator[T]{e.start, e.end, e.step, e.start}
}
