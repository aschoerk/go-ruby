package ruby

type Enumerator[T any] interface {
	HasNext() bool
	Next() (*T, bool)
	Clone() Enumerator[T]
}

type Predicate[T any] func(T) bool

type Comparer[T any] func(T, T) int

type Hash[K any, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Delete(key K)
	Len() int
}

func GenEquals[T comparable](u T) Predicate[T] {
	return func(v T) bool {
		return u == v
	}
}

func GenTrue[T comparable]() Predicate[T] {
	return func(v T) bool {
		return true
	}
}

func GenFalse[T comparable]() Predicate[T] {
	return func(v T) bool {
		return false
	}
}

type Enumerable[T any] interface {
	// Querying
	Includes(T, ...func(T, T) bool) bool
	// All() bool
	All(Predicate[T]) bool
	Any(Predicate[T]) bool
	None(Predicate[T]) bool
	One(Predicate[T]) bool
	Count(Predicate[T]) int
	CountAll() int
	TallyTo(Hash[T, int]) Hash[T, int]
	Tally() Hash[T, int]

	// Fetching
	First() (T, bool)
	FirstN(int) Enumerable[T]
	Drop(int) Enumerable[T]
	TakeWhile(Predicate[T]) Enumerable[T]
	DropWhile(Predicate[T]) Enumerable[T]

	// Iterating
	Each(func(T))
	EachWithIndex(func(int, T))
	Entries() []T

	Min(Comparer[T]) (T, bool)

	Filter(Predicate[T]) Enumerable[T]
}
