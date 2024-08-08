package ruby

func E[T any](slice []T) Enumerable[T] {
	return &enumerableImpl[T]{enumerator: &sliceEnumerator[T]{&slice, 0}}
}

type sliceEnumerator[T any] struct {
	data *[]T
	pos  int
}

func (g *sliceEnumerator[T]) HasNext() bool {
	return g.pos < len(*g.data)
}

func (g *sliceEnumerator[T]) Next() (*T, bool) {
	if g.HasNext() {
		res := (*g.data)[g.pos]
		g.pos++
		return &res, true
	} else {
		return nil, false
	}
}

func (g *sliceEnumerator[T]) Clone() Enumerator[T] {
	return &sliceEnumerator[T]{g.data, 0}
}
