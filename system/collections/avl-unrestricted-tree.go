package collections

import (
	"fmt"

	"github.com/zergon321/mempool"
	"golang.org/x/exp/constraints"
)

// AVLTree[TKey constraints.Ordered, TValue any] structure. Public methods are Add, Remove, Update, Search, DisplayTreeInOrder.
type UnrestrictedAVLTree[TKey Comparable, TValue any] struct {
	root *UnrestrictedAVLNode[TKey, TValue]
	pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]]
}

func (t *UnrestrictedAVLTree[TKey, TValue]) Erase() error {
	t.root = nil
	t.pool = nil

	return nil
}

func (t *UnrestrictedAVLTree[TKey, TValue]) SetPool(pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]]) {
	t.pool = pool
}

func (t *UnrestrictedAVLTree[TKey, TValue]) Add(key TKey, value TValue) {
	t.root = t.root.add(key, value, t.pool)
}

func (t *UnrestrictedAVLTree[TKey, TValue]) AddOrUpdate(
	key TKey, value TValue,
	upd func(oldValue TValue) (TValue, error),
) error {
	root, err := t.root.addOrUpdate(key, value, upd, t.pool)

	if err != nil {
		return err
	}

	t.root = root

	return nil
}

func (t *UnrestrictedAVLTree[TKey, TValue]) Remove(key TKey) {
	t.root = t.root.remove(key, t.pool)
}

func (t *UnrestrictedAVLTree[TKey, TValue]) Update(oldKey TKey, newKey TKey, newValue TValue) {
	t.root = t.root.remove(oldKey, t.pool)
	t.root = t.root.add(newKey, newValue, t.pool)
}

func (t *UnrestrictedAVLTree[TKey, TValue]) Search(key TKey) (node *UnrestrictedAVLNode[TKey, TValue]) {
	return t.root.search(key)
}

func (t *UnrestrictedAVLTree[TKey, TValue]) VisitInOrder(visit func(node *UnrestrictedAVLNode[TKey, TValue]) error) error {
	return t.visitInOrder(t.root, visit)
}

func (t *UnrestrictedAVLTree[TKey, TValue]) visitInOrder(node *UnrestrictedAVLNode[TKey, TValue], visit func(node *UnrestrictedAVLNode[TKey, TValue]) error) error {
	if node == nil {
		return nil
	}

	if node.left != nil {
		err := t.visitInOrder(node.left, visit)

		if err != nil {
			return err
		}
	}

	if node != nil {
		err := visit(node)

		if err != nil {
			return err
		}
	}

	if node.right != nil {
		err := t.visitInOrder(node.right, visit)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *UnrestrictedAVLTree[TKey, TValue]) DisplayInOrder() {
	t.root.displayNodesInOrder()
}

// AVLNode structure
type UnrestrictedAVLNode[TKey Comparable, TValue any] struct {
	key   TKey
	Value TValue

	// height counts nodes (not edges)
	height int
	left   *UnrestrictedAVLNode[TKey, TValue]
	right  *UnrestrictedAVLNode[TKey, TValue]
}

// Key returns the key of the AVL tree node.
func (node *UnrestrictedAVLNode[TKey, TValue]) Key() TKey {
	return node.key
}

// Erase nullifies all the
// fields of the AVL tree node.
func (node *UnrestrictedAVLNode[TKey, TValue]) Erase() error {
	var (
		zeroValTKey   TKey
		zeroValTValue TValue
	)

	node.key = zeroValTKey
	node.Value = zeroValTValue
	node.height = 0
	node.left = nil
	node.right = nil

	return nil
}

// Adds a new node
func (n *UnrestrictedAVLNode[TKey, TValue]) add(key TKey, value TValue, pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]]) *UnrestrictedAVLNode[TKey, TValue] {
	if n == nil {
		if pool != nil {
			node := pool.Get()

			node.key = key
			node.Value = value
			node.height = 1

			return node
		}

		return &UnrestrictedAVLNode[TKey, TValue]{key, value, 1, nil, nil}
	}

	if key.Less(n.key) {
		n.left = n.left.add(key, value, pool)
	} else if key.Greater(n.key) {
		n.right = n.right.add(key, value, pool)
	} else {
		// if same key exists update value
		n.Value = value
	}
	return n.rebalanceTree()
}

func (n *UnrestrictedAVLNode[TKey, TValue]) addOrUpdate(
	key TKey, value TValue,
	upd func(oldValue TValue) (TValue, error),
	pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]],
) (*UnrestrictedAVLNode[TKey, TValue], error) {
	var err error

	if n == nil {
		if pool != nil {
			node := pool.Get()

			node.key = key
			node.Value = value
			node.height = 1

			return node, nil
		}

		return &UnrestrictedAVLNode[TKey, TValue]{key, value, 1, nil, nil}, nil
	}

	if key.Less(n.key) {
		n.left, err = n.left.addOrUpdate(key, value, upd, pool)

		if err != nil {
			return nil, err
		}
	} else if key.Greater(n.key) {
		n.right, err = n.right.addOrUpdate(key, value, upd, pool)

		if err != nil {
			return nil, err
		}
	} else {
		// if same key exists update value
		value, err := upd(n.Value)

		if err != nil {
			return nil, err
		}

		n.Value = value
	}

	return n.rebalanceTree(), nil
}

// Removes a node
func (n *UnrestrictedAVLNode[TKey, TValue]) remove(key TKey, pool *mempool.Pool[*UnrestrictedAVLNode[TKey, TValue]]) *UnrestrictedAVLNode[TKey, TValue] {
	if n == nil {
		return nil
	}
	if key.Less(n.key) {
		n.left = n.left.remove(key, pool)
	} else if key.Greater(n.key) {
		n.right = n.right.remove(key, pool)
	} else {
		if n.left != nil && n.right != nil {
			// node to delete found with both children;
			// replace values with smallest node of the right sub-tree
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.Value = rightMinNode.Value
			// delete smallest node that we replaced
			n.right = n.right.remove(rightMinNode.key, pool)
		} else if n.left != nil {
			// node only has left child
			node := n
			n = n.left

			if pool != nil {
				pool.Put(node)
			}
		} else if n.right != nil {
			// node only has right child
			node := n
			n = n.right

			if pool != nil {
				pool.Put(node)
			}
		} else {
			// node has no children
			node := n
			n = nil

			if pool != nil {
				pool.Put(node)
			}

			return n
		}

	}
	return n.rebalanceTree()
}

// Searches for a node
func (n *UnrestrictedAVLNode[TKey, TValue]) search(key TKey) *UnrestrictedAVLNode[TKey, TValue] {
	if n == nil {
		return nil
	}
	if key.Less(n.key) {
		return n.left.search(key)
	} else if key.Greater(n.key) {
		return n.right.search(key)
	} else {
		return n
	}
}

// Displays nodes left-depth first (used for debugging)
func (n *UnrestrictedAVLNode[TKey, TValue]) displayNodesInOrder() {
	if n.left != nil {
		n.left.displayNodesInOrder()
	}
	fmt.Print(n.key, " ")
	if n.right != nil {
		n.right.displayNodesInOrder()
	}
}

func (n *UnrestrictedAVLNode[TKey, TValue]) getHeight() int {
	if n == nil {
		return 0
	}
	return n.height
}

func (n *UnrestrictedAVLNode[TKey, TValue]) recalculateHeight() {
	n.height = 1 + maxElem(n.left.getHeight(), n.right.getHeight())
}

// Returns maxElem number - TODO: std lib seemed to only have a method for floats!
func maxElem[TKey constraints.Ordered](a TKey, b TKey) TKey {
	if a > b {
		return a
	}
	return b
}

// Checks if node is balanced and rebalance
func (n *UnrestrictedAVLNode[TKey, TValue]) rebalanceTree() *UnrestrictedAVLNode[TKey, TValue] {
	if n == nil {
		return n
	}
	n.recalculateHeight()

	// check balance factor and rotateLeft if right-heavy and rotateRight if left-heavy
	balanceFactor := n.left.getHeight() - n.right.getHeight()
	if balanceFactor == -2 {
		// check if child is left-heavy and rotateRight first
		if n.right.left.getHeight() > n.right.right.getHeight() {
			n.right = n.right.rotateRight()
		}
		return n.rotateLeft()
	} else if balanceFactor == 2 {
		// check if child is right-heavy and rotateLeft first
		if n.left.right.getHeight() > n.left.left.getHeight() {
			n.left = n.left.rotateLeft()
		}
		return n.rotateRight()
	}
	return n
}

// Rotate nodes left to balance node
func (n *UnrestrictedAVLNode[TKey, TValue]) rotateLeft() *UnrestrictedAVLNode[TKey, TValue] {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

// Rotate nodes right to balance node
func (n *UnrestrictedAVLNode[TKey, TValue]) rotateRight() *UnrestrictedAVLNode[TKey, TValue] {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

// Finds the smallest child (based on the key) for the current node
func (n *UnrestrictedAVLNode[TKey, TValue]) findSmallest() *UnrestrictedAVLNode[TKey, TValue] {
	if n.left != nil {
		return n.left.findSmallest()
	} else {
		return n
	}
}

// NewAVLTree creates a new
// AVL tree with the specified options.
func NewUnrestrictedAVLTree[
	TKey Comparable, TValue any,
](
	options ...UnrestrictedAVLTreeOption[TKey, TValue],
) (
	*UnrestrictedAVLTree[TKey, TValue], error,
) {
	tree := &UnrestrictedAVLTree[TKey, TValue]{}

	for i := 0; i < len(options); i++ {
		option := options[i]
		err := option(tree)

		if err != nil {
			return nil, err
		}
	}

	return tree, nil
}
