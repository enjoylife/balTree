package gotree

import (
	"fmt"
)

type color bool

const (
	Red   color = false //default
	Black color = true
)

// Pretty output for errors, debuging, etc.
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

// Exports a key and value
type Node struct {
	Key, Value  interface{}
	left, right *Node
	color       color
}

// internal
type rbIter struct {
	current *Node
	stack   []*Node
	next    func() *Node
}

// Exports a Height, Size, and
type RBTree struct {
	Height      int
	Size        int
	first, last *Node
	iter        rbIter
	root        *Node
	cmp         CompareFunc
}

// creates a RBTree for use.
// Must give it a compare function, see common.go for an example
func New(cmp CompareFunc) *RBTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &RBTree{root: nil, first: nil, last: nil, cmp: cmp}
}

func (t *RBTree) Min() *Node {
	return t.first
}

func (t *RBTree) Max() *Node {
	return t.last
}

func (t *RBTree) Next() *Node {

	n := t.iter.next()
	if n == nil {
	}
	return n

}

func (t *RBTree) InitIter(order TravOrder) *Node {
	current := t.root
	stack := []*Node{}
	switch order {
	case InOrder:
		t.iter.next = func() (out *Node) {
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
			return out
		}

	case PreOrder:
		t.iter.next = func() (out *Node) {
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
			return out
		}
	case PostOrder:
		if current == nil {
			return current
		}
		stack = append(stack, current)
		var prevNode *Node = nil

		t.iter.next = func() (out *Node) {
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
			return out

		}
	default:
		s := fmt.Sprintf("rbTree has not implemented %s for iteration.", order)
		panic(s)
	}
	// return our first node
	return t.iter.next()

}

func (t *RBTree) Traverse(order TravOrder, f IterFunc) {

	n := t.root
	switch order {
	case InOrder:
		var inorder func(node *Node)
		inorder = func(node *Node) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node.Key, node.Value)
			inorder(node.right)
		}
		inorder(n)
	case PreOrder:
		var preorder func(node *Node)
		preorder = func(node *Node) {
			if node == nil {
				return
			}
			f(node.Key, node.Value)
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
			f(node.Key, node.Value)
		}
		postorder(n)
	default:
		s := fmt.Sprintf("rbTree has not implemented %s.", order)
		panic(s)
	}

}

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

func (t *RBTree) Search(key interface{}) (value interface{}, ok bool) {
	if key == nil {
		return
	}
	x := t.root
	for x != nil {

		switch t.cmp(x.Key, key) {
		case EQ:
			return x.Value, true
		case LT:
			x = x.left
		case GT:
			x = x.right
		default:
			panic("Compare result of undefined")
		}
	}
	return
}

// Get an Item in the tree.
func (t *RBTree) SearchRecurse(key interface{}) (value interface{}, ok bool) {

	if key == nil {
		return
	}
	return t.get(t.root, key)
}

func (t *RBTree) get(node *Node, key interface{}) (value interface{}, ok bool) {
	if node == nil {
		return nil, false
	}

	switch t.cmp(node.Key, key) {
	case EQ:
		return node.Value, true
	case LT:
		return t.get(node.left, key)
	case GT:
		return t.get(node.right, key)
	default:
		panic("Compare result of undefined")
	}
}

func (t *RBTree) Insert(key interface{}, value interface{}) (old interface{}, ok bool) {
	if key == nil {
		return
	}

	if t.root == nil {
		t.root = &Node{color: Red, Key: key, Value: value, left: nil, right: nil}
		t.first = t.root
		t.last = t.root
	} else {
		t.root = t.insert(t.root, key, value)
		//t.root = t.insertIter(t.root, key, value)
	}
	if t.root.color == Red {
		t.Height++
	}
	t.root.color = Black // maintain rb invariants
	return
}

func (t *RBTree) insert(h *Node, key interface{}, value interface{}) *Node {
	if h == nil {
		// base case, insert do stuff on new node
		n := &Node{color: Red, Key: key, Value: value}
		// set Min
		switch t.cmp(t.first.Key, key) {
		case GT:
			t.first = n
		}
		// set Max
		switch t.cmp(t.last.Key, key) {
		case LT:
			t.last = n
		}
		return n
	}

	switch t.cmp(h.Key, key) {
	case EQ:
		h.Value = value
	case LT:
		h.left = t.insert(h.left, key, value)
	case GT:
		h.right = t.insert(h.right, key, value)
	default:
		panic("Compare result of undefined")
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
	return h
}

func (h *Node) isRed() bool {
	return h != nil && h.color == Red
}
func (h *Node) colorFlip() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

/* Experinmental */
//Return a complete copy(view) of tree to work on.
// TODO Deep copy
func (t *RBTree) CopyiedSliceIter(order TravOrder) []*Node {
	tCopy := []*Node{}
	f := func(key interface{}, value interface{}) {
		tCopy = append(tCopy, &Node{Key: key, Value: value})
	}
	switch order {
	case InOrder:
		t.Traverse(InOrder, f)
	case PreOrder:
		t.Traverse(PreOrder, f)
	case PostOrder:
		t.Traverse(PostOrder, f)
	default:
		s := fmt.Sprintf("rbTree has not implemented %s for iteration.", order)
		panic(s)
	}
	t.iter.next = func() (out *Node) {
		return
	}
	// return our first node
	return tCopy

}
