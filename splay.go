package gotree

import (
	"fmt"
)

var _ = fmt.Printf

type Node struct {
	Elem        Interface
	left, right *Node
}

type Tree struct {
	Size        int // Number of inserted elements
	first, last *Node
	iterNext    func() *Node // initially nil
	root        *Node
}

func (t *Tree) Search(item Interface) (found Interface, err error) {
	if item == nil {
		var e InvalidInterfaceError
		return nil, e
	}
	t.root = t.root.splay(item)
	switch t.root.Elem.Compare(item) {
	case EQ:
		return t.root.Elem, nil
	}
	return nil, nil
}
func (t *Tree) Insert(item Interface) (old Interface, err error) {
	var n *Node
	// TODO min and max update
	if item == nil {
		var err InvalidInterfaceError
		return nil, err
	}

	if t.root == nil {
		t.Size++
		t.root = &Node{Elem: item, left: nil, right: nil}
		t.first = t.root
		t.last = t.root
		return
	}
	t.root = t.root.splay(item)
	switch t.root.Elem.Compare(item) {
	case GT:
		n = &Node{Elem: item, left: t.root.left, right: t.root}
		t.root.left = nil
		t.root = n
		t.Size++
	case LT:
		n = &Node{Elem: item, left: t.root, right: t.root.right}
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
func (t *Tree) Remove(item Interface) (old Interface, err error) {
	var x *Node
	if item == nil {
		var err InvalidInterfaceError
		return nil, err
	}
	if t.root == nil {
		var err NonexistentElemError
		return nil, err
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
	return old, nil

}

// Min returns the smallest inserted element if possible. If the smallest value is not
// found(empty tree), then Min returns a nil.
func (t *Tree) Min() Interface {
	if t.first != nil {
		return t.first.Elem
	}
	return nil
}

// Max returns the largest inserted element if possible. If the largest value is not
// found(empty tree), then Max returns a nil.
func (t *Tree) Max() Interface {
	if t.last != nil {
		return t.last.Elem
	}
	return nil
}

type IterFunc func(*Node)

func (t *Tree) Next() *Node {

	if t.iterNext == nil {
		return nil
	}
	return t.iterNext() // func set by call to InitIter(TravOrder)

}

// InitIter is the initializer which setups the tree for iterating over it's elements in
// a specific order. It setups the internal data, and then returns the first Node to be looked at. See Next for an example.
func (t *Tree) InitIter(order TravOrder) *Node {

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
// Given a TravOrder and a function which conforms to the RBIterFunc type:
//
//      type IterFunc func(*Node)
//
// Map calls the function for each Node  in the specified order.
func (t *Tree) Map(order TravOrder, f IterFunc) {

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
		s := fmt.Sprintf("Tree has not implemented %s.", order)
		panic(s)
	}

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

func (t *Node) splay(item Interface) (out *Node) {
	var left, right, parent *Node
	var n Node
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
