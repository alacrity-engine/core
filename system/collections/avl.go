package collections

import (
	avltree "github.com/zergon321/go-avltree"
	"golang.org/x/exp/constraints"
)

type AVLTree[TKey constraints.Ordered, TValue any] struct {
	tree *avltree.AVLTree[TKey, TValue]
}

func (avl *AVLTree[TKey, TValue]) Add(key TKey, value TValue) error {
	avl.tree.Add(key, value)
	return nil
}

func (avl *AVLTree[TKey, TValue]) Update(key TKey, upd func(value TValue, found bool) TValue) error {
	node := avl.tree.Search(key)
	newValue := upd(node.Value, node != nil)
	node.Value = newValue

	return nil
}

func (avl *AVLTree[TKey, TValue]) Remove(key TKey) error {
	avl.tree.Remove(key)
	return nil
}

func (avl *AVLTree[TKey, TValue]) Search(key TKey) (TValue, bool, error) {
	var zeroVal TValue
	node := avl.tree.Search(key)

	if node == nil {
		return zeroVal, false, nil
	}

	return node.Value, true, nil
}

func (avl *AVLTree[TKey, TValue]) VisitInOrder(visit func(key TKey, value TValue)) error {
	avl.tree.VisitInOrder(func(node *avltree.AVLNode[TKey, TValue]) {
		visit(node.Key(), node.Value)
	})

	return nil
}

func NewAVLTree[TKey constraints.Ordered, TValue any](options ...AVLTreeOption[TKey, TValue]) (*AVLTree[TKey, TValue], error) {
	avl := &AVLTree[TKey, TValue]{}
	params := avlTreeParams[TKey, TValue]{}

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
