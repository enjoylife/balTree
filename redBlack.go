package gotree

import (
	"fmt"
)

// Color is the used to maintain the redblack tree balance.
type color bool

const (
	Red   color = false // we rely on default for node initializations
	Black color = true
)

// Pretty output for errors, debugging, etc.
func (c color) String() string {
	var s string
	switch c {
	case Red:
		s = "red"
	case Black:
		s = "black"
	}
	return s
}

// A Node is the type manipulated within the tree. It holds the inserted elements.
// It is exposed whenever the tree traversal functions are used.
type Node struct {
	Elem Interface
	//
	left, right *Node
	color       color
}

// A RBTree is our main type our redblack tree methods are defined on.
type RBTree struct {
	Height      int // Height from root to leaf
	Size        int // Number of inserted elements
	first, last *Node
	iterNext    func() *Node // initially ni
	root        *Node
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
// Prior to a call to Next, a call to InitIter needs to be made to set up the necessary
// data to allow for traversal of the tree. Example:
//
//      // (exInt is simple int type)
//  	sum := 0
//  	for i, n := 0, tree.InitIter(PreOrder); n != nil; i, n = i+1, tree.Next() {
//  		sum += int(n.Elem.(exInt)) + i
//  	}
// Note: If one was to break out of the loop prior to a complete traversal,
// and start another loop without calling InitIter, then the previously uncompleted iterator is continued again.
func (t *RBTree) Next() *Node {

	if t.iterNext == nil {
		return nil
	}
	return t.iterNext() // func set by call to InitIter(TravOrder)

}

// InitIter is the initializer which setups the tree for iterating over it's elements in
// a specific order. It setups the internal data, and then returns the first Node to be looked at. See Next for an example.
func (t *RBTree) InitIter(order TravOrder) *Node {

	current := t.root
	stack := []*Node{}
	switch order {
	case InOrder:
		t.iterNext = func() (out *Node) {
			for len(stack) > 0 || current != nil {
				if current != nil {
					stack = append(stack, current)
					current = current.left
				} else {
					// pop
					stackIndex := len(stack) - 1
					out = stack[stackIndex]
					current = out
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
		t.iterNext = func() (out *Node) {
			for len(stack) > 0 || current != nil {
				if current != nil {
					out = current
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
		if current == nil {
			return current
		}
		stack = append(stack, current)
		var prevNode *Node = nil

		t.iterNext = func() (out *Node) {
			for len(stack) > 0 {
				// peek
				stackIndex := len(stack) - 1
				current = stack[stackIndex]
				if (prevNode == nil) ||
					(prevNode.left == current) ||
					(prevNode.right == current) {
					if current.left != nil {
						stack = append(stack, current.left)
					} else if current.right != nil {
						stack = append(stack, current.right)
					}
				} else if current.left == prevNode {
					if current.right != nil {
						stack = append(stack, current.right)
					}
				} else {
					out = current
					// pop, but no assignment
					stackIndex := len(stack) - 1
					stack = stack[0:stackIndex]
					prevNode = current
					break
				}
				prevNode = current
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
//      type IterFunc func(*Node)
//
// Map calls the function for each Node  in the specified order.
func (t *RBTree) Map(order TravOrder, f IterFunc) {

	n := t.root
	switch order {
	case InOrder:
		var inorder func(node *Node)
		inorder = func(node *Node) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node)
			inorder(node.right)
		}
		inorder(n)
	case PreOrder:
		var preorder func(node *Node)
		preorder = func(node *Node) {
			if node == nil {
				return
			}
			f(node)
			preorder(node.left)
			preorder(node.right)
		}
		preorder(n)
	case PostOrder:
		var postorder func(node *Node)
		postorder = func(node *Node) {
			if node == nil {
				return
			}
			postorder(node.left)
			postorder(node.right)
			f(node)
		}
		postorder(n)
	default:
		s := fmt.Sprintf("rbTree has not implemented %s.", order)
		panic(s)
	}

}

// Search takes as input any type implementing the Interface interface and returns either:
// a matching Interface element as based upon that types Compare function along wih a nil error.
// If given an item which can't be successfully compared within the array, found is returned with a nil, and
// error is set to InvalidInterfaceError.
// If a search within the tree comes up empty, found is nil, but error is populated with a NonexistentElemError.
func (t *RBTree) Search(item Interface) (found Interface, err error) {
	if item == nil {
		var e InvalidInterfaceError
		return nil, e
	}
	x := t.root
	for x != nil {
		switch x.Elem.Compare(item) {
		case GT:
			x = x.left
		case LT:
			x = x.right
		case EQ:
			return x.Elem, nil
		}
	}
	var e NonexistentElemError
	return nil, e
}

// Insert takes a type implementing the Interface interface, this type is then inserted into the
// tree. If there was a previous entry at the same insertion point as the item to be inserted,
// the old element is returned.
// If given an item which can't be successfully compared within the array, old is returned with a nil, and
// error is set to InvalidInterfaceError.
func (t *RBTree) Insert(item Interface) (old Interface, err error) {
	if item == nil {
		var err InvalidInterfaceError
		return nil, err
	}

	if t.root == nil {
		t.Size++
		t.root = &Node{Elem: item, left: nil, right: nil}
		t.first = t.root
		t.last = t.root
	} else {
		t.root, old = t.insert(t.root, item)
	}

	if t.root.color == Red {
		t.Height++
	}
	t.root.color = Black // maintain rb invariants
	return
}

func (t *RBTree) insert(h *Node, item Interface) (root *Node, old Interface) {
	if h == nil {
		t.Size++
		// base case, insert do stuff on new node
		n := &Node{Elem: item, left: nil, right: nil}
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
		return n, nil
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

	if h.right.isRed() && !(h.left.isRed()) {
		h = h.rotateLeft()
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = h.rotateRight()
	}

	if h.left.isRed() && h.right.isRed() {
		h.colorFlip()
	}
	return h, old
}

// Remove takes a type implementing the Interface interface, this type is then searched on inside the tree.
// If a matching entry is found the item is removed from the tree and old is populated with said removed item. error is nil in this case.
// If when searching within the tree comes up empty, old is nil, but error is populated with a NonexistentElemError.
func (t *RBTree) Remove(item Interface) (old Interface, err error) {
	if item == nil || t.root == nil {
		var err InvalidInterfaceError
		return nil, err
	}
	t.root, old, err = t.remove(t.root, item)
	if err != nil {
		return nil, err
	} else if old != nil {
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
	if t.root != nil && t.root.color == Red {
		t.root.color = Black // maintain rb invariants
		t.Height--
	} else if t.root == nil {
		t.Height--
	}
	return old, nil

}

// TODO Test error returns
func (t *RBTree) remove(h *Node, item Interface) (root *Node, old Interface, err error) {

	var e NonexistentElemError
	switch h.Elem.Compare(item) {
	case LT, EQ:
		if h.left.isRed() {
			h = h.rotateRight()
		}
		if result := h.Elem.Compare(item); result == EQ && h.right == nil {
			old = h.Elem
			t.Size--
			return nil, old, nil
		}
		if h.right != nil {
			if !h.right.isRed() && !(h.right.left.isRed()) {
				h = h.moveRedRight()
			}
			if result := h.Elem.Compare(item); result == EQ {
				old = h.Elem
				t.Size--

				x := h.right.min()
				h.Elem = x.Elem
				h.right = h.right.removeMin()
			} else {
				h.right, old, err = t.remove(h.right, item)
			}
		} else {

			err = e
		}
	case GT:
		if h.left != nil {
			if !h.left.isRed() && !(h.left.left.isRed()) {
				h = h.moveRedLeft()
			}
			h.left, old, err = t.remove(h.left, item)
		} else {
			err = e
		}
	}
	h = h.fixUp()
	return h, old, err
}

// Left Leaning Red Black Tree functions and helpers to maintain public methods

func (h *Node) rotateLeft() (x *Node) {
	x = h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = Red
	return
}

func (h *Node) rotateRight() (x *Node) {
	x = h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = Red
	return
}

func (h *Node) isRed() bool {
	return h != nil && h.color == Red
}

func (h *Node) moveRedLeft() *Node {
	h.colorFlip()
	if h.right.left.isRed() {
		h.right = h.right.rotateRight()
		h = h.rotateLeft()
		h.colorFlip()
	}
	return h
}

func (h *Node) moveRedRight() *Node {
	h.colorFlip()
	if h.left.left.isRed() {
		h = h.rotateRight()
		h.colorFlip()
	}
	return h
}

func (h *Node) colorFlip() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

func (h *Node) fixUp() *Node {
	if h.right.isRed() {
		h = h.rotateLeft()
	}

	if h.left.isRed() && h.left.left.isRed() {
		h = h.rotateRight()
	}
	if h.left.isRed() && h.right.isRed() {
		h.colorFlip()
	}
	return h
}
func (h *Node) min() *Node {
	for ; h.left != nil; h = h.left {
	}
	return h
}
func (h *Node) max() *Node {
	for ; h.right != nil; h = h.right {
	}
	return h
}

func (h *Node) removeMin() *Node {
	if h.left == nil {
		return nil
	}
	if !h.left.isRed() && !h.left.left.isRed() {
		h = h.moveRedLeft()
	}

	h.left = h.left.removeMin()

	return h.fixUp()
}
