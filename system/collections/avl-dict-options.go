package collections

import (
	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLSortedDictionaryOption[
	TKey constraints.Ordered, TValue any,
] func(
	tree *AVLSortedDictionary[TKey, TValue],
	params *avlSortedDictionaryParams[TKey, TValue],
) error

type avlSortedDictionaryParams[TKey constraints.Ordered, TValue any] struct {
	innerTreeOptions []avltree.AVLTreeOption[TKey, TValue]
}

func AVLSortedDictionaryOptionWithMemoryPool[
	TKey constraints.Ordered, TValue any,
](
	pool *mempool.Pool[*avltree.AVLNode[TKey, TValue]],
) AVLSortedDictionaryOption[TKey, TValue] {
	return func(tree *AVLSortedDictionary[TKey, TValue], params *avlSortedDictionaryParams[TKey, TValue]) error {
		params.innerTreeOptions = append(params.innerTreeOptions,
			avltree.AVLTreeOptionWithMemoryPool(pool))

		return nil
	}
}
