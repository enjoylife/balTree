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
	iterChan    chan *Node
	cmp         CompareFunc
}

// creates a RBTree for use.
// Must give it a compare function, see common.go for an example
func New(cmp CompareFunc) *RBTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &RBTree{root: nil, first: nil, last: nil, Height: 0, cmp: cmp}
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

func (t *RBTree) Insert(key interface{}, value interface{}) (old interface{}, ok bool) {
	if key == nil {
		return
	}

	if t.root == nil {
		t.Size++
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
		t.Size++
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
func (h *Node) moveRedLeft() *Node {
	h.colorFlip()
	if h.right.left.isRed() {
		h.right = h.right.rotateRight()
		h = h.rotateLeft()
		h.colorFlip()
		/*if h.right.right.isRed() {
			h.right = h.right.rotateLeft()
		}*/
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

func (h *Node) removeMin() *Node {
	if h == nil {
		panic("WE GOT A NIL")
	}
	if h.left == nil {
		return nil
	}
	if !h.left.isRed() && !h.left.left.isRed() {
		h = h.moveRedLeft()
	}

	h.left = h.left.removeMin()
	return h.fixUp()
}
func (h *Node) fixUp() *Node {
	//fmt.Println("#1", h)
	/*if h.right.isRed() && !h.left.isRed() {
		h = h.rotateLeft()
	}*/
	if h.right.isRed() {
		h = h.rotateLeft()
	}

	//fmt.Println("#2", h)

	if h.left.isRed() && h.left.left.isRed() {
		h = h.rotateRight()
	}
	//fmt.Println("#3", h)
	if h.left.isRed() && h.right.isRed() {
		h.colorFlip()
	}
	return h
}

func (h *Node) min() *Node {
	for h.left != nil {
		h = h.left
	}
	return h
}

func (t *RBTree) Remove(key interface{}) (ok bool) {
	if key == nil {
		return
	}
	if _, check := t.Search(key); !check {
		return
	}
	t.root = t.remove(t.root, key)
	if t.root != nil && t.root.color == Red {
		t.root.color = Black // maintain rb invariants
		t.Height--
	} else if t.root == nil {
		t.Height--
	}
	return true

}

func (t *RBTree) remove(h *Node, key interface{}) *Node {

	switch t.cmp(h.Key, key) {
	case LT:
		if h.left != nil {
			if !h.left.isRed() && !(h.left.left.isRed()) {
				h = h.moveRedLeft()
			}
			h.left = t.remove(h.left, key)
		}
	default:
		if h.left.isRed() {
			h = h.rotateRight()
		}
		if result := t.cmp(h.Key, key); result == EQ && h.right == nil {
			t.Size--
			return nil
		}

		if h.right != nil {
			if !h.right.isRed() && !(h.right.left.isRed()) {
				h = h.moveRedRight()
			}
			if result := t.cmp(h.Key, key); result == EQ {
				t.Size--
				x := h.right.min()
				h.Key = x.Key
				h.Value = x.Value
				h.right = h.right.removeMin()
			} else {
				h.right = t.remove(h.right, key)
			}
		}
	}
	return h.fixUp()
}
