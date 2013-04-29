package gorbtree

import (
	"../"
	"fmt"
	"sync"
)

var _ = fmt.Println

const (
	Red   = true
	Black = false
	Left  = true
	Right = false
)

type rbNode struct {
	left, right *rbNode
	key, value  interface{}
	color       bool
}

type RbTree struct {
	Height      int
	Size        int
	first, last *rbNode
	root        *rbNode
	cmp         gotree.CompareFunc
	lock        sync.RWMutex
}

func (n *rbNode) Key() interface{} {
	return n.key
}

func (n *rbNode) Value() interface{} {
	return n.value
}

func (n *rbNode) MinChild() *rbNode {
	for n.left != nil {
		n = n.left
	}
	return n

}

func (n *rbNode) MaxChild() *rbNode {
	for n.right != nil {
		n = n.right
	}
	return n
}

// Cant modify the tree as we go down
// for if the stack gets out of sync as we head
// back up it will be all bad
func (t *RbTree) Traverse(f gotree.IterFunc) {
	node := t.root
	stack := []*rbNode{nil}
	for len(stack) != 1 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.left
		} else {
			stackIndex := len(stack) - 1
			node = stack[stackIndex]
			f(node.key, node.value)
			stack = stack[0:stackIndex]
			//fmt.Println("stack size", len(stack))
			node = node.right
		}
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

func New(cmp gotree.CompareFunc) *RbTree {
	if cmp == nil {
		panic("Must define a compare function")
	}
	return &RbTree{root: nil, first: nil, last: nil, cmp: cmp}
}

func (t *RbTree) Search(key interface{}) (value interface{}, ok bool) {

	if key == nil {
		return

	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	x := t.root
	for x != nil {
		cmp := t.cmp(x.key, key)
		if cmp == 0 {
			return x.value, true
		} else if cmp > 0 {
			x = x.left
		} else {
			x = x.right
		}
	}
	return
}

func (t *RbTree) Insert(key interface{}, value interface{}) (old interface{}, ok bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	//t.root = t.insert(t.root, key, value)
	t.root = t.insertIter(t.root, key, value)
	if t.root.color == Red {
		t.Height++
	}
	t.root.color = Black
	return
}

func (t *RbTree) insert(h *rbNode, key interface{}, value interface{}) *rbNode {

	// empty tree
	if h == nil {
		return &rbNode{color: Red, key: key, value: value}
	}

	switch cmp := t.cmp(h.key, key); {
	case cmp == 0:
		h.value = value
	case cmp > 0:
		h.left = t.insert(h.left, key, value)
	default:
		h.right = t.insert(h.right, key, value)
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

func (t *RbTree) insertIter(h *rbNode, key interface{}, value interface{}) *rbNode {

	stack := []*rbNode{}
	count := 0
	for {
		// empty tree
		if h == nil {
			h = &rbNode{color: Red, key: key, value: value}
			stack = append(stack, h)
			count++
			break
		}
		stack = append(stack, h)
		count++
		switch cmp := t.cmp(h.key, key); {
		case cmp == 0:
			h.value = value
			break
		case cmp > 0:
			h = h.left
		default:
			h = h.right
		}
	}
	for count > 0 {
		count--
		h = stack[count]
		if h.right.isRed() && !(h.left.isRed()) {
			h = h.rotateLeft()
		}
		if h.left.isRed() && h.left.left.isRed() {
			h = h.rotateRight()
		}

		if h.left.isRed() && h.right.isRed() {
			h.colorFlip()
		}
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
