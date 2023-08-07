package collections

import "github.com/zergon321/mempool"

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
