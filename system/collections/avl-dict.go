package collections

import (
	avltree "github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLSortedDictionary[TKey constraints.Ordered, TValue any] struct {
	tree *avltree.AVLTree[TKey, TValue]
}

func (avl *AVLSortedDictionary[TKey, TValue]) Erase() error {
	return avl.tree.Erase()
}

func (avl *AVLSortedDictionary[TKey, TValue]) SetPool(pool *mempool.Pool[*avltree.AVLNode[TKey, TValue]]) {
	avl.tree.SetPool(pool)
}

func (avl *AVLSortedDictionary[TKey, TValue]) Add(key TKey, value TValue) error {
	avl.tree.Add(key, value)
	return nil
}

func (avl *AVLSortedDictionary[TKey, TValue]) AddOrUpdate(key TKey, value TValue, upd func(oldValue TValue) (TValue, error)) error {
	return avl.tree.AddOrUpdate(key, value, upd)
}

func (avl *AVLSortedDictionary[TKey, TValue]) Update(key TKey, upd func(value TValue, found bool) (TValue, error)) error {
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

func (avl *AVLSortedDictionary[TKey, TValue]) Remove(key TKey) error {
	avl.tree.Remove(key)
	return nil
}

func (avl *AVLSortedDictionary[TKey, TValue]) Search(key TKey) (TValue, bool, error) {
	var zeroVal TValue
	node := avl.tree.Search(key)

	if node == nil {
		return zeroVal, false, nil
	}

	return node.Value, true, nil
}

func (avl *AVLSortedDictionary[TKey, TValue]) VisitInOrder(visit func(key TKey, value TValue) error) error {
	return avl.tree.VisitInOrder(func(node *avltree.AVLNode[TKey, TValue]) error {
		return visit(node.Key(), node.Value)
	})
}

func NewAVLSortedDictionary[TKey constraints.Ordered, TValue any](options ...AVLSortedDictionaryOption[TKey, TValue]) (*AVLSortedDictionary[TKey, TValue], error) {
	avl := &AVLSortedDictionary[TKey, TValue]{}
	params := avlSortedDictionaryParams[TKey, TValue]{}

	for i := 0; i < len(options); i++ {
		option := options[i]
		err := option(avl, &params)

		if err != nil {
			return nil, err
		}
	}

	innerTree, err := avltree.NewAVLTree(
		params.innerTreeOptions...)

	if err != nil {
		return nil, err
	}

	avl.tree = innerTree

	return avl, nil
}
