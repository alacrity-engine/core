package collections

import "golang.org/x/exp/constraints"

// TODO: add enumerators.

type Comparable interface {
	Less(Comparable) bool
	Greater(Comparable) bool
	Equal(Comparable) bool
}

type SortedDictionary[TKey constraints.Ordered, TValue any] interface {
	// Add doesn't return an error if the key already exists;
	// it updates the underlying value instead.
	Add(key TKey, value TValue) error
	AddOrUpdate(key TKey, value TValue, upd func(oldValue TValue) (TValue, error)) error
	Update(key TKey, upd func(value TValue, found bool) (TValue, error)) error
	// Remove doesn't return an error if the key doesn't exist.
	Remove(key TKey) error
	Search(key TKey) (TValue, bool, error)
	VisitInOrder(func(key TKey, value TValue))
}

type SortedDictionaryProducer[TKey constraints.Ordered, TValue any] interface {
	Produce() (SortedDictionary[TKey, TValue], error)
	Dispose(dict SortedDictionary[TKey, TValue]) error
}

/******************************************************************************************/

type UnrestrictedSortedDictionary[TKey Comparable, TValue any] interface {
	// Add doesn't return an error if the key already exists;
	// it updates the underlying value instead.
	Add(key TKey, value TValue) error
	AddOrUpdate(key TKey, value TValue, upd func(oldValue TValue) (TValue, error)) error
	Update(key TKey, upd func(value TValue, found bool) (TValue, error)) error
	// Remove doesn't return an error if the key doesn't exist.
	Remove(key TKey) error
	Search(key TKey) (TValue, bool, error)
	VisitInOrder(func(key TKey, value TValue))
}

type UnrestrictedSortedSet[TKey Comparable] interface {
	// Add doesn't return an error if the key already exists.
	Add(key TKey) error
	// Remove doesn't return an error if the key doesn't exist.
	Remove(key TKey) error
	Search(key TKey) (bool, error)
	VisitInOrder(func(key TKey))
}
