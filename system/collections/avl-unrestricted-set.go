package collections

import (
	"github.com/zergon321/mempool"
)

type AVLUnrestrictedSortedSet[TKey Comparable] struct {
	tree *UnrestrictedAVLTree[TKey, TKey]
}

func (avl *AVLUnrestrictedSortedSet[TKey]) SetPool(pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TKey]]) {
	avl.tree.SetPool(pool)
}

func (set *AVLUnrestrictedSortedSet[TKey]) Erase() error {
	return set.tree.Erase()
}

func (set *AVLUnrestrictedSortedSet[TKey]) Add(key TKey) error {
	set.tree.Add(key, key)
	return nil
}

func (set *AVLUnrestrictedSortedSet[TKey]) Remove(key TKey) error {
	set.tree.Remove(key)
	return nil
}

func (set *AVLUnrestrictedSortedSet[TKey]) Search(key TKey) (bool, error) {
	return set.tree.Search(key) != nil, nil
}

func (set *AVLUnrestrictedSortedSet[TKey]) VisitInOrder(visit func(key TKey) error) error {
	return set.tree.VisitInOrder(func(node *UnrestrictedAVLNode[TKey, TKey]) error {
		return visit(node.Key())
	})
}
