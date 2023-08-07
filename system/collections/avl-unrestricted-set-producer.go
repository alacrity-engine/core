package collections

import (
	"errors"

	"github.com/zergon321/mempool"
)

type AVLUnrestrictedSetPooledProducer[TKey Comparable] struct {
	pool     *mempool.Pool[*AVLUnrestrictedSortedSet[TKey]]
	nodePool *mempool.Pool[*UnrestrictedAVLNode[TKey, TKey]]
}

func (prod *AVLUnrestrictedSetPooledProducer[TKey]) Produce() (UnrestrictedSortedSet[TKey], error) {
	tree := prod.pool.Get()
	tree.SetPool(prod.nodePool)

	return tree, nil
}

func (prod *AVLUnrestrictedSetPooledProducer[TKey]) Dispose(dict UnrestrictedSortedSet[TKey]) error {
	tree, ok := dict.(*AVLUnrestrictedSortedSet[TKey])

	if !ok {
		return errors.New("incorrect type")
	}

	return prod.pool.Put(tree)
}

func NewAVLUnrestrictedSetPooledProducer[TKey Comparable](
	pool *mempool.Pool[*AVLUnrestrictedSortedSet[TKey]],
	nodePool *mempool.Pool[*UnrestrictedAVLNode[TKey, TKey]],
) *AVLUnrestrictedSetPooledProducer[TKey] {
	return &AVLUnrestrictedSetPooledProducer[TKey]{
		pool:     pool,
		nodePool: nodePool,
	}
}
