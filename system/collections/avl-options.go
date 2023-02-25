package collections

import (
	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLTreeOption[
	TKey constraints.Ordered, TValue any,
] func(
	tree *AVLTree[TKey, TValue],
	params *avlTreeParams[TKey, TValue],
) error

type avlTreeParams[TKey constraints.Ordered, TValue any] struct {
	innerTreeOptions []avltree.AVLTreeOption[TKey, TValue]
}

func AVLTreeOptionWithMemoryPool[
	TKey constraints.Ordered, TValue any,
](
	pool *mempool.Pool[*avltree.AVLNode[TKey, TValue]],
) AVLTreeOption[TKey, TValue] {
	return func(tree *AVLTree[TKey, TValue], params *avlTreeParams[TKey, TValue]) error {
		params.innerTreeOptions = append(params.innerTreeOptions,
			avltree.AVLTreeOptionWithMemoryPool(pool))

		return nil
	}
}
