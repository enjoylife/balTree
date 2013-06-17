package gotree

import (
	"fmt"
	"runtime"
)

// Color is the used to maintain the redblack tree balance.
type color bool

const (
	red   color = false // we rely on default for node initializations
	black color = true
)

// Pretty output for errors, debugging, etc.
func (c color) String() string {
	var s string
	switch c {
	case red:
		s = "red"
	case black:
		s = "black"
	}
	return s
}

// A RBNode is the type manipulated within the tree. It holds the inserted elements.
// It is exposed whenever the tree traversal functions are used.
type RBNode struct {
	Elem        Interface
	left, right *RBNode
	color       color
}

// A RBTree is our main type our redblack tree methods are defined on.
type RBTree struct {
	height      int // height from root to leaf
	size        int // Number of inserted elements
	first, last *RBNode
	iterNext    func() Interface // initially nil
	root        *RBNode
}

// Height returns the max depth of any branch of the tree
func (t *RBTree) Height() int {
	return t.height
}
func (t *RBTree) Size() int {
	return t.size
}
func (t *RBTree) Clear() {
	t.root = nil
	t.last = nil
	t.first = nil
	t.size = 0
	t.height = 0
	t.iterNext = nil
	runtime.GC()
}

// Min returns the smallest inserted element if possible. If the smallest value is not
// found(empty tree), then Min returns a nil.
func (t *RBTree) Min() Interface {
	if t.first != nil {
		return t.first.Elem
	}
	return nil
}

// Max returns the largest inserted element if possible. If the largest value is not
// found(empty tree), then Max returns a nil.
func (t *RBTree) Max() Interface {
	if t.last != nil {
		return t.last.Elem
	}
	return nil
}

// Next is called when individual elements are wanted to be traversed over.
// Prior to a call to Next, a call to IterInit needs to be made to set up the necessary
// data to allow for traversal of the tree. Example:
//
//    sum := 0
//    for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {
//        elem := n.(exInt)  // (exInt is simple int type)
//        sum += int(elem) + i
//    }
// Note: If one was to break out of the loop prior to a complete traversal,
// and start another loop without calling IterInit, then the previously uncompleted iterator is continued again.
func (t *RBTree) Next() Interface {

	if t.iterNext == nil {
		return nil
	}
	return t.iterNext() // func set by call to IterInit(TravOrder)

}

// IterInit is the initializer which setups the tree for iterating over it's elements in
// a specific order. It setups the internal data, and then returns the first RBNode to be looked at. See Next for an example.
func (t *RBTree) IterInit(order TravOrder) Interface {

	current := t.root
	stack := []*RBNode{}
	switch order {
	case InOrder:
		t.iterNext = func() (out Interface) {
			for len(stack) > 0 || current != nil {
				if current != nil {
					stack = append(stack, current)
					current = current.left
				} else {
					// pop
					stackIndex := len(stack) - 1
					current = stack[stackIndex]
					out = current.Elem
					stack = stack[0:stackIndex]
					current = current.right
					break
				}
			}
			// last node, reset
			if out == nil {
				t.iterNext = nil
			}
			return out
		}

	case PreOrder:
		t.iterNext = func() (out Interface) {
			for len(stack) > 0 || current != nil {
				if current != nil {
					out = current.Elem
					stack = append(stack, current.right)
					current = current.left
					break
				} else {
					// pop
					stackIndex := len(stack) - 1
					current = stack[stackIndex]
					stack = stack[0:stackIndex]
				}
			}

			// last node, reset
			if out == nil {
				t.iterNext = nil
			}
			return out
		}
	case PostOrder:
		stack = append(stack, current)
		var prevRBNode *RBNode

		t.iterNext = func() (out Interface) {
			for len(stack) > 0 {
				// peek
				stackIndex := len(stack) - 1
				current = stack[stackIndex]
				if (prevRBNode == nil) ||
					(prevRBNode.left == current) ||
					(prevRBNode.right == current) {
					if current.left != nil {
						stack = append(stack, current.left)
					} else if current.right != nil {
						stack = append(stack, current.right)
					}
				} else if current.left == prevRBNode {
					if current.right != nil {
						stack = append(stack, current.right)
					}
				} else {
					out = current.Elem
					// pop, but no assignment
					stackIndex := len(stack) - 1
					stack = stack[0:stackIndex]
					prevRBNode = current
					break
				}
				prevRBNode = current
			}

			// last node, reset
			if out == nil {
				t.iterNext = nil
			}
			return out

		}
	default:
		s := fmt.Sprintf("rbTree has not implemented %s for iteration.", order)
		panic(s)
	}
	// return our first node
	return t.iterNext()

}

// Map is a more performance orientated way to iterate over the elements of the tree.
// Given a TravOrder and a function which conforms to the IterFunc type:
//
//      type IterFunc func(Interface)
//
// Map calls the function for each RBNode  in the specified order.
func (t *RBTree) Map(order TravOrder, f IterFunc) {

	n := t.root
	switch order {
	case InOrder:
		var inorder func(node *RBNode)
		inorder = func(node *RBNode) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node.Elem)
			inorder(node.right)
		}
		inorder(n)
	case PreOrder:
		var preorder func(node *RBNode)
		preorder = func(node *RBNode) {
			if node == nil {
				return
			}
			f(node.Elem)
			preorder(node.left)
			preorder(node.right)
		}
		preorder(n)
	case PostOrder:
		var postorder func(node *RBNode)
		postorder = func(node *RBNode) {
			if node == nil {
				return
			}
			postorder(node.left)
			postorder(node.right)
			f(node.Elem)
		}
		postorder(n)
	default:
		s := fmt.Sprintf("rbTree has not implemented %s.", order)
		panic(s)
	}

}

// Search returns the matching item if found, otherwise nil is returned.
func (t *RBTree) Search(item Interface) (found Interface) {
	if item == nil {
		return
	}
	x := t.root
	for x != nil {
		switch x.Elem.Compare(item) {
		case EQ:
			found = x.Elem
			return
		case GT:
			x = x.left
		case LT:
			x = x.right
		}
	}
	return
}

// Insert will either insert a new entry into the tree, and return nil. Or if there was a previous entry already inserted, then in addition to inserting the new item, the previously inserted item will be returned.
func (t *RBTree) Insert(item Interface) (old Interface) {
	if item == nil {
		return
	}

	if t.root == nil {
		t.size++
		t.root = &RBNode{Elem: item, left: nil, right: nil}
		t.first = t.root
		t.last = t.root
	} else {
		t.root, old = t.insert(t.root, item)
	}

	if t.root.color == red {
		t.height++
	}
	t.root.color = black // maintain rb invariants
	return
}

func (t *RBTree) insert(h *RBNode, item Interface) (root *RBNode, old Interface) {
	if h == nil {
		t.size++
		// base case, insert do stuff on new node
		n := &RBNode{Elem: item, left: nil, right: nil}
		// set Min
		switch t.first.Elem.Compare(item) {
		case GT:
			t.first = n
		}
		// set Max
		switch t.last.Elem.Compare(item) {
		case LT:
			t.last = n
		}
		root = n
		return
	}

	switch h.Elem.Compare(item) {
	case GT:
		h.left, old = t.insert(h.left, item)
	case LT:
		h.right, old = t.insert(h.right, item)
	case EQ:
		old = h.Elem
		h.Elem = item
	}

	if h.right.isred() && !(h.left.isred()) {
		h = h.rotateLeft()
	}
	if h.left.isred() && h.left.left.isred() {
		h = h.rotateRight()
	}

	if h.left.isred() && h.right.isred() {
		h.colorFlip()
	}
	root = h
	return
}

// Remove looks for a matching entry, and if found, the item is removed from the tree and old is populated with the removed item. If the item is not matched in the tree, nil is returned.
func (t *RBTree) Remove(item Interface) (old Interface) {
	if item == nil || t.root == nil {
		return
	}
	t.root, old = t.remove(t.root, item)
	if old != nil {
		if t.root == nil {
			t.first = nil
			t.last = nil
		} else {
			// set Min
			switch t.first.Elem.Compare(old) {
			case EQ:
				t.first = t.root.min()
			}
			// set Max
			switch t.last.Elem.Compare(old) {
			case EQ:
				t.last = t.root.max()

			}

		}
	} else {

	}
	if t.root != nil && t.root.color == red {
		t.root.color = black // maintain rb invariants
		t.height--
	} else if t.root == nil {
		t.height--
	}
	return

}

func (t *RBTree) remove(h *RBNode, item Interface) (root *RBNode, old Interface) {

	switch h.Elem.Compare(item) {
	case LT, EQ:
		if h.left.isred() {
			h = h.rotateRight()
		}
		if result := h.Elem.Compare(item); result == EQ && h.right == nil {
			t.size--
			old = h.Elem
			h = nil
			root = nil
			return
		}
		if h.right != nil {
			if !h.right.isred() && !(h.right.left.isred()) {
				h = h.moveredRight()
			}
			if result := h.Elem.Compare(item); result == EQ {
				old = h.Elem
				t.size--

				x := h.right.min()
				h.Elem = x.Elem
				h.right = h.right.removeMin()
			} else {
				h.right, old = t.remove(h.right, item)
			}
		}
	case GT:
		if h.left != nil {
			if !h.left.isred() && !(h.left.left.isred()) {
				h = h.moveredLeft()
			}
			h.left, old = t.remove(h.left, item)
		}

	}
	root = h.fixUp()
	return
}

// Left Leaning red black Tree functions and helpers to maintain public methods

func (h *RBNode) rotateLeft() (x *RBNode) {
	x = h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = red
	return
}

func (h *RBNode) rotateRight() (x *RBNode) {
	x = h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = red
	return
}

func (h *RBNode) isred() bool {
	return h != nil && h.color == red
}

func (h *RBNode) moveredLeft() *RBNode {
	h.colorFlip()
	if h.right.left.isred() {
		h.right = h.right.rotateRight()
		h = h.rotateLeft()
		h.colorFlip()
	}
	return h
}

func (h *RBNode) moveredRight() *RBNode {
	h.colorFlip()
	if h.left.left.isred() {
		h = h.rotateRight()
		h.colorFlip()
	}
	return h
}

func (h *RBNode) colorFlip() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

func (h *RBNode) fixUp() *RBNode {
	if h.right.isred() {
		h = h.rotateLeft()
	}

	if h.left.isred() && h.left.left.isred() {
		h = h.rotateRight()
	}
	if h.left.isred() && h.right.isred() {
		h.colorFlip()
	}
	return h
}
func (h *RBNode) min() *RBNode {
	for ; h.left != nil; h = h.left {
	}
	return h
}
func (h *RBNode) max() *RBNode {
	for ; h.right != nil; h = h.right {
	}
	return h
}

func (h *RBNode) removeMin() *RBNode {
	if h.left == nil {
		return nil
	}
	if !h.left.isred() && !h.left.left.isred() {
		h = h.moveredLeft()
	}

	h.left = h.left.removeMin()

	return h.fixUp()
}
