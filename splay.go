package gotree

import (
	"fmt"
)

var _ = fmt.Printf

type SplayNode struct {
	Elem        Interface
	left, right *SplayNode
}

type SplayTree struct {
	Size        int // Number of inserted elements
	first, last *SplayNode
	iterNext    func() *SplayNode // initially nil
	root        *SplayNode
}

// Search returns the matching item if found, otherwise nil is returned.
func (t *SplayTree) Search(item Interface) (found Interface) {
	if item == nil {
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
		t.Size++
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
		t.Size++
	case LT:
		n = &SplayNode{Elem: item, left: t.root, right: t.root.right}
		t.root.right = nil
		t.root = n
		t.Size++
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
			x = t.root.left
		} else {
			x = t.root.left.splay(item)
			x.right = t.root.right
		}
		t.Size--
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

type SplayIterFunc func(*SplayNode)

// Next is called when individual elements are wanted to be traversed over.
// Prior to a call to Next, a call to IterInit needs to be made to set up the necessary
// data to allow for traversal of the tree. Example:
//
//    sum := 0
//    for i, n := 0, tree.IterInit(InOrder); n != nil; i, n = i+1, tree.Next() {
//        elem := n.Elem.(exInt)  // (exInt is simple int type)
//        sum += int(elem)
//    }
// Note: If one was to break out of the loop prior to a complete traversal,
// and start another loop without calling IterInit, then the previously uncompleted iterator is continued again.
func (t *SplayTree) Next() *SplayNode {

	if t.iterNext == nil {
		return nil
	}
	return t.iterNext() // func set by call to IterInit(TravOrder)

}

// IterInit is the initializer which setups the tree for iterating over it's elements in
// a specific order. It setups the internal data, and then returns the first SplayNode to be looked at. See Next for an example.
func (t *SplayTree) IterInit(order TravOrder) *SplayNode {

	current := t.root
	stack := []*SplayNode{}
	switch order {
	case InOrder:
		t.iterNext = func() (out *SplayNode) {
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
		t.iterNext = func() (out *SplayNode) {
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
		var prevSplayNode *SplayNode = nil

		t.iterNext = func() (out *SplayNode) {
			for len(stack) > 0 {
				// peek
				stackIndex := len(stack) - 1
				current = stack[stackIndex]
				if (prevSplayNode == nil) ||
					(prevSplayNode.left == current) ||
					(prevSplayNode.right == current) {
					if current.left != nil {
						stack = append(stack, current.left)
					} else if current.right != nil {
						stack = append(stack, current.right)
					}
				} else if current.left == prevSplayNode {
					if current.right != nil {
						stack = append(stack, current.right)
					}
				} else {
					out = current
					// pop, but no assignment
					stackIndex := len(stack) - 1
					stack = stack[0:stackIndex]
					prevSplayNode = current
					break
				}
				prevSplayNode = current
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
// Given a TravOrder and a function which conforms to the RBSplayIterFunc type:
//
//      type SplayIterFunc func(*SplayNode)
//
// Map calls the function for each SplayNode  in the specified order.
func (t *SplayTree) Map(order TravOrder, f SplayIterFunc) {

	n := t.root
	switch order {
	case InOrder:
		var inorder func(node *SplayNode)
		inorder = func(node *SplayNode) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node)
			inorder(node.right)
		}
		inorder(n)
	case PreOrder:
		var preorder func(node *SplayNode)
		preorder = func(node *SplayNode) {
			if node == nil {
				return
			}
			f(node)
			preorder(node.left)
			preorder(node.right)
		}
		preorder(n)
	case PostOrder:
		var postorder func(node *SplayNode)
		postorder = func(node *SplayNode) {
			if node == nil {
				return
			}
			postorder(node.left)
			postorder(node.right)
			f(node)
		}
		postorder(n)
	default:
		s := fmt.Sprintf("SplayTree has not implemented %s.", order)
		panic(s)
	}

}

// Height returns the max depth of any branch of the tree
func (t *SplayTree) Height() int {
	var calc func(n *SplayNode) int
	calc = func(n *SplayNode) int {
		if n == nil {
			return 0
		}
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
