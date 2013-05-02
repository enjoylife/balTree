package gorbtree

import (
	"fmt"
	"gotree"
	"sync"
)

type color bool

const (
	Red   color = true
	Black color = false
)

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

type rbNode struct {
	left, right *rbNode
	Key, Value  interface{}
	color       color
}

type rbIter struct {
	current *rbNode
	stack   []*rbNode
	next    func() *rbNode
}

type RbTree struct {
	Height      int
	Size        int
	first, last *rbNode
	iter        rbIter
	root        *rbNode
	cmp         gotree.CompareFunc
	lock        sync.RWMutex
}

func New(cmp gotree.CompareFunc) *RbTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &RbTree{root: nil, first: nil, last: nil, cmp: cmp}
}

func (t *RbTree) Next() *rbNode {
	return t.iter.next()
}

func (t *RbTree) Min() (key interface{}, value interface{}) {
	n := t.root
	for n.left != nil {
		n = n.left
	}
	return n.Key, n.Value

}

func (t *RbTree) Max() (key interface{}, value interface{}) {
	n := t.root
	for n.right != nil {
		n = n.right
	}
	return n.Key, n.Value

}

func (t *RbTree) InitIter(order gotree.TravOrder) *rbNode {
	current := t.root
	stack := []*rbNode{}
	switch order {
	case gotree.InOrder:
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

	case gotree.PreOrder:
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
	case gotree.PostOrder:
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

func (t *RbTree) Traverse(order gotree.TravOrder, f gotree.IterFunc) {

	n := t.root
	switch order {
	case gotree.InOrder:
		var inorder func(node *rbNode)
		inorder = func(node *rbNode) {
			if node == nil {
				return
			}
			inorder(node.left)
			f(node.Key, node.Value)
			inorder(node.right)
		}
		inorder(n)
	case gotree.PreOrder:
		var preorder func(node *rbNode)
		preorder = func(node *rbNode) {
			if node == nil {
				return
			}
			f(node.Key, node.Value)
			preorder(node.left)
			preorder(node.right)
		}
		preorder(n)
	case gotree.PostOrder:
		var postorder func(node *rbNode)
		postorder = func(node *rbNode) {
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

		switch t.cmp(x.Key, key) {
		case gotree.EQ:
			return x.Value, true
		case gotree.LT:
			x = x.left
		case gotree.GT:
			x = x.right
		default:
			panic("Compare result of undefined")
		}
	}
	return
}

func (t *RbTree) Insert(key interface{}, value interface{}) (old interface{}, ok bool) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.root == nil {
		t.root = &rbNode{color: Red, Key: key, Value: value, left: nil, right: nil}
	} else {
		t.root = t.insert(t.root, key, value)
		//t.root = t.insertIter(t.root, key, value)
	}
	if t.root.color == Red {
		t.Height++
	}
	t.root.color = Black
	return
}

func (t *RbTree) insert(h *rbNode, key interface{}, value interface{}) *rbNode {

	// empty tree
	if h == nil {
		return &rbNode{color: Red, Key: key, Value: value}
	}

	switch t.cmp(h.Key, key) {
	case gotree.EQ:
		h.Value = value
	case gotree.LT:
		h.left = t.insert(h.left, key, value)
	case gotree.GT:
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

func (t *RbTree) insertIter(h *rbNode, key interface{}, value interface{}) *rbNode {

	// empty tree
	if h == nil {
		return &rbNode{color: Red, Key: key, Value: value}
	}

	// setup our own stack and helpers
	var (
		stack     = []*rbNode{}
		count int = 0
		prior *rbNode
	)

L:
	for {
		switch t.cmp(h.Key, key) {
		case gotree.EQ:
			h.Value = value
			return t.root // no need for rest of the fix code
		case gotree.LT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.left
			if h == nil {
				h = &rbNode{color: Red, Key: key, Value: value}
				prior.left = h
				break L
			}
		case gotree.GT:
			prior = h
			stack = append(stack, prior)
			count++
			h = h.right
			if h == nil {
				h = &rbNode{color: Red, Key: key, Value: value}
				prior.right = h
				break L
			}
		default:
			panic("Compare result undefined")
		}

		if prior == h {
			panic("Shouldn't be equal last check")
		}

	}

	// h is parent of new node at this point
	h = prior
L2:
	for {
		count--

		if h.right.isRed() && !(h.left.isRed()) {
			h = h.rotateLeft()
		}
		if h.left.isRed() && h.left.left.isRed() {
			h = h.rotateRight()
		}
		if h.left.isRed() && h.right.isRed() {
			h.colorFlip()
		}

		if count == 0 {
			break L2
		}

		if count > 0 {

			switch t.cmp(stack[count-1].Key, h.Key) {
			case gotree.LT:
				stack[count-1].left = h
			case gotree.GT:
				stack[count-1].right = h
			}
			h = stack[count-1]
		}

	}

	return h
}

func (t *RbTree) TraverseIter(f gotree.IterFunc) {
	node := t.root
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.left
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.Key, node.Value)
			stack = stack[0:stackIndex]
			node = node.right
		}
	}
}
