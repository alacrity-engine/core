package collections

import (
	avltree "github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

type AVLDictionary[TKey constraints.Ordered, TValue any] struct {
	tree *avltree.AVLTree[TKey, TValue]
}

func (avl *AVLDictionary[TKey, TValue]) Erase() error {
	return avl.tree.Erase()
}

func (avl *AVLDictionary[TKey, TValue]) SetPool(pool *mempool.Pool[*avltree.AVLNode[TKey, TValue]]) {
	avl.tree.SetPool(pool)
}

func (avl *AVLDictionary[TKey, TValue]) Add(key TKey, value TValue) error {
	avl.tree.Add(key, value)
	return nil
}

func (avl *AVLDictionary[TKey, TValue]) AddOrUpdate(key TKey, value TValue, upd func(oldValue TValue) (TValue, error)) error {
	return avl.tree.AddOrUpdate(key, value, upd)
}

func (avl *AVLDictionary[TKey, TValue]) Update(key TKey, upd func(value TValue, found bool) (TValue, error)) error {
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

func (avl *AVLDictionary[TKey, TValue]) Remove(key TKey) error {
	avl.tree.Remove(key)
	return nil
}

func (avl *AVLDictionary[TKey, TValue]) Search(key TKey) (TValue, bool, error) {
	var zeroVal TValue
	node := avl.tree.Search(key)

	if node == nil {
		return zeroVal, false, nil
	}

	return node.Value, true, nil
}

func (avl *AVLDictionary[TKey, TValue]) VisitInOrder(visit func(key TKey, value TValue)) {
	avl.tree.VisitInOrder(func(node *avltree.AVLNode[TKey, TValue]) {
		visit(node.Key(), node.Value)
	})
}

func NewAVLDictionary[TKey constraints.Ordered, TValue any](options ...AVLDictionaryOption[TKey, TValue]) (*AVLDictionary[TKey, TValue], error) {
	avl := &AVLDictionary[TKey, TValue]{}
	params := avlDictionaryParams[TKey, TValue]{}

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
