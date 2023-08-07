package collections

import (
	"github.com/zergon321/mempool"
)

type AVLUnrestrictedSortedDictionaryOption[
	TKey Comparable, TValue any,
] func(
	tree *AVLUnrestrictedSortedDictionary[TKey, TValue],
	params *avlUnrestrictedSortedDictionaryParams[TKey, TValue],
) error

type avlUnrestrictedSortedDictionaryParams[TKey Comparable, TValue any] struct {
	innerTreeOptions []UnrestrictedAVLTreeOption[TKey, TValue]
}

func AVLUnrestrictedSortedDictionaryOptionWithMemoryPool[
	TKey Comparable, TValue any,
](
	pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]],
) AVLUnrestrictedSortedDictionaryOption[TKey, TValue] {
	return func(tree *AVLUnrestrictedSortedDictionary[TKey, TValue], params *avlUnrestrictedSortedDictionaryParams[TKey, TValue]) error {
		params.innerTreeOptions = append(params.innerTreeOptions,
			UnrestrictedAVLTreeOptionWithMemoryPool(pool))

		return nil
	}
}
