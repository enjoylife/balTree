package gotree

import (
	"fmt"
	"sync"
)

type color bool

const (
	Red   color = true
	Black color = false
)

/* For pretty output when debuging */
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
type rbNode struct {
	key, value  interface{}
	left, right *rbNode
	color       color
}

func (n *rbNode) Key() interface{} {
	return n.key
}
func (n *rbNode) Value() interface{} {
	return n.value
}

// internal
type rbIter struct {
	current *rbNode
	stack   []*rbNode
	next    func() *rbNode
}

// Exports a Height, Size, and
type RbTree struct {
	Height      int
	Size        int
	first, last *rbNode
	iter        rbIter
	root        *rbNode
	cmp         CompareFunc
	lock        sync.RWMutex
}

// creates a RbTree for use.
// Must give it a compare function, see common.go for an example
func New(cmp CompareFunc) *RbTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &RbTree{root: nil, first: nil, last: nil, cmp: cmp}
}

func (t *RbTree) Min() Node {
	return t.first
}

func (t *RbTree) Max() Node {
	return t.last
}

func (t *RbTree) Next() *rbNode {

	return (t.iter.next())
}

func (t *RbTree) InitIter(order TravOrder) *rbNode {
	current := t.root
	stack := []*rbNode{}
	switch order {
	case InOrder:
		t.iter.next = func() (out *rbNode) {
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
		t.iter.next = func() (out *rbNode) {
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
		var prevNode *rbNode = nil

		t.iter.next = func() (out *rbNode) {
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

func (t *RbTree) Traverse(order TravOrder, f IterFunc) {

	n := t.root
	switch order {
	case InOrder:
		var inorder func(node *rbNode)
		inorder = func(node *rbNode) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node.key, node.value)
			inorder(node.right)
		}
		inorder(n)
	case PreOrder:
		var preorder func(node *rbNode)
		preorder = func(node *rbNode) {
			if node == nil {
				return
			}
			f(node.key, node.value)
			preorder(node.left)
			preorder(node.right)
		}
		preorder(n)
	case PostOrder:
		var postorder func(node *rbNode)
		postorder = func(node *rbNode) {
			if node == nil {
				return
			}
			postorder(node.left)
			postorder(node.right)
			f(node.key, node.value)
		}
		postorder(n)
	default:
		s := fmt.Sprintf("rbTree has not implemented %s.", order)
		panic(s)
	}

}

func (h *rbNode) rotateLeft() (x *rbNode) {
	x = h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = Red
	return
}

func (h *rbNode) rotateRight() (x *rbNode) {
	x = h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = Red
	return
}

func (t *RbTree) Search(key interface{}) (value interface{}, ok bool) {
	if key == nil {
		return
	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	x := t.root
	for x != nil {

		switch t.cmp(x.key, key) {
		case EQ:
			return x.value, true
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

func (t *RbTree) Insert(key interface{}, value interface{}) (old interface{}, ok bool) {
	if key == nil {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.root == nil {
		t.root = &rbNode{color: Red, key: key, value: value, left: nil, right: nil}
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

func (t *RbTree) insert(h *rbNode, key interface{}, value interface{}) *rbNode {
	if h == nil {
		// base case, insert do stuff on new node
		n := &rbNode{color: Red, key: key, value: value}
		// set Min
		switch t.cmp(t.first.key, key) {
		case GT:
			t.first = n
		}
		// set Max
		switch t.cmp(t.last.key, key) {
		case LT:
			t.last = n
		}
		return n
	}

	switch t.cmp(h.key, key) {
	case EQ:
		h.value = value
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

func (h *rbNode) isRed() bool {
	return h != nil && h.color == Red
}
func (h *rbNode) colorFlip() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
