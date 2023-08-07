package collections

import (
	"github.com/zergon321/mempool"
)

type AVLUnrestrictedSortedSetOption[
	TKey Comparable,
] func(
	tree *AVLUnrestrictedSortedSet[TKey],
	params *avlUnrestrictedSortedSetParams[TKey],
) error

type avlUnrestrictedSortedSetParams[TKey Comparable] struct {
	innerTreeOptions []UnrestrictedAVLTreeOption[TKey, TKey]
}

func AVLUnrestrictedSortedSetOptionWithMemoryPool[
	TKey Comparable,
](
	pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TKey]],
) AVLUnrestrictedSortedSetOption[TKey] {
	return func(tree *AVLUnrestrictedSortedSet[TKey], params *avlUnrestrictedSortedSetParams[TKey]) error {
		params.innerTreeOptions = append(params.innerTreeOptions,
			UnrestrictedAVLTreeOptionWithMemoryPool[TKey](pool))

		return nil
	}
}
