package collections

import (
	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLDictionaryOption[
	TKey constraints.Ordered, TValue any,
] func(
	tree *AVLDictionary[TKey, TValue],
	params *avlDictionaryParams[TKey, TValue],
) error

type avlDictionaryParams[TKey constraints.Ordered, TValue any] struct {
	innerTreeOptions []avltree.AVLTreeOption[TKey, TValue]
}

func AVLDictionaryOptionWithMemoryPool[
	TKey constraints.Ordered, TValue any,
](
	pool *mempool.Pool[*avltree.AVLNode[TKey, TValue]],
) AVLDictionaryOption[TKey, TValue] {
	return func(tree *AVLDictionary[TKey, TValue], params *avlDictionaryParams[TKey, TValue]) error {
		params.innerTreeOptions = append(params.innerTreeOptions,
			avltree.AVLTreeOptionWithMemoryPool(pool))

		return nil
	}
}
