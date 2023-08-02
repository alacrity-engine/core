package collections

import (
	"github.com/zergon321/mempool"
)

type AVLUnrestrictedSetOption[
	TKey Comparable,
] func(
	tree *AVLUnrestrictedSet[TKey],
	params *avlUnrestrictedSetParams[TKey],
) error

type avlUnrestrictedSetParams[TKey Comparable] struct {
	innerTreeOptions []UnrestrictedAVLTreeOption[TKey, TKey]
}

func AVLUnrestrictedSetOptionWithMemoryPool[
	TKey Comparable,
](
	pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TKey]],
) AVLUnrestrictedSetOption[TKey] {
	return func(tree *AVLUnrestrictedSet[TKey], params *avlUnrestrictedSetParams[TKey]) error {
		params.innerTreeOptions = append(params.innerTreeOptions,
			UnrestrictedAVLTreeOptionWithMemoryPool[TKey](pool))

		return nil
	}
}

/******************************************************************/

type UnrestrictedAVLTreeOption[
	TKey Comparable, TValue any,
] func(tree *UnrestrictedAVLTree[TKey, TValue]) error

func UnrestrictedAVLTreeOptionWithMemoryPool[
	TKey Comparable, TValue any,
](
	pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]],
) UnrestrictedAVLTreeOption[TKey, TValue] {
	return func(tree *UnrestrictedAVLTree[TKey, TValue]) error {
		tree.pool = pool
		return nil
	}
}
