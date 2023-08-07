package collections

import (
	"github.com/zergon321/mempool"
)

type AVLUnrestrictedSortedDictionary[TKey Comparable, TValue any] struct {
	tree *UnrestrictedAVLTree[TKey, TValue]
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) Erase() error {
	return avl.tree.Erase()
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) SetPool(pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]]) {
	avl.tree.SetPool(pool)
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) Add(key TKey, value TValue) error {
	avl.tree.Add(key, value)
	return nil
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) AddOrUpdate(key TKey, value TValue, upd func(oldValue TValue) (TValue, error)) error {
	return avl.tree.AddOrUpdate(key, value, upd)
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) Update(key TKey, upd func(value TValue, found bool) (TValue, error)) error {
	node := avl.tree.Search(key)

	var value TValue

	if node != nil {
		value = node.Value
	}

	newValue, err := upd(value, node != nil)

	if err != nil {
		return err
	}

	if node != nil {
		node.Value = newValue
	}

	return nil
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) Remove(key TKey) error {
	avl.tree.Remove(key)
	return nil
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) Search(key TKey) (TValue, bool, error) {
	var zeroVal TValue
	node := avl.tree.Search(key)

	if node == nil {
		return zeroVal, false, nil
	}

	return node.Value, true, nil
}

func (avl *AVLUnrestrictedSortedDictionary[TKey, TValue]) VisitInOrder(visit func(key TKey, value TValue)) {
	avl.tree.VisitInOrder(func(node *UnrestrictedAVLNode[TKey, TValue]) {
		visit(node.Key(), node.Value)
	})
}

func NewAVLUnrestrictedSortedDictionary[TKey Comparable, TValue any](options ...AVLUnrestrictedSortedDictionaryOption[TKey, TValue]) (*AVLUnrestrictedSortedDictionary[TKey, TValue], error) {
	avl := &AVLUnrestrictedSortedDictionary[TKey, TValue]{}
	params := avlUnrestrictedSortedDictionaryParams[TKey, TValue]{}

	for i := 0; i < len(options); i++ {
		option := options[i]
		err := option(avl, &params)

		if err != nil {
			return nil, err
		}
	}

	innerTree, err := NewUnrestrictedAVLTree(
		params.innerTreeOptions...)

	if err != nil {
		return nil, err
	}

	avl.tree = innerTree

	return avl, nil
}
