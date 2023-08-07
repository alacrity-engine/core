package collections

import (
	"errors"

	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLSortedDictionaryPooledProducer[TKey constraints.Ordered, TValue any] struct {
	pool     *mempool.Pool[*AVLSortedDictionary[TKey, TValue]]
	nodePool *mempool.Pool[*avltree.AVLNode[TKey, TValue]]
}

func (prod *AVLSortedDictionaryPooledProducer[TKey, TValue]) Produce() (SortedDictionary[TKey, TValue], error) {
	tree := prod.pool.Get()
	tree.SetPool(prod.nodePool)

	return tree, nil
}

func (prod *AVLSortedDictionaryPooledProducer[TKey, TValue]) Dispose(dict SortedDictionary[TKey, TValue]) error {
	tree, ok := dict.(*AVLSortedDictionary[TKey, TValue])

	if !ok {
		return errors.New("incorrect type")
	}

	return prod.pool.Put(tree)
}

func NewAVLSortedDictionaryProducer[TKey constraints.Ordered, TValue any](
	pool *mempool.Pool[*AVLSortedDictionary[TKey, TValue]],
	nodePool *mempool.Pool[*avltree.AVLNode[TKey, TValue]],
) *AVLSortedDictionaryPooledProducer[TKey, TValue] {
	return &AVLSortedDictionaryPooledProducer[TKey, TValue]{
		pool:     pool,
		nodePool: nodePool,
	}
}
