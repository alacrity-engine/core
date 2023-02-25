package collections

import "golang.org/x/exp/constraints"

// TODO: add enumerators; AddOrUpdate method.

type SortedDictionary[TKey constraints.Ordered, TValue any] interface {
	Add(key TKey, value TValue) error
	Update(key TKey, upd func(value TValue, found bool) TValue) error
	Remove(key TKey) error
	Search(key TKey) (TValue, bool, error)
	VisitInOrder(func(key TKey, value TValue)) error
}
