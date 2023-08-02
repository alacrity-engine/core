package collections

import (
	"errors"

	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLTreePooledProducer[TKey constraints.Ordered, TValue any] struct {
	pool     *mempool.Pool[*AVLDictionary[TKey, TValue]]
	nodePool *mempool.Pool[*avltree.AVLNode[TKey, TValue]]
}

func (prod *AVLTreePooledProducer[TKey, TValue]) Produce() (SortedDictionary[TKey, TValue], error) {
	tree := prod.pool.Get()
	tree.SetPool(prod.nodePool)

	return tree, nil
}

func (prod *AVLTreePooledProducer[TKey, TValue]) Dispose(dict SortedDictionary[TKey, TValue]) error {
	tree, ok := dict.(*AVLDictionary[TKey, TValue])

	if !ok {
		return errors.New("incorrect type")
	}

	return prod.pool.Put(tree)
}

func NewAVLSortedDictionaryProducer[TKey constraints.Ordered, TValue any](
	pool *mempool.Pool[*AVLDictionary[TKey, TValue]],
	nodePool *mempool.Pool[*avltree.AVLNode[TKey, TValue]],
) *AVLTreePooledProducer[TKey, TValue] {
	return &AVLTreePooledProducer[TKey, TValue]{
		pool:     pool,
		nodePool: nodePool,
	}
}
