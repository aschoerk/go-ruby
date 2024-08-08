package ruby

type enumerableImpl[T any] struct {
	enumerator Enumerator[T]
}

func e[T any](enumerator Enumerator[T]) Enumerable[T] {
	return &enumerableImpl[T]{enumerator: enumerator}
}

func firstParam[T any](f func() (T, bool)) T {
	res, exists := f()
	if !exists {
		panic("expected to exist")
	}
	return res
}

func existsParam[T any](f func() (T, bool)) bool {
	_, exists := f()
	return exists
}

type commonEnumerator[T any] struct {
	enum       Enumerator[T]
	found      *T
	exists     bool
	didHasNext bool
}

func initCommonEnumerator[T any](e Enumerator[T]) commonEnumerator[T] {
	return commonEnumerator[T]{e, nil, false, false}
}

func (e *commonEnumerator[T]) clone() commonEnumerator[T] {
	return commonEnumerator[T]{e.enum.Clone(), nil, false, false}
}
