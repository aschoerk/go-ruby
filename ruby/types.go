package ruby

type Enumerator[T any] interface {
	hasNext() bool
	next() T
}

type EnumeratorGenerator[T any] interface {
	create() Enumerator[T]
}

type Predicate[T any] func(T) bool

type Hash[K any, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Delete(key K)
}

type Enumerable[T any] interface {
	// Querying
	Includes(T, ...func(T, T) bool) bool
	// All() bool
	All(...Predicate[T]) bool
	Any(...Predicate[T]) bool
	None(...Predicate[T]) bool
	One(...Predicate[T]) bool
	Count(...Predicate[T]) int
	Tally(...Hash[T, int]) Hash[T, int]

	// Iterating
	Each(func(T))
	EachWithIndex(func(int, T))
	Entries() []T

	Filter(...Predicate[T]) Enumerable[T]
}
