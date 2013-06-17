package gotree

import (
	"fmt"
	"runtime"
)

var _ = fmt.Printf

type SplayNode struct {
	Elem        Interface
	left, right *SplayNode
}

type SplayTree struct {
	size        int // Number of inserted elements
	first, last *SplayNode
	iterNext    func() Interface // initially nil
	root        *SplayNode
}

func (t *SplayTree) Clear() {
	t.root = nil
	t.last = nil
	t.first = nil
	t.size = 0
	t.iterNext = nil
	runtime.GC()
}

// Search returns the matching item if found, otherwise nil is returned.
func (t *SplayTree) Search(item Interface) (found Interface) {
	if item == nil || t.root == nil {
		return
	}
	t.root = t.root.splay(item)
	switch t.root.Elem.Compare(item) {
	case EQ:
		return t.root.Elem
	}
	return
}

// Insert will either insert a new entry into the tree, and return nil. Or if there was a previous entry already inserted, then in addition to inserting the new item, the previously inserted item will be returned.
func (t *SplayTree) Insert(item Interface) (old Interface) {
	var n *SplayNode
	// TODO min and max update
	if item == nil {
		return nil
	}

	if t.root == nil {
		t.size++
		t.root = &SplayNode{Elem: item, left: nil, right: nil}
		t.first = t.root
		t.last = t.root
		return
	}
	t.root = t.root.splay(item)
	switch t.root.Elem.Compare(item) {
	case GT:
		n = &SplayNode{Elem: item, left: t.root.left, right: t.root}
		t.root.left = nil
		t.root = n
		t.size++
	case LT:
		n = &SplayNode{Elem: item, left: t.root, right: t.root.right}
		t.root.right = nil
		t.root = n
		t.size++
	case EQ:
		old = t.root.Elem
		t.root.Elem = item

	}
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
	return
}

// Remove looks for a matching entry, and if found, the item is removed from the tree and old is populated with the removed item. If the item is not matched in the tree, nil is returned.
func (t *SplayTree) Remove(item Interface) (old Interface) {
	var x *SplayNode
	if item == nil || t.root == nil {
		return
	}

	t.root = t.root.splay(item)

	switch t.root.Elem.Compare(item) {
	// TODO NP case
	case EQ:
		old = t.root.Elem
		if t.root.left == nil {
			x = t.root.right
		} else {
			x = t.root.left.splay(item)
			x.right = t.root.right
		}
		t.root = x
		t.size--
		if t.root != nil {
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
		} else {
			t.first = nil
			t.last = nil
		}

	}
	return old

}

// Min returns the smallest inserted element if possible. If the smallest value is not
// found(empty tree), then Min returns a nil.
func (t *SplayTree) Min() Interface {
	if t.first != nil {
		return t.first.Elem
	}
	return nil
}

// Max returns the largest inserted element if possible. If the largest value is not
// found(empty tree), then Max returns a nil.
func (t *SplayTree) Max() Interface {
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
func (t *SplayTree) Next() Interface {

	if t.iterNext == nil {
		return nil
	}
	return t.iterNext() // func set by call to IterInit(TravOrder)

}

// IterInit is the initializer which setups the tree for iterating over it's elements in
// a specific order. It setups the internal data, and then returns the first Interface to be looked at. See Next for an example.
func (t *SplayTree) IterInit(order TravOrder) Interface {

	current := t.root
	stack := []*SplayNode{}
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
	default:
		s := fmt.Sprintf("rbSplayTree has not implemented %s for iteration.", order)
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
// Map calls the function for each Interface type in the specified order.
func (t *SplayTree) Map(order TravOrder, f IterFunc) {

	if t.root == nil {
		return
	}
	n := t.root
	switch order {
	case InOrder:
		var inorder func(node *SplayNode)
		inorder = func(node *SplayNode) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node.Elem)
			inorder(node.right)
		}
		inorder(n)
	default:
		s := fmt.Sprintf("SplayTree has not implemented %s.", order)
		panic(s)
	}

}

func (t *SplayTree) Size() int {
	return t.size
}

// Height returns the max depth of any branch of the tree.
// Note: Runs in O(n) where n is the maximum depthed branch.
func (t *SplayTree) Height() int {
	var calc func(n *SplayNode) int
	calc = func(n *SplayNode) int {
		if n == nil {
			return 0
		}
		// math.Max for int
		if a, b := calc(n.left), calc(n.right); a >= b {
			return 1 + a
		} else {
			return 1 + b
		}
	}
	return calc(t.root)
}

func (h *SplayNode) min() *SplayNode {
	for ; h.left != nil; h = h.left {
	}
	return h
}
func (h *SplayNode) max() *SplayNode {
	for ; h.right != nil; h = h.right {
	}
	return h
}

func (t *SplayNode) splay(item Interface) (out *SplayNode) {
	var left, right, parent *SplayNode
	var n SplayNode
	left = &n
	right = &n

L:
	for {
		switch t.Elem.Compare(item) {
		//TODO NP case
		case GT:
			//fmt.Println("Madit LEft")
			if t.left == nil {
				break L
			}
			switch t.left.Elem.Compare(item) {
			//TODO NP case
			case GT:
				// rotate right
				parent = t.left
				t.left = parent.right
				parent.right = t
				t = parent
				if t.left == nil {
					break L
				}
			}
			// link right
			right.left = t
			right = t
			t = t.left
		case LT:
			if t.right == nil {
				//fmt.Println("Madit Right")
				break L
			}
			switch t.right.Elem.Compare(item) {
			case LT:
				// rotate left
				parent = t.right
				t.right = parent.left
				parent.left = t
				t = parent
				if t.right == nil {
					break L
				}
			}
			// link left
			left.right = t
			left = t
			t = t.right
		case EQ:
			break L
		}
	}
	// assemble
	left.right = t.left
	right.left = t.right
	t.left = n.right
	t.right = n.left
	return t
}
