package collections

import (
	"errors"

	"github.com/zergon321/mempool"
)

type AVLunrestrictedSortedDictionaryPooledProducer[TKey Comparable, TValue any] struct {
	pool     *mempool.Pool[*AVLUnrestrictedSortedDictionary[TKey, TValue]]
	nodePool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]]
}

func (prod *AVLunrestrictedSortedDictionaryPooledProducer[TKey, TValue]) Produce() (UnrestrictedSortedDictionary[TKey, TValue], error) {
	tree := prod.pool.Get()
	tree.SetPool(prod.nodePool)

	return tree, nil
}

func (prod *AVLunrestrictedSortedDictionaryPooledProducer[TKey, TValue]) Dispose(dict UnrestrictedSortedDictionary[TKey, TValue]) error {
	tree, ok := dict.(*AVLUnrestrictedSortedDictionary[TKey, TValue])

	if !ok {
		return errors.New("incorrect type")
	}

	return prod.pool.Put(tree)
}

func NewAVLUnrestrictedSortedDictionaryProducer[TKey Comparable, TValue any](
	pool *mempool.Pool[*AVLUnrestrictedSortedDictionary[TKey, TValue]],
	nodePool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]],
) *AVLunrestrictedSortedDictionaryPooledProducer[TKey, TValue] {
	return &AVLunrestrictedSortedDictionaryPooledProducer[TKey, TValue]{
		pool:     pool,
		nodePool: nodePool,
	}
}
